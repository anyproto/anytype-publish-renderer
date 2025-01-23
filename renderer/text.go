package renderer

import (
	"cmp"
	"fmt"
	"html"
	"slices"
	"strconv"
	"strings"
	"unicode/utf16"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"github.com/anyproto/anytype-publish-renderer/renderer/markintervaltree"
	"github.com/anyproto/anytype-publish-renderer/utils"
	"go.uber.org/zap"
)

const bulbEmoji = 0x1F4A1

type TextRenderParams struct {
	Classes     string
	Id          string
	InnerFlex   []templ.Component
	OuterFlex   []templ.Component
	ChildrenIds []string
}

func cmpMarks(a, b *model.BlockContentTextMark) int {
	return cmp.Compare(a.Range.From, b.Range.From)
}

func (r *Renderer) applyMark(s string, mark *model.BlockContentTextMark) string {
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
		tag := fmt.Sprintf(`<markupbgcolor class="bgColor bgColor-%s">`, color)
		return tag + s + "</markupbgcolor>"
	case model.BlockContentTextMark_Mention:
		return "<markupmention>" + s + "</markupmention>"
	case model.BlockContentTextMark_Emoji:
		code := []rune(mark.Param)[0]
		emojiSrc := r.GetEmojiUrl(code)
		emojiHtml, err := utils.TemplToString(InlineEmojiTemplate(emojiSrc, "c28"))
		if err != nil {
			log.Error("Failed to render emoji template", zap.Error(err))
			return ""
		} else {
			return emojiHtml
		}
	}

	return "<markupobject>" + s + "</markupobject>"
}

func StrToUTF16(str string) []uint16 {
	return utf16.Encode([]rune(str))
}

// - make borders
//   - make set from ranges, from-to
//   - sort
//   - for each range, find overlapping intervals
//     add props from each of this ranges to this range
func (r *Renderer) applyNonOverlapingMarks(text string, marks []*model.BlockContentTextMark) string {
	if len(marks) == 0 {
		text = html.EscapeString(text)
		return text
	}

	rText := utf16.Decode(StrToUTF16(text))
	root := &markintervaltree.MarkIntervalTreeNode{
		Mark:        marks[0],
		MaxUpperVal: marks[0].Range.To,
	}

	for i := 1; i < len(marks); i++ {
		root.Insert(marks[i])
	}

	rangeSet := make(map[int32]bool)
	rangeSet[0] = true
	rangeSet[int32(len(rText))] = true
	for _, mark := range marks {
		rangeSet[mark.Range.From] = true
		rangeSet[mark.Range.To] = true
	}

	rangeRay := make([]int32, len(rangeSet))
	i := 0
	for k := range rangeSet {
		rangeRay[i] = k
		i++
	}

	slices.Sort(rangeRay)

	var markedText strings.Builder

	log.Debug("rangeRay", zap.String("ray", fmt.Sprintf("%#v", rangeRay)))
	for i := 0; i < len(rangeRay)-1; i++ {
		curRange := &model.Range{
			From: rangeRay[i],
			To:   rangeRay[i+1],
		}
		marksToApply := make([]*model.BlockContentTextMark, 0)
		markintervaltree.SearchOverlaps(root, curRange, &marksToApply)

		markedPart := string(rText[curRange.From:curRange.To])
		log.Debug("apply marks",
			zap.String("markedPart", markedPart),
			zap.Int32("from", curRange.From),
			zap.Int32("to", curRange.To))
		markedPart = html.EscapeString(markedPart)
		for _, m := range marksToApply {
			markedPart = r.applyMark(markedPart, m)
			log.Debug("apply mark", zap.String("markedPart", markedPart), zap.Int32("from", m.Range.From), zap.Int32("to", m.Range.To))
		}
		log.Debug("final marked part", zap.String("m", markedPart))
		markedText.WriteString(markedPart)
	}

	return markedText.String()
}

func (r *Renderer) MakeRenderTextParams(b *model.Block) (params *TextRenderParams) {
	blockText := b.GetText()
	style := blockText.GetStyle()
	textClass := "text" + style.String()
	align := "align" + strconv.Itoa(int(b.GetAlign()))
	classes := []string{textClass, align}

	if bgColor := b.GetBackgroundColor(); bgColor != "" {
		classes = append(classes, "bgColor", "bgColor-"+bgColor)
	}

	text := blockText.Text
	var textComp templ.Component
	if style != model.BlockContentText_Code {
		marks := blockText.GetMarks().Marks
		text = r.applyNonOverlapingMarks(text, marks)
		textComp = PlainTextWrapTemplate(templ.Raw(text))
	} else {
		fields := b.GetFields()
		lang := pbtypes.GetString(fields, "lang")
		textComp = TextCodeTemplate(text, lang)
	}

	var outerFlex []templ.Component
	var innerFlex []templ.Component
	switch style {
	case model.BlockContentText_Toggle:
		externalComp := ToggleMarkerTemplate()
		innerFlex = append(innerFlex, externalComp, textComp)
	case model.BlockContentText_Numbered:
		number := r.BlockNumbers[b.Id]
		log.Debug("number", zap.Int("num", number), zap.String("id", b.Id))
		externalComp := NumberMarkerTemplate(fmt.Sprintf("%d", number))
		innerFlex = append(innerFlex, externalComp, textComp)
	case model.BlockContentText_Marked:
		externalComp := BulletMarkerTemplate()
		innerFlex = append(innerFlex, externalComp, textComp)
	case model.BlockContentText_Callout:
		emojiSrc := r.GetEmojiUrl(bulbEmoji)
		externalComp := AdditionalEmojiTemplate(emojiSrc)
		innerFlex = append(innerFlex, externalComp, textComp)
	case model.BlockContentText_Quote:
		externalComp := AdditionalQuoteTemplate()
		outerFlex = append(outerFlex, externalComp)
		innerFlex = append(innerFlex, textComp)
	case model.BlockContentText_Checkbox:
		var checkboxComp templ.Component
		if blockText.Checked {
			checkboxComp = CheckboxCheckedTemplate()
		} else {
			checkboxComp = CheckboxUncheckedTemplate()
		}
		innerFlex = append(innerFlex, checkboxComp, textComp)
	default:
		innerFlex = append(innerFlex, textComp)
	}

	params = &TextRenderParams{
		Id:          b.Id,
		Classes:     strings.Join(classes, " "),
		ChildrenIds: b.ChildrenIds,
		OuterFlex:   outerFlex,
		InnerFlex:   innerFlex,
	}
	return

}
func (r *Renderer) RenderText(b *model.Block) templ.Component {
	params := r.MakeRenderTextParams(b)
	return TextTemplate(r, params)
}
