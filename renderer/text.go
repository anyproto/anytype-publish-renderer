package renderer

import (
	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

type TextRenderParams struct {
	Classes string
	Id string
	Text string
}

func (r *Renderer) RenderText(b *model.Block) templ.Component {
	textClass := "text" + b.GetText().GetStyle().String()
	params := TextRenderParams{
		Id: "block" + b.Id,
		Classes: "block blockText " + textClass,
		Text: b.GetText().Text,
	}

	return TextTemplate(r, &params)
}
