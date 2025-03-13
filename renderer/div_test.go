package renderer

import (
	"testing"

	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"

	"github.com/stretchr/testify/assert"
)

func TestMakeRenderDivParams(t *testing.T) {
	id := "66c5b61a7e4bcd764b24c213"
	r := NewTestRenderer(
		WithBlocksById(map[string]*model.Block{
			id: {
				Id: id,
				Content: &model.BlockContentOfDiv{
					Div: &model.BlockContentDiv{Style: model.BlockContentDiv_Dots},
				},
			}},
		),
	)

	divBlock := r.BlocksById[id]
	expected := &BlockParams{
		Id:      id,
		Classes: []string{"block", "align0", "blockDiv", "divDot"},
	}

	actual := r.makeRenderDivParams(divBlock)

	assert.Equal(t, expected.Id, actual.Id)
	assert.Equal(t, expected.Classes, actual.Classes)
}
