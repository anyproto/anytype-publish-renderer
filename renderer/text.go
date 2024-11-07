package renderer

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"go.uber.org/zap"
)

type TextRenderParams struct {
	Classes     string
	Id          string
	TextComp    templ.Component
	ChildrenIds []string
}

func cmpMarks(a, b *model.BlockContentTextMark) int {
	return cmp.Compare(a.Range.From, b.Range.From)
}

func applyMark(s string, mark *model.BlockContentTextMark) string {
	return "<b>" + s + "</b>"
}

func applyMarks(text string, marks []*model.BlockContentTextMark) string {
	if len(marks) == 0 {
		return text
	}
	var markedText strings.Builder
	var lastPos int32
	rText := []rune(text)
	slices.SortFunc(marks, cmpMarks)
	for _, mark := range marks {
		log.Debug("Marks:", zap.String("pos", fmt.Sprintf("%d: %d-%d", lastPos, mark.Range.From, mark.Range.To)))

		before := rText[lastPos:mark.Range.From]
		markedText.WriteString(string(before))

		markedPart := rText[mark.Range.From:mark.Range.To]
		markedText.WriteString(applyMark(string(markedPart), mark))
		lastPos = mark.Range.To
	}
	return markedText.String()
}

func (r *Renderer) RenderText(b *model.Block) templ.Component {
	blockText := b.GetText()
	style := blockText.GetStyle()
	textClass := "text" + style.String()
	align := "align" + strconv.Itoa(int(b.GetAlign()))
	classes := []string{"block", "blockText", textClass, align}

	text := blockText.Text
	var comp templ.Component

	if style != model.BlockContentText_Code {
		marks := blockText.GetMarks().Marks
		text = applyMarks(text, marks)
		comp = templ.Raw(text)
	} else {
		comp = PlainTextTemplate(text)
	}

	params := TextRenderParams{
		Id:          "block-" + b.Id,
		Classes:     strings.Join(classes, " "),
		TextComp:    comp,
		ChildrenIds: b.ChildrenIds,
	}

	return TextTemplate(r, &params)
}
