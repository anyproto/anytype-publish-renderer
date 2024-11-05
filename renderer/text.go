package renderer

import (
	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

func (r *Renderer) RenderText(b *model.Block) templ.Component {
	return TextTemplate(r, b.GetText().Text)
}
