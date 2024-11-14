package renderer

import (
	"strconv"
	"strings"

	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"

	"github.com/a-h/templ"
)

type LayoutRenderParams struct {
	Classes     string
	Id          string
	ChildrenIds []string
}

func (r *Renderer) MakeRenderLayoutParams(b *model.Block) (params *LayoutRenderParams) {
	layoutClass := "layout" + b.GetLayout().GetStyle().String()
	align := "align" + strconv.Itoa(int(b.GetAlign()))
	classes := []string{layoutClass, align}
	params = &LayoutRenderParams{
		Id:          b.Id,
		Classes:     strings.Join(classes, " "),
		ChildrenIds: b.ChildrenIds,
	}
	return

}
func (r *Renderer) RenderLayout(b *model.Block) templ.Component {
	params := r.MakeRenderLayoutParams(b)
	return LayoutTemplate(r, params)
}
