package service

import (
	"bytes"
	"context"
	"fmt"
	"hms/gateway/pkg/common"
	"strings"
	"time"

	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/scrypt"

	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service/processing"
	proc "hms/gateway/pkg/docs/service/processing"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/infrastructure"
)

type Service struct {
	Infra *infrastructure.Infra
	Proc  *processing.Proc
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

const (
	N      = 1048576
	r      = 8
	p      = 1
	keyLen = 32
)

func NewUserService(cfg *config.Config, infra *infrastructure.Infra) *Service {
	process := processing.New(
		infra.LocalDB,
		infra.EthClient,
		infra.FilecoinClient,
		infra.IpfsClient,
		cfg.Storage.Localfile.Path,
	)

	process.Start()

	return &Service{
		Infra: infra,
		Proc:  process,
	}
}

func (s *Service) Register(ctx context.Context, procRequest *proc.Request, user *model.UserCreateRequest) (err error) {
	ehrSystemID := ctx.(*gin.Context).GetString("ehrSystemID")
	address, pwdHash, err := s.getUserAddressAndHash(ehrSystemID, user.UserID, user.Password)
	if err != nil {
		return err
	}

	requestID := ctx.(*gin.Context).GetString("reqId")

	txHash, err := s.Infra.Index.UserAdd(requestID, address, user.UserID, user.Role, pwdHash)
	if err != nil {
		return fmt.Errorf("Index.UserAdd error: %w", err)
	}

	procRequest.AddEthereumTx(proc.TxDeleteDoc, txHash)

	return nil
}

func (s *Service) Login(ctx context.Context, user *model.UserAuthRequest) (err error) {
	ehrSystemID := ctx.(*gin.Context).GetString("ehrSystemID")
	address, pwdHash, err := s.getUserAddressAndHash(ehrSystemID, user.UserID, user.Password)
	if err != nil {
		return err
	}

	userHash, err := s.Infra.Index.GetUserPasswordHash(ctx, address)
	if err != nil {
		return fmt.Errorf("Login.GetUserPasswordHash error: %w", err)
	}

	if bytes.Compare(pwdHash, userHash) != 0 {
		return errors.ErrFieldIsIncorrect("Password")
	}

	return nil
}

func (s *Service) getUserAddressAndHash(ehrSystemID, userID, password string) (eth_common.Address, []byte, error) {
	_, userPrivateKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return eth_common.Address{}, nil, fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	privateUserKey := userPrivateKey[:]

	privateKey, err := crypto.ToECDSA(privateUserKey)
	if err != nil {
		return eth_common.Address{}, nil, fmt.Errorf("crypto.ToECDSA error: %w userID %s", err, userID)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	pwdHash, err := s.generateHash(privateUserKey, userID, ehrSystemID, password)
	if err != nil {
		return eth_common.Address{}, nil, fmt.Errorf("register s.generateHash error: %w userID %s, password: %s", err, userID, password)
	}
	return address, pwdHash, nil
}

// method should be idempotent
func (s *Service) generateHash(salt []byte, phrases ...string) ([]byte, error) {
	hash := strings.Join(phrases, "")

	pwdHash, err := scrypt.Key([]byte(hash), salt, N, r, p, keyLen)
	if err != nil {
		return nil, fmt.Errorf("scrypt.Key error: %w", err)
	}

	return pwdHash, nil
}

func (s *Service) CreateToken(userID string) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(common.JWTExpires).Unix()
	td.AccessUuid = uuid.New().String()

	td.RtExpires = time.Now().Add(common.JWTRefreshExpires).Unix()
	td.RefreshUuid = td.AccessUuid + "++" + userID

	var err error
	//Creating Access Token
	accessTokenSecret, _, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("CreateToken Keystore.Get error: %w userID %s", err, userID)
	}

	refreshTokenSecret, _, err := s.Infra.Keystore.Get(userID + "_refresh")
	if err != nil {
		return nil, fmt.Errorf("CreateRefreshToken Keystore.Get error: %w userID %s", err, userID)
	}

	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["userId"] = userID
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString((*accessTokenSecret)[:])
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString((*refreshTokenSecret)[:])
	if err != nil {
		return nil, err
	}

	return td, nil
}

func (s *Service) CreateAuth(userid int64, td *TokenDetails) error {
	//at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	//rt := time.Unix(td.RtExpires, 0)
	//now := time.Now()

	//var client *redis.Client
	//errAccess := client.Set(td.AccessUuid, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	//if errAccess != nil {
	//	return errAccess
	//}
	//errRefresh := client.Set(td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	//if errRefresh != nil {
	//	return errRefresh
	//}
	return nil
}
