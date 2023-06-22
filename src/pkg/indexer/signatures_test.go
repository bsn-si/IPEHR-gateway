package indexer

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func Test_MakeSignature(t *testing.T) {
	tests := []struct {
		name     string
		data     func() []byte
		deadline *big.Int
		pk       string
		want     string
		wantErr  bool
	}{
		{
			"1. make signature correct",
			func() []byte {
				return []byte("Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean massa. ")
			},
			big.NewInt(1687339154),
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
			"ed64670b1ca78aa9fdedf3656a28d941e2452247b791be26209d04dafa666e537107abb0aa0ef67f1b0bd804011a4f0913720fb1c2131cf111feba40208b91451c",
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

			sig, err := makeSignature(tt.data(), pk, tt.deadline)
			if err != nil {
				t.Fatal(err)
			}

			if hex.EncodeToString(sig) != tt.want {
				t.Errorf("Expected signature: %s received: %x", tt.want, sig)
			}
		})
	}
}
