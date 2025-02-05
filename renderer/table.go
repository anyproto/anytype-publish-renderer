package renderer

import (
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
)

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

const DEFAULT_COLUMN_WIDTH = 140

func (r *Renderer) MakeRenderTableParams(b *model.Block) (params *RenderTableParams) {

	var columnSizes []string
	columns := r.BlocksById[b.ChildrenIds[0]]

	for _, colId := range columns.ChildrenIds {
		col := r.BlocksById[colId]
		fields := col.GetFields()
		width := pbtypes.GetInt64(fields, "width")
		if width == 0 {
			width = DEFAULT_COLUMN_WIDTH
		}
		columnSizes = append(columnSizes, fmt.Sprintf("%dpx", width))
	}

	rows := r.BlocksById[b.ChildrenIds[1]]

	classes := []string{"block", "blockTable"}
	if b.BackgroundColor != "" {
		classes = append(classes, fmt.Sprintf("bgColor bgColor-%s", b.BackgroundColor))
	}

	params = &RenderTableParams{
		Classes:     strings.Join(classes, " "),
		Id:          b.Id,
		Rows:        rows,
		Columns:     columns,
		ColumnSizes: strings.Join(columnSizes, " "),
	}

	return
}

func (r *Renderer) RenderTable(b *model.Block) templ.Component {
	params := r.MakeRenderTableParams(b)
	return TableTemplate(r, params)
}

func (r *Renderer) MakeRenderTableRowCellParams(b *model.Block) (params *RenderTableRowCellParams) {

	textComp := r.RenderBlock(b.Id)
	params = &RenderTableRowCellParams{
		Classes:  "",
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
