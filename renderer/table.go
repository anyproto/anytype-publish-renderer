package renderer

import (
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
)

const DefaultColumnWidth = 140

type RenderTableParams struct {
	Classes string
	Id      string

	ColumnSizes string

	Rows    *model.Block
	Columns *model.Block
}

type RenderTableRowCellParams struct {
	Classes string
	Id      string

	TextComp templ.Component
}

func (r *Renderer) MakeRenderTableParams(b *model.Block) *BlockParams {
	var columnSizes []string
	columns := r.BlocksById[b.ChildrenIds[0]]

	for _, colId := range columns.ChildrenIds {
		col := r.BlocksById[colId]
		fields := col.GetFields()
		width := pbtypes.GetInt64(fields, "width")
		if width == 0 {
			width = DefaultColumnWidth
		}
		columnSizes = append(columnSizes, fmt.Sprintf("%dpx", width))
	}

	rows := r.BlocksById[b.ChildrenIds[1]]

	blockParams := makeDefaultBlockParams(b)
	if b.BackgroundColor != "" {
		blockParams.Classes = append(blockParams.Classes, fmt.Sprintf("bgColor bgColor-%s", b.BackgroundColor))
	}
	tableTemplate := TableTemplate(r, &RenderTableParams{
		Rows:        rows,
		Columns:     columns,
		ColumnSizes: strings.Join(columnSizes, " "),
	})
	blockParams.Content = BlocksWrapper(&BlockWrapperParams{
		Classes: []string{"scrollWrap"},
		Components: []templ.Component{
			BlocksWrapper(&BlockWrapperParams{
				Classes:    []string{"inner"},
				Components: []templ.Component{tableTemplate},
			}),
		},
	})
	return blockParams
}

func (r *Renderer) RenderTable(b *model.Block) templ.Component {
	params := r.MakeRenderTableParams(b)
	return BlockTemplate(r, params)
}

func (r *Renderer) MakeRenderTableRowCellParams(b *model.Block) (params *RenderTableRowCellParams) {
	align := fmt.Sprintf("align-h%d", b.GetAlign())
	vAlign := fmt.Sprintf("align-v%d", b.GetVerticalAlign())
	classes := []string{"cell", align, vAlign}

	textComp := r.RenderBlock(b.Id)
	params = &RenderTableRowCellParams{
		Classes:  strings.Join(classes, " "),
		Id:       b.Id,
		TextComp: textComp,
	}
	return
}

func (r *Renderer) RenderTableRowCell(cellId string) templ.Component {
	cellBlock, ok := r.BlocksById[cellId]
	if !ok {
		return TableRowCellEmptyTemplate()
	}
	params := r.MakeRenderTableRowCellParams(cellBlock)
	return TableRowCellTemplate(r, params)
}

func (r *Renderer) rowHeaderClass(rowId string) string {
	classes := []string{"row"}

	if r.BlocksById[rowId].GetTableRow().IsHeader {
		classes = append(classes, "isHeader")
	}

	return strings.Join(classes, " ")
}
