package renderer

import (
	"strconv"
	"strings"

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
	align := "align" + strconv.Itoa(int(b.GetAlign()))
	classes := []string{"block", "blockLayout", layoutClass, align}
	params := LayoutRenderParams{
		Id: "block-" + b.Id,
		Classes: strings.Join(classes, " "),
		ChildrenIds: b.ChildrenIds,
	}

	return LayoutTemplate(r, &params)
}
