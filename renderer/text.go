package renderer

import (
	"fmt"
	"html"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode/utf16"

	"github.com/a-h/templ"
	"github.com/gogo/protobuf/types"
	"go.uber.org/zap"

	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"github.com/anyproto/anytype-publish-renderer/renderer/markintervaltree"
	"github.com/anyproto/anytype-publish-renderer/utils"
)

func emojiParam(t model.BlockContentTextStyle) int32 {
	switch t {
	case model.BlockContentText_Header1:
		return 30
	case model.BlockContentText_Header2:
		return 26
	case model.BlockContentText_Header3:
		return 22
	default:
		return 20
	}
}

func applyHeader(style model.BlockContentTextStyle, s string) string {
	var tagName string
	switch style {
	case model.BlockContentText_Header1:
	case model.BlockContentText_Title:
		tagName = "h1"
	case model.BlockContentText_Header2:
		tagName = "h2"
	case model.BlockContentText_Header3:
		tagName = "h3"

	}

	if tagName == "" {
		return s
	}

	return fmt.Sprintf("<%s>%s</%s>", tagName, s, tagName)
}

func (r *Renderer) applyMark(style model.BlockContentTextStyle, s string, mark *model.BlockContentTextMark) string {
	emojiSize := emojiParam(style)

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
		return fmt.Sprintf(`<a href="%s" class="markuplink" target="_blank">`, mark.Param) + s + "</a>"

	case model.BlockContentTextMark_TextColor:
		color := mark.Param
		tag := fmt.Sprintf(`<markupcolor class="textColor textColor-%s">`, color)
		return tag + s + "</markupcolor>"

	case model.BlockContentTextMark_BackgroundColor:
		color := mark.Param
		tag := fmt.Sprintf(`<markupbgcolor class="bgColor bgColor-%s">`, color)
		return tag + s + "</markupbgcolor>"

	case model.BlockContentTextMark_Mention:
		details := r.findTargetDetails(mark.Param)

		if details != nil && len(details.Fields) != 0 {
			iconParams := r.MakeRenderIconObjectParams(details, &IconObjectProps{Size: emojiSize})
			classes := []string{}
			if iconParams.Src != "" || iconParams.SvgSrc != "" {
				classes = append(classes, "withImage")
			}
			if iconParams.SvgSrc != "" {
				iconParams.IconClasses = append(iconParams.IconClasses, "withSvg")
			}
			link := r.makeAnytypeLink(details, mark.Param)
			html, err := utils.TemplToString(TextMarkupMention(r, templ.SafeURL(link), s, classes, iconParams))
			if err != nil {
				log.Error("Failed to render mention icon", zap.Error(err))
			}
			return html
		}

		log.Error("Failed to render mention icon: details are missing", zap.String("mark.Param", mark.Param))
		return ""

	case model.BlockContentTextMark_Emoji:
		code := []rune(mark.Param)[0]
		emojiSrc := r.GetEmojiUrl(code)
		emojiHtml, err := utils.TemplToString(InlineEmojiTemplate(emojiSrc, fmt.Sprintf("c%d", emojiSize)))
		if err != nil {
			log.Error("Failed to render emoji template", zap.Error(err))
			return ""
		} else {
			return emojiHtml
		}
	case model.BlockContentTextMark_Object:
		details := r.findTargetDetails(mark.Param)
		if details == nil || len(details.Fields) == 0 {
			return "<markupobject>" + s + "</markupobject>"
		}
		link := r.makeAnytypeLink(details, mark.Param)
		return fmt.Sprintf(`<a href="%s" class="markuplink" target="_blank">`, link) + s + "</a>"
	}

	return "<markupobject>" + s + "</markupobject>"
}

// Convert a string into "JS-like" rune slices (surrogate pairs split)
//
// When we get Range from anytype-ts, it is calculates emojies by codepoints.
// Which means, that js calculates `":man-woman-boy-girl: asdf".length == 16`, but not 5
func toJSRunes(s string) []rune {
	var jsRunes []rune
	for _, r := range s {
		if r > 0xFFFF {
			// Convert to surrogate pair (two runes)
			high, low := utf16.EncodeRune(r)
			jsRunes = append(jsRunes, rune(high), rune(low))
		} else {
			jsRunes = append(jsRunes, r)
		}
	}

	return jsRunes
}

