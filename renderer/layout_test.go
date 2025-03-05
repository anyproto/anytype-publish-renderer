package renderer

import (
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeRenderLayoutParams(t *testing.T) {
	id := "div-66c82ef37e4bcdd7891d4276"
	r := NewTestRenderer(WithBlocksById(makeLayout(id)))
	layoutBlock := r.BlocksById[id]

	expected := &BlockParams{
		Id:          id,
		Classes:     []string{"block", "align0", "blockLayout", "layoutDiv"},
		ChildrenIds: []string{"id1", "id2"},
	}

	actual := r.makeLayoutBlockParams(layoutBlock)
	assert.Equal(t, expected.Id, actual.Id)
	assert.Equal(t, expected.Classes, actual.Classes)
	assert.EqualValues(t, expected.ChildrenIds, actual.ChildrenIds)
}

func makeLayout(id string) map[string]*model.Block {
	return map[string]*model.Block{
		id: {
			Id:          id,
			ChildrenIds: []string{"id1", "id2"},
			Content: &model.BlockContentOfLayout{
				Layout: &model.BlockContentLayout{Style: model.BlockContentLayout_Div},
			},
		},
		"id1": {
			Id:      "id1",
			Content: &model.BlockContentOfText{Text: &model.BlockContentText{Text: "test1"}},
		},
		"id2": {
			Id:      "id2",
			Content: &model.BlockContentOfText{Text: &model.BlockContentText{Text: "test2"}},
		},
	}
}
