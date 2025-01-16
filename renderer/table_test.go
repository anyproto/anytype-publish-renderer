package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeRenderTable(t *testing.T) {
	t.Run("default column sizes", func(t *testing.T) {
		r := getTestRenderer("test-tables")
		id := "67892d2f0e5b176d4fd8fd25"
		tableBlock := r.BlocksById[id]

		expected := &RenderTableParams{
			ColumnSizes: "72px 140px 140px",
		}

		actual := r.MakeRenderTableParams(tableBlock)
		assert.EqualValues(t, expected.ColumnSizes, actual.ColumnSizes)
	})
}
