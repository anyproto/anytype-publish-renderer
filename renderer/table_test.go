package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeRenderTable(t *testing.T) {
	t.Run("table rendering", func(t *testing.T) {
		r := getTestRenderer("test-tables")
		id := "67892d2f0e5b176d4fd8fd25"
		tableBlock := r.BlocksById[id]

		expected := &BlockParams{Id: id, Classes: []string{"block", "align0", "blockTable"}}

		actual := r.makeTableBlockParams(tableBlock)
		assert.Equal(t, expected.Id, actual.Id)
		assert.Equal(t, expected.Classes, actual.Classes)
	})
}
