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

	Rows *model.Block
}

type RenderTableRowCellParams struct {
	Classes string
	Id      string

	TextComp templ.Component
}

func (r *Renderer) MakeRenderTableParams(b *model.Block) (params *RenderTableParams) {

	var columnSizes []string
	columnIds := r.BlocksById[b.ChildrenIds[0]].ChildrenIds
	for _, colId := range columnIds {
		col := r.BlocksById[colId]
		fields := col.GetFields()
		width := pbtypes.GetInt64(fields, "width")
		columnSizes = append(columnSizes, fmt.Sprintf("%dpx", width))
	}

	rows := r.BlocksById[b.ChildrenIds[1]]
	var classes string
	if b.BackgroundColor != "" {
		classes = fmt.Sprintf("bgColor bgColor-%s", b.BackgroundColor)
	}
	params = &RenderTableParams{
		Classes:     classes,
		Id:          "",
		Rows:        rows,
		ColumnSizes: strings.Join(columnSizes, " "),
	}

	return
}

func (r *Renderer) RenderTable(b *model.Block) templ.Component {
	params := r.MakeRenderTableParams(b)
	return TableTemplate(r, params)
}

func (r *Renderer) MakeRenderTableRowCellParams(b *model.Block) (params *RenderTableRowCellParams) {

	textComp := r.RenderText(b)
	params = &RenderTableRowCellParams{
		Classes:  "",
		Id:       b.Id,
		TextComp: textComp,
	}
	return
}

func (r *Renderer) RenderTableRowCell(cellId string) templ.Component {
	cellBlock := r.BlocksById[cellId]
	params := r.MakeRenderTableRowCellParams(cellBlock)
	return TableRowCellTemplate(r, params)
}

func gridSizes(sizes string) templ.SafeCSSProperty {
	return templ.SafeCSSProperty(sizes)
}

func (r *Renderer) rowHeaderClass(rowId string) string {
	var headerClass string
	if r.BlocksById[rowId].GetTableRow().IsHeader {
		headerClass = "isHeader"
	}
	return headerClass
}
