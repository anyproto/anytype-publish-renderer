package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeRenderText(t *testing.T) {
	r := getTestRenderer()
	id := "66c58b2a7e4bcd764b24c205"
	textBlock := r.BlocksById[id]

	expected := &TextRenderParams{
		Id:          id,
		Classes:     "textParagraph align0",
		ChildrenIds: nil,
	}

	actual := r.MakeRenderTextParams(textBlock)
	assert.Equal(t, expected.Id, actual.Id)
	assert.Equal(t, expected.Classes, actual.Classes)
	assert.EqualValues(t, expected.ChildrenIds, actual.ChildrenIds)
}