func fromJSRunes(jsRunes []rune) string {
	var utf16Units []uint16
	for _, r := range jsRunes {
		utf16Units = append(utf16Units, uint16(r))
	}
	// Decode UTF-16 (reconstruct surrogate pairs into full runes)
	runes := utf16.Decode(utf16Units)

	// Convert runes back to a string
	return string(runes)
}

func makeMarksRangeRay(marks []*model.BlockContentTextMark, textLen int32) []int32 {
	rangeSet := make(map[int32]bool)
	rangeSet[0] = true
	rangeSet[textLen] = true
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
	return rangeRay
}

// - make borders
//   - make set from ranges, from-to
//   - sort
//   - for each range, find overlapping intervals
//     add props from each of this ranges to this range
func (r *Renderer) applyNonOverlapingMarks(style model.BlockContentTextStyle, text string, marks []*model.BlockContentTextMark) string {
	if len(marks) == 0 {
		text = html.EscapeString(text)
		return text
	}

	var markedText strings.Builder

	// convert to JSRunes to cut marks.Range in the same way as JS does
	rText := toJSRunes(text)
	marksIntervalTree := markintervaltree.New(marks)
	rtextLen := int32(len(rText))
	rangeRay := makeMarksRangeRay(marks, rtextLen)

	for i := range len(rangeRay) - 1 {
		curRange := &model.Range{
			From: rangeRay[i],
			To:   rangeRay[i+1],
		}
		// skip marks out of range, catch info in logs

		if curRange.From > rtextLen || curRange.To > rtextLen {
			var sb strings.Builder
			sb.WriteString("[ ")
			for _, m := range marks {
				sb.WriteString("{ type: ")
				sb.WriteString(m.Type.String())
				sb.WriteString(", param: ")
				sb.WriteString(m.Param)
				sb.WriteString(", from: ")
				sb.WriteString(strconv.Itoa(int(m.Range.From)))
				sb.WriteString(", to: ")
				sb.WriteString(strconv.Itoa(int(m.Range.To)))
				sb.WriteString("}; ")
			}
			sb.WriteString(" ]")
			log.Error("markup index out of range, skipping markup block",
				zap.String("PublishFilesPath", r.Config.PublishFilesPath),
				zap.Int32("from", curRange.From), zap.Int32("to", curRange.To),
				zap.String("text", text), zap.Int("len", len(text)),
				zap.String("all marks", sb.String()),
			)

			continue
		}

		marksToApply := marksIntervalTree.SearchOverlaps(curRange)
		markedPart := fromJSRunes(rText[curRange.From:curRange.To])
		markedPart = html.EscapeString(markedPart)
		for _, m := range marksToApply {
			markedPart = r.applyMark(style, markedPart, m)
		}

		markedText.WriteString(markedPart)
	}

	return markedText.String()
}

func replaceNewlineBr(text string) string {
	r := regexp.MustCompile(`\r?\n`)
	text = r.ReplaceAllString(text, "<br>")
	return text
}

