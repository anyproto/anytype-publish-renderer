package renderer

import (
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"

	"github.com/a-h/templ"
)


func (r *Renderer) RenderLayout(b *model.Block) templ.Component {
	return LayoutTemplate(r, b)
}
