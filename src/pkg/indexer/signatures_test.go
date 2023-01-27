package indexer

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func Test_MakeSignature(t *testing.T) {
	tests := []struct {
		name    string
		data    func() []byte
		nonce   *big.Int
		pk      string
		want    string
		wantErr bool
	}{
		{
			"1. make signature correct",
			func() []byte {
				return []byte("Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa. ")
			},
			big.NewInt(0),
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
			"8bb3c32716200a33fa59b92878adb2970b379e615780e8eea7bc6cc91b11b3da3f6d1262bd31f48607f0e37cb12d493f44782c646b53c80582c87220aab0b0a91b",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkBytes, err := hex.DecodeString(tt.pk)
			if err != nil {
				t.Fatal(err)
			}

			pk, err := crypto.ToECDSA(pkBytes)
			if err != nil {
				t.Fatal(err)
			}

			sig, err := makeSignature(tt.data(), tt.nonce, pk)
			if err != nil {
				t.Fatal(err)
			}

			if hex.EncodeToString(sig) != tt.want {
				t.Errorf("Expected signature: %s received: %x", tt.want, sig)
			}
		})
	}
}