func (r *Renderer) makeTextBlockParams(b *model.Block) (params *BlockParams) {
	blockText := b.GetText()
	style := blockText.GetStyle()
	if style == model.BlockContentText_Title {
		b.Align = model.BlockAlign(r.LayoutAlign)
	}
	bgColor := b.GetBackgroundColor()
	color := blockText.GetColor()
	iconEmoji := blockText.GetIconEmoji()
	iconImage := blockText.GetIconImage()
	var contentClasses []string
	classes := []string{"text" + style.String()}

	blockParams := makeDefaultBlockParams(b)
	if bgColor != "" {
		if (style == model.BlockContentText_Callout) ||
			(style == model.BlockContentText_Quote) {
			classes = append(classes, "bgColor", "bgColor-"+bgColor)
		} else {
			contentClasses = append(contentClasses, "bgColor", "bgColor-"+bgColor)
		}
	}

	if color != "" {
		contentClasses = append(contentClasses, "textColor", "textColor-"+color)
	}

	text := blockText.Text
	var textComp templ.Component
	if style != model.BlockContentText_Code {
		if blockText.GetMarks() != nil {
			marks := blockText.GetMarks().Marks
			text = r.applyNonOverlapingMarks(style, text, marks)
			text = replaceNewlineBr(text)
		}
		text = applyHeader(style, text)
		textComp = PlainTextWrapTemplate(templ.Raw(text))
	} else {
		fields := b.GetFields()
		lang := pbtypes.GetString(fields, "lang")
		textComp = TextCodeTemplate(text, lang)
	}

	var innerFlex []templ.Component
	switch style {
	case model.BlockContentText_Toggle:
		externalComp := ToggleMarkerTemplate(utils.GetColor(color))
		innerFlex = append(innerFlex, externalComp, textComp)
	case model.BlockContentText_Numbered:
		number := r.BlockNumbers[b.Id]
		externalComp := NumberMarkerTemplate(fmt.Sprintf("%d", number))
		innerFlex = append(innerFlex, externalComp, textComp)
	case model.BlockContentText_Marked:
		externalComp := BulletMarkerTemplate(color)
		innerFlex = append(innerFlex, externalComp, textComp)
	case model.BlockContentText_Callout:
		if iconEmoji == "" && iconImage == "" {
			iconEmoji = "💡"
		}

		details := &types.Struct{
			Fields: map[string]*types.Value{
				bundle.RelationKeyIconEmoji.String():      pbtypes.String(iconEmoji),
				bundle.RelationKeyIconImage.String():      pbtypes.String(iconImage),
				bundle.RelationKeyResolvedLayout.String(): pbtypes.Float64(float64(model.ObjectType_basic)),
			},
		}

		params := r.MakeRenderIconObjectParams(details, &IconObjectProps{Size: 20})
		iconTemplate := IconObjectTemplate(r, params)
		additionalTemplate := AdditionalIconTemplate(iconTemplate)

		innerFlex = append(innerFlex, additionalTemplate, textComp)
	case model.BlockContentText_Title:
		details := r.Sp.Snapshot.Data.GetDetails()
		done := getRelationField(details, bundle.RelationKeyDone, relationToBool)
		additionalTemplate := NoneTemplate("")

		if isTodoLayout(r.ResolvedLayout) {
			iconDetails := &types.Struct{
				Fields: map[string]*types.Value{
					bundle.RelationKeyDone.String():           pbtypes.Bool(done),
					bundle.RelationKeyResolvedLayout.String(): pbtypes.Float64(float64(model.ObjectType_todo)),
				},
			}

			params := r.MakeRenderIconObjectParams(iconDetails, &IconObjectProps{Size: 30})
			iconTemplate := IconObjectTemplate(r, params)
			additionalTemplate = AdditionalIconTemplate(iconTemplate)
		}

		innerFlex = append(innerFlex, additionalTemplate, textComp)
	case model.BlockContentText_Quote:
		blockParams.Additional = AdditionalQuoteTemplate(color)
		blockParams.AdditionalClasses = append(blockParams.AdditionalClasses, "textColor-"+color)
		innerFlex = append(innerFlex, textComp)
	case model.BlockContentText_Checkbox:
		var checkboxComp templ.Component
		if blockText.Checked {
			checkboxComp = CheckboxCheckedTemplate()
			classes = append(classes, "isChecked")
		} else {
			checkboxComp = CheckboxUncheckedTemplate()
		}
		innerFlex = append(innerFlex, checkboxComp, textComp)
	default:
		innerFlex = append(innerFlex, textComp)
	}

	blockParams.Classes = append(blockParams.Classes, classes...)
	if len(innerFlex) != 0 {
		blockParams.Content = BlocksWrapper(&BlockWrapperParams{Classes: []string{"flex"}, Components: innerFlex})
	}
	blockParams.ContentClasses = append(blockParams.ContentClasses, contentClasses...)
	return blockParams

}
func (r *Renderer) RenderText(b *model.Block) templ.Component {
	params := r.makeTextBlockParams(b)
	return BlockTemplate(r, params)
}
