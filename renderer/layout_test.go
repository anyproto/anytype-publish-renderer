package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeRenderLayoutParams(t *testing.T) {
	r := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
	id := "div-66c82ef37e4bcdd7891d4276"
	layoutBlock := r.BlocksById[id]

	expected := &LayoutRenderParams{
		Id:          id,
		Classes:     "layoutDiv align0",
		ChildrenIds: []string{"66c58b0a7e4bcd764b24c1ff", "66c5bf7d7e4bcd764b24c217", "66c58b2a7e4bcd764b24c205", "66c5bf9f7e4bcd764b24c218", "672b6bc10e5b174db34fb5ab", "672b6bc40e5b174db34fb5ae", "r-b50d569fc4be036cdfda9ca4096aeec7", "66c5b5ed7e4bcd764b24c208", "66c5b5f27e4bcd764b24c209", "66c5b5fa7e4bcd764b24c20a", "66c5b5fd7e4bcd764b24c20b", "66c5b60a7e4bcd764b24c20f", "66c5b60b7e4bcd764b24c210", "66c5b6117e4bcd764b24c211", "66c5b6297e4bcd764b24c214", "66c5b62c7e4bcd764b24c215", "66c6f8047e4bcd7bc81f3f0a", "66c829e37e4bcdd7891d422e", "66c829e97e4bcdd7891d422f"},
	}

	actual := r.MakeRenderLayoutParams(layoutBlock)
	assert.Equal(t, expected.Id, actual.Id)
	assert.Equal(t, expected.Classes, actual.Classes)
	assert.EqualValues(t, expected.ChildrenIds, actual.ChildrenIds)
}
