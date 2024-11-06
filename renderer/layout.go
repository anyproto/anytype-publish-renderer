package renderer

import (
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"

	"github.com/a-h/templ"
)
type LayoutRenderParams struct {
	Classes string
	Id string
	ChildrenIds []string
}

func (r *Renderer) RenderLayout(b *model.Block) templ.Component {
	layoutClass := "layout" + b.GetLayout().GetStyle().String()
	params := LayoutRenderParams{
		Id: "block" + b.Id,
		Classes: "block blockLayout " + layoutClass,
		ChildrenIds: b.ChildrenIds,
	}

	return LayoutTemplate(r, &params)
}
