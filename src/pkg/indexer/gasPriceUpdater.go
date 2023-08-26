package indexer

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
)

func (i *Index) gasPriceUpdater(gasMultiplier float64) {
	ticker := time.NewTicker(common.GasPriceUpdateInterval)
	defer ticker.Stop()

	for {
		<-ticker.C
		i.gasPriceUpdate(gasMultiplier)
	}
}

func (i *Index) gasPriceUpdate(gasMultiplier float64) {
	timeout := common.GasPriceUpdateInterval / 2

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	gasPrice, err := i.client.SuggestGasPrice(ctx)
	if err != nil {
		log.Println("[INDEX] GasPrice update error:", err)
		return
	}

	if gasMultiplier != 0 {
		gasPriceFloat := new(big.Float).SetInt(gasPrice)
		gasPriceFloat.Mul(gasPriceFloat, big.NewFloat(gasMultiplier))
		gasPrice, _ = gasPriceFloat.Int(new(big.Int))
	}

	i.Lock()
	i.transactOpts.GasFeeCap = gasPrice
	i.Unlock()

	log.Println("[INDEX] GasPrice:", gasPrice)
}
