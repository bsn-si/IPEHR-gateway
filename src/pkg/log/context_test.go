package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	t.Run("context with values", func(t *testing.T) {
		expected := Fields{
			"foo": "bar",
		}
		ctx := ContextWithFields(context.Background(), expected)
		actual := FieldsFromContext(ctx)
		assert.EqualValues(t, expected, actual)
	})
	t.Run("empty context", func(t *testing.T) {
		assert.Nil(t, FieldsFromContext(context.Background()))
	})
	t.Run("nil context", func(t *testing.T) {
		assert.Nil(t, FieldsFromContext(nil)) //nolint staticcheck
	})
}
