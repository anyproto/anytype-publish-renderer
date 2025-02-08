package renderer

import (
	"strconv"
	"strings"

	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"

	"github.com/a-h/templ"
)

type LayoutRenderParams struct {
	Classes     string
	Id          string
	ChildrenIds []string
	Width       float64
}

func (r *Renderer) MakeRenderLayoutParams(b *model.Block) (params *LayoutRenderParams) {
	layoutClass := "layout" + b.GetLayout().GetStyle().String()
	align := "align" + strconv.Itoa(int(b.GetAlign()))
	classes := []string{layoutClass, align}
	fields := b.GetFields()
	width := pbtypes.GetFloat64(fields, "width")

	params = &LayoutRenderParams{
		Id:          b.Id,
		Classes:     strings.Join(classes, " "),
		ChildrenIds: b.ChildrenIds,
		Width:       width,
	}
	return

}
func (r *Renderer) RenderLayout(b *model.Block) templ.Component {
	params := r.MakeRenderLayoutParams(b)
	return LayoutTemplate(r, params)
}
