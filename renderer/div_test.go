package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeRenderDivParams(t *testing.T) {
	r := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
	divBlock := r.BlocksById["66c5b61a7e4bcd764b24c213"]

	expected := &BlockParams{
		Id:      "66c5b61a7e4bcd764b24c213",
		Classes: []string{"block", "align0", "blockDiv", "divDot"},
	}

	actual := r.makeRenderDivParams(divBlock)

	assert.Equal(t, expected.Id, actual.Id)
	assert.Equal(t, expected.Classes, actual.Classes)
}
