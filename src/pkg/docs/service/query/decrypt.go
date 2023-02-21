package query

import (
	"fmt"
	"math/big"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/chachaPoly"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/hm"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/storage/treeindex"
)

func decryptQueryResponse(resp *model.QueryResponse, key *chachaPoly.Key, nonce *chachaPoly.Nonce) error {
	for i, row := range resp.Rows {
		var tmp map[string]any

		switch row := row.(type) {
		case map[string]any:
			tmp = row
		default:
			return errors.Errorf("Unsupported QueryResponse.Rows type: %T. Row index: %d", row, i)
		}

		fields := map[string]*struct {
			value    any
			required bool
		}{
			"Data":        {nil, true},
			"IsEncrypted": {nil, true},
			"DataType":    {nil, false},
		}

		for key, field := range fields {
			val, ok := tmp[key]
			if !ok && field.required {
				return errors.Errorf("Field %s is missing. Row index: %d", key, i)
			}

			field.value = val
		}

		switch isEncrypted := fields["IsEncrypted"].value.(type) {
		case bool:
			if !isEncrypted {
				resp.Rows[i] = fields["Data"].value
				continue
			}
		default:
			return errors.Errorf("Unsupported QueryResponse.Rows. Row index %d IsEncryped type: %T", i, isEncrypted)
		}

		var data []byte

		switch val := fields["Data"].value.(type) {
		case []byte:
			data = val
		default:
			return errors.Errorf("Unsupported QueryResponse.Row type for encrypted data: %T. Row index: %d", val, i)
		}

		var dataType treeindex.DataType

		switch val := fields["DataType"].value.(type) {
		case treeindex.DataType:
			dataType = val
		case nil:
			dataType = treeindex.DataTypeBasic
		default:
			return errors.Errorf("Unsupported QueryResponse.Rows. Row index: %d DataType type: %T", i, val)
		}

		var err error

		switch dataType {
		case treeindex.DataTypeBasic:
			resp.Rows[i], err = hm.DecryptString(data, key, nonce)
			if err != nil {
				return fmt.Errorf("DecryptString error: %w row index: %d", err, i)
			}
		case treeindex.DataTypeBigInt:
			val := new(big.Int).SetBytes(data)

			resp.Rows[i], err = hm.DecryptInt(val, key)
			if err != nil {
				return fmt.Errorf("DecryptString error: %w row index: %d", err, i)
			}
		case treeindex.DataTypeBigFloat:
			val := new(big.Float)

			err = val.UnmarshalText(data)
			if err != nil {
				return fmt.Errorf("bigFloat UnmarshalText error: %w. Row index: %d", err, i)
			}

			resp.Rows[i], err = hm.DecryptFloat(val, key)
			if err != nil {
				return fmt.Errorf("DecryptString error: %w row index: %d", err, i)
			}
		default:
			return errors.Errorf("Unsupported QueryResponse.Row DataType: %v. Row index: %d", dataType, i)
		}
	}

	return nil
}
