package syncer

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"ipehr/stat/pkg/localDB"
)

type Config struct {
	Endpoint   string
	StartBlock uint64
	Contracts  []struct {
		Name    string
		Address string
		AbiPath string
	}
}

type Syncer struct {
	db        *localDB.DB
	ethClient *ethclient.Client
	addrList  sync.Map
	ehrABI    *abi.ABI
	usersABI  *abi.ABI
	blockNum  *big.Int
}

const (
	BlockNotFoundTimeout = time.Second * 15
	BlockGetErrorTimeout = time.Second * 30
)

func New(db *localDB.DB, ethClient *ethclient.Client, cfg Config) *Syncer {
	s := Syncer{
		db:        db,
		ethClient: ethClient,
		addrList:  sync.Map{},
		blockNum:  big.NewInt(int64(cfg.StartBlock)),
	}

	lastBlock, err := db.SyncLastBlockGet()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = db.SyncLastBlockSet(cfg.StartBlock)
			if err != nil {
				log.Fatal("[SYNC] SyncLastBlockSet error: ", err)
			}
		} else {
			log.Fatal("SyncLastBlockGet error: ", err)
		}
	}

	if lastBlock > s.blockNum.Uint64() {
		s.blockNum = big.NewInt(int64(lastBlock))
	}

	for _, c := range cfg.Contracts {
		s.addrList.Store(common.HexToAddress(c.Address), c.Name)

		abiJSON, err := os.ReadFile(c.AbiPath)
		if err != nil {
			log.Fatalf("abiSON read file '%s' error: %v", c.AbiPath, err)
		}

		abi, err := abi.JSON(bytes.NewReader(abiJSON))
		if err != nil {
			log.Fatal("abi.JSON error: ", err)
		}

		switch c.Name {
		case "ehrIndex":
			s.ehrABI = &abi
		case "users":
			s.usersABI = &abi
		}
	}

	return &s
}

func (s *Syncer) Start() {
	var bigInt1 = big.NewInt(1)

	log.Printf("[SYNC] Starting sync from block number: %d", s.blockNum)

	ctx := context.Background()

	go func() {
		for {
			// get the full block details, using a custom jsonrpc ID as a test
			block, err := s.ethClient.BlockByNumber(ctx, s.blockNum)
			if err != nil {
				if err.Error() == "not found" {
					time.Sleep(BlockNotFoundTimeout)
					continue
				} else {
					log.Printf("[SYNC] Block %d %v get error:", s.blockNum, err)
					log.Printf("[SYNC] BlockByNumber error: %v Sleeping %s...", err, BlockGetErrorTimeout)
					time.Sleep(BlockGetErrorTimeout)
					continue
				}
			}

			ts := time.Unix(int64(block.Time()), 0)

			for _, blockTx := range block.Transactions() {
				if blockTx.To() == nil {
					// contract creation
					continue
				}

				contractName, ok := s.addrList.Load(*blockTx.To())
				if !ok {
					continue
				}

				receipt, err := s.ethClient.TransactionReceipt(ctx, blockTx.Hash())
				if err != nil {
					log.Printf("[SYNC] tx %s receipt get error: %v", blockTx.Hash().String(), err)
				}

				if receipt.Status == types.ReceiptStatusFailed {
					continue
				}

				decodedSig := blockTx.Data()[:4]
				decodedData := blockTx.Data()[4:]

				var _abi *abi.ABI

				switch contractName {
				case "ehrIndex":
					_abi = s.ehrABI
				case "users":
					_abi = s.usersABI
				}

				method, err := _abi.MethodById(decodedSig)
				if err != nil {
					log.Println("abi.MethodById error: ", err)
					continue
				}

				switch method.Name {
				case "multicall":
					err = s.procMulticall(_abi, method, decodedData, ts)
					if err != nil {
						log.Fatal("[SYNC] procMulticall error: ", err)
					}
				case "addEhrDoc":
					err = s.procAddEhrDoc(method, decodedData, ts)
					if err != nil {
						log.Fatal("[SYNC] procAddEhrDoc error: ", err)
					}
				case "userNew":
					err = s.procUserNew(method, decodedData, ts)
					if err != nil {
						log.Fatal("[SYNC] procUserNew error: ", err)
					}
				}
			}

			log.Printf("[SYNC] new block %v %v txs %d", block.Number().Int64(), time.Unix(int64(block.Time()), 0).Format("2006-01-02 15:04:05"), len(block.Transactions()))

			err = s.db.SyncLastBlockSet(s.blockNum.Uint64())
			if err != nil {
				log.Fatal("[SYNC] SyncLastBlockSet error: ", err)
			}

			s.blockNum.Add(s.blockNum, bigInt1)
		}
	}()
}

func (s *Syncer) procMulticall(_abi *abi.ABI, method *abi.Method, inputData []byte, ts time.Time) error {
	args, err := method.Inputs.Unpack(inputData)
	if err != nil {
		return fmt.Errorf("UnpackValues error: %w", err)
	}

	for _, m := range args[0].([][]byte) {
		decodedSig := m[:4]
		decodedData := m[4:]

		method, err = _abi.MethodById(decodedSig)
		if err != nil {
			return fmt.Errorf("abi.MethodById error: %w", err)
		}

		switch method.Name {
		case "addEhrDoc":
			err = s.procAddEhrDoc(method, decodedData, ts)
			if err != nil {
				return fmt.Errorf("procAddEhrDoc error: %w", err)
			}
		case "userNew":
			err = s.procUserNew(method, decodedData, ts)
			if err != nil {
				return fmt.Errorf("procUserNew error: %w", err)
			}
		}
	}

	return nil
}

func (s *Syncer) procAddEhrDoc(method *abi.Method, inputData []byte, ts time.Time) error {
	log.Println("[STAT] new EHR document registered")

	err := s.db.StatDocumentsCountIncrement(ts)
	if err != nil {
		return fmt.Errorf("StatDocumentsCountIncrement error: %w", err)
	}

	return nil
}

func (s *Syncer) procUserNew(method *abi.Method, inputData []byte, ts time.Time) error {
	log.Println("[STAT] new patient registered")

	err := s.db.StatPatientsCountIncrement(ts)
	if err != nil {
		return fmt.Errorf("StatPatientsCountIncrement error: %w", err)
	}

	return nil
}
