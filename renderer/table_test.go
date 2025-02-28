package renderer

import (
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeRenderTable(t *testing.T) {
	t.Run("table rendering", func(t *testing.T) {
		id := "table"
		r := NewTestRenderer(WithBlocksById(makeTable(id)))
		tableBlock := r.BlocksById[id]

		expected := &BlockParams{Id: id, Classes: []string{"block", "align0", "blockTable"}}

		actual := r.makeTableBlockParams(tableBlock)
		assert.Equal(t, expected.Id, actual.Id)
		assert.Equal(t, expected.Classes, actual.Classes)
	})
}

func makeTable(id string) map[string]*model.Block {
	return map[string]*model.Block{
		id: {
			Id:          id,
			ChildrenIds: []string{"id1", "id2"},
			Content: &model.BlockContentOfTable{
				Table: &model.BlockContentTable{},
			},
		},
		"id1": {
			Id:          "id1",
			ChildrenIds: []string{"column1", "column2"},
			Content:     &model.BlockContentOfLayout{Layout: &model.BlockContentLayout{Style: model.BlockContentLayout_TableColumns}},
		},
		"id2": {
			Id:          "id2",
			ChildrenIds: []string{"row1"},
			Content:     &model.BlockContentOfLayout{Layout: &model.BlockContentLayout{Style: model.BlockContentLayout_TableRows}},
		},
		"row1": {
			Id:          "row1",
			ChildrenIds: []string{"row1-column1", "row1-column2"},
			Content:     &model.BlockContentOfTableRow{TableRow: &model.BlockContentTableRow{}},
		},
		"row1-column1": {
			Id:      "row1-column1",
			Content: &model.BlockContentOfText{Text: &model.BlockContentText{Text: "test1"}},
		},
		"row1-column2": {
			Id:      "row1-column2",
			Content: &model.BlockContentOfText{Text: &model.BlockContentText{Text: "test2"}},
		},
		"column1": {
			Id:      "column1",
			Content: &model.BlockContentOfTableColumn{TableColumn: &model.BlockContentTableColumn{}},
		},
		"column2": {
			Id:      "column2",
			Content: &model.BlockContentOfTableColumn{TableColumn: &model.BlockContentTableColumn{}},
		},
	}
}
