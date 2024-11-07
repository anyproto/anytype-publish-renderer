package renderer

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-publish-renderer/utils"
	"go.uber.org/zap"
)

const CDN_URL = "https://anytype-static.fra1.cdn.digitaloceanspaces.com"
const bulbEmoji = 0x1F4A1

func emojiUrl(code rune) string {
	return fmt.Sprintf("%s/emojies/%x.png", CDN_URL, code)
}

type TextRenderParams struct {
	Classes      string
	Id           string
	TextComp     templ.Component
	ExternalComp templ.Component
	ChildrenIds  []string
}

func cmpMarks(a, b *model.BlockContentTextMark) int {
	return cmp.Compare(a.Range.From, b.Range.From)
}

func applyMark(s string, mark *model.BlockContentTextMark) string {
	switch mark.Type {
	case model.BlockContentTextMark_Strikethrough:
		return "<markupstrike>" + s + "</markupstrike>"
	case model.BlockContentTextMark_Keyboard:
		return "<markupcode>" + s + "</markupcode>"
	case model.BlockContentTextMark_Italic:
		return "<markupitalic>" + s + "</markupitalic>"
	case model.BlockContentTextMark_Bold:
		return "<markupbold>" + s + "</markupbold>"
	case model.BlockContentTextMark_Underscored:
		return "<markupunderline>" + s + "</markupunderline>"
	case model.BlockContentTextMark_Link:
		url := mark.Param
		tag := fmt.Sprintf(`<a href="%s">`, url)
		return tag + s + "</a>"
	case model.BlockContentTextMark_TextColor:
		color := mark.Param
		tag := fmt.Sprintf(`<markupcolor class="textColor textColor-%s">`, color)
		return tag + s + "</markupcolor>"
	case model.BlockContentTextMark_BackgroundColor:
		color := mark.Param
		tag := fmt.Sprintf(`<markubgpcolor class="bgColor bgColor-%s">`, color)
		return tag + s + "</markupbgcolor>"
	case model.BlockContentTextMark_Mention:
		return "<markupmention>" + s + "</markupmention>"
	case model.BlockContentTextMark_Emoji:
		code := []rune(mark.Param)[0]
		emojiHtml, err := utils.TemplToString(InlineEmojiTemplate(code, "c28"))
		if err != nil {
			log.Error("Failed to render emoji template", zap.Error(err))
			return ""
		} else {
			return emojiHtml
		}
	}

	return "<markupobject>" + s + "</markupobject>"
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

	if bgColor := b.GetBackgroundColor(); bgColor != "" {
		classes = append(classes, "bgColor", "bgColor-"+bgColor)
	}

	text := blockText.Text
	var comp templ.Component

	if style != model.BlockContentText_Code {
		marks := blockText.GetMarks().Marks
		text = applyMarks(text, marks)
		comp = templ.Raw(text)
	} else {
		comp = PlainTextTemplate(text)
	}

	var externalComp templ.Component
	if style == model.BlockContentText_Marked {
		externalComp = BulletMarkerTemplate()
	}
	if style == model.BlockContentText_Callout {
		externalComp = AdditionalEmojiTemplate()
	}

	params := TextRenderParams{
		Id:           "block-" + b.Id,
		Classes:      strings.Join(classes, " "),
		TextComp:     comp,
		ExternalComp: externalComp,
		ChildrenIds:  b.ChildrenIds,
	}

	return TextTemplate(r, &params)
}
