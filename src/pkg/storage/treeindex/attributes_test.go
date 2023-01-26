package treeindex

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack/v5"
)

func TestAttributes_EncodeDecodeMsgpack(t *testing.T) {
	tests := []struct {
		name       string
		attributes Attributes
	}{
		{
			"1. empty map",
			Attributes{},
		},
		{
			"2. simple node value",
			Attributes{
				"value_node": newValueNode("123"),
			},
		},
		{
			"3. multi values",
			Attributes{
				"key1": newValueNode("123"),
				"key2": newValueNode(12.34),
			},
		},
		{
			"4. node with slice node",
			Attributes{
				"slice": &SliceNode{
					BaseNode: BaseNode{
						NodeType: SliceNodeType,
					},
					Data: Attributes{
						"key1": newValueNode("123"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := msgpack.Marshal(tt.attributes)
			assert.Nil(t, err)

			got := Attributes{}
			assert.Nil(t, msgpack.Unmarshal(data, &got))

			assert.Equal(t, tt.attributes, got)
		})
	}
}
