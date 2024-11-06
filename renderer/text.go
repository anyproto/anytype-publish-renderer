package renderer

import (
	"strconv"
	"strings"

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
	align := "align" + strconv.Itoa(int(b.GetAlign()))
	classes := []string{"block", "blockText", textClass, align}

	params := TextRenderParams{
		Id: "block-" + b.Id,
		Classes: strings.Join(classes, " "),
		Text: b.GetText().Text,
	}

	return TextTemplate(r, &params)
}
