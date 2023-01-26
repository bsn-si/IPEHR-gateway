package treeindex

import (
	"testing"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack/v5"
)

func TestEHRNode_EncodeDecodeMsgpack(t *testing.T) {
	origin := EHRNode{
		BaseNode: BaseNode{
			ID:   "some_id",
			Type: base.EHRItemType,
			Name: "some_name",
		},
		Attributes:   Attributes{},
		Compositions: Container{},
	}

	data, err := msgpack.Marshal(origin)
	assert.Nil(t, err)

	got := EHRNode{}
	assert.Nil(t, msgpack.Unmarshal(data, &got))

	assert.Equal(t, origin, got)
}
