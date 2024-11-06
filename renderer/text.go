package renderer

import (
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

type TextRenderParams struct {
	Classes     string
	Id          string
	Text        string
	ChildrenIds []string
}

func processMarks(text string, marks []*model.BlockContentTextMark) string {
	var markedText strings.Builder

	return text
}

func (r *Renderer) RenderText(b *model.Block) templ.Component {
	textClass := "text" + b.GetText().GetStyle().String()
	align := "align" + strconv.Itoa(int(b.GetAlign()))
	classes := []string{"block", "blockText", textClass, align}

	marks := b.GetText().GetMarks().Marks
	textWithMarkup := processMarks(b.GetText().Text, marks)
	params := TextRenderParams{
		Id:          "block-" + b.Id,
		Classes:     strings.Join(classes, " "),
		Text:        textWithMarkup,
		ChildrenIds: b.ChildrenIds,
	}

	return TextTemplate(r, &params)
}
