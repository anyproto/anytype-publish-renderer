package renderer

import (
	"cmp"
	"fmt"
	"html"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode/utf16"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"github.com/anyproto/anytype-publish-renderer/renderer/markintervaltree"
	"github.com/anyproto/anytype-publish-renderer/utils"
	"github.com/gogo/protobuf/types"
	"go.uber.org/zap"
)

type TextRenderParams struct {
	Classes        string
	ContentClasses string
	Id             string
	InnerFlex      []templ.Component
	OuterFlex      []templ.Component
	ChildrenIds    []string
}

func cmpMarks(a, b *model.BlockContentTextMark) int {
	return cmp.Compare(a.Range.From, b.Range.From)
}

func emojiParam(t model.BlockContentTextStyle) int32 {
	switch t {
	default:
		return 20
	case model.BlockContentText_Header1:
		return 30
	case model.BlockContentText_Header2:
		return 26
	case model.BlockContentText_Header3:
		return 22
	}
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

		var iconHtml, class, link string

		if details != nil && len(details.Fields) != 0 {
			params := r.MakeRenderIconObjectParams(details, &IconObjectProps{Size: emojiSize})

			var err error
			iconHtml, err = utils.TemplToString(IconObjectTemplate(r, params))

			if err != nil {
				log.Error("Failed to render mention icon", zap.Error(err))
			}

			if iconHtml != "" {
				class = "withImage"
			}
			spaceId := getRelationField(details, bundle.RelationKeySpaceId, relationToString)
			link = fmt.Sprintf(linkTemplate, mark.Param, spaceId)
		}

		return `<a href=` + link + ` target="_blank" class="markupmention ` + class + `"><span class="smile">` + iconHtml + `</span><img src="./static/img/space.svg" class="space" /><span class="name">` + s + `</span></a>`

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
		spaceId := getRelationField(details, bundle.RelationKeySpaceId, relationToString)
		link := fmt.Sprintf(linkTemplate, mark.Param, spaceId)
		return fmt.Sprintf(`<a href="%s" class="markuplink" target="_blank">`, link) + s + "</a>"
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
func (r *Renderer) applyNonOverlapingMarks(style model.BlockContentTextStyle, text string, marks []*model.BlockContentTextMark) string {
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
			markedPart = r.applyMark(style, markedPart, m)
			log.Debug("apply mark", zap.String("markedPart", markedPart), zap.Int32("from", m.Range.From), zap.Int32("to", m.Range.To))
		}
		log.Debug("final marked part", zap.String("m", markedPart))
		markedText.WriteString(markedPart)
	}

	return markedText.String()
}

func replaceNewlineBr(text string) string {
	r := regexp.MustCompile(`\r?\n`)
	text = r.ReplaceAllString(text, "<br>")
	return text
}

func (r *Renderer) MakeRenderTextParams(b *model.Block) (params *TextRenderParams) {
	blockText := b.GetText()
	style := blockText.GetStyle()
	bgColor := b.GetBackgroundColor()
	color := blockText.GetColor()
	iconEmoji := blockText.GetIconEmoji()
	iconImage := blockText.GetIconImage()
	classes := []string{"block", "blockText"}
	contentClasses := []string{"content"}

	classes = append(classes, "text"+style.String())
	classes = append(classes, "align"+strconv.Itoa(int(b.GetAlign())))

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
		marks := blockText.GetMarks().Marks
		text = r.applyNonOverlapingMarks(style, text, marks)
		text = replaceNewlineBr(text)
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
		externalComp := ToggleMarkerTemplate(utils.GetColor(color))
		innerFlex = append(innerFlex, externalComp, textComp)
	case model.BlockContentText_Numbered:
		number := r.BlockNumbers[b.Id]
		log.Debug("number", zap.Int("num", number), zap.String("id", b.Id))
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
				bundle.RelationKeyIconEmoji.String(): pbtypes.String(iconEmoji),
				bundle.RelationKeyIconImage.String(): pbtypes.String(iconImage),
				bundle.RelationKeyLayout.String():    pbtypes.Float64(float64(model.ObjectType_basic)),
			},
		}

		params := r.MakeRenderIconObjectParams(details, &IconObjectProps{Size: 20})
		iconTemplate := IconObjectTemplate(r, params)
		additionalTemplate := AdditionalIconTemplate(iconTemplate)

		innerFlex = append(innerFlex, additionalTemplate, textComp)
	case model.BlockContentText_Title:
		details := r.Sp.Snapshot.Data.GetDetails()
		layout := getRelationField(details, bundle.RelationKeyLayout, relationToObjectTypeLayout)
		done := getRelationField(details, bundle.RelationKeyDone, relationToBool)
		additionalTemplate := NoneTemplate("")

		if isTodoLayout(layout) {
			iconDetails := &types.Struct{
				Fields: map[string]*types.Value{
					bundle.RelationKeyDone.String():   pbtypes.Bool(done),
					bundle.RelationKeyLayout.String(): pbtypes.Float64(float64(model.ObjectType_todo)),
				},
			}

			params := r.MakeRenderIconObjectParams(iconDetails, &IconObjectProps{Size: 30})
			iconTemplate := IconObjectTemplate(r, params)
			additionalTemplate = AdditionalIconTemplate(iconTemplate)
		}

		innerFlex = append(innerFlex, additionalTemplate, textComp)
	case model.BlockContentText_Quote:
		externalComp := AdditionalQuoteTemplate(color)
		outerFlex = append(outerFlex, externalComp)
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

	params = &TextRenderParams{
		Id:             b.Id,
		Classes:        strings.Join(classes, " "),
		ContentClasses: strings.Join(contentClasses, " "),
		ChildrenIds:    b.ChildrenIds,
		OuterFlex:      outerFlex,
		InnerFlex:      innerFlex,
	}
	return

}
func (r *Renderer) RenderText(b *model.Block) templ.Component {
	params := r.MakeRenderTextParams(b)
	return TextTemplate(r, params)
}
