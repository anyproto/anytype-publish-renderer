package renderer

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"go.uber.org/zap"
)

type EmbedRenderParams struct {
	Id      string
	Classes string
	Content string
}

type JsSVGString struct {
	Content string `json:"content,omitempty"`
}

func removeIframeWidthHeight(text string) string {
	if !strings.HasPrefix(text, "<iframe") {
		return text
	}

	r := regexp.MustCompile(`width="[0-9]*"|height="[0-9]*"`)
	text = r.ReplaceAllString(text, "")
	return text
}

func (r *Renderer) MakeEmbedRenderParams(b *model.Block) *EmbedRenderParams {
	latex := b.GetLatex()
	processor := latex.Processor
	text := latex.Text
	style := processor.String()
	bgColor := b.GetBackgroundColor()
	embedClass := "is" + style
	align := "align" + strconv.Itoa(int(b.GetAlign()))
	classes := []string{embedClass, align}

	if bgColor != "" {
		classes = append(classes, "bgColor", "bgColor-" + bgColor)
	}

	text = removeIframeWidthHeight(text)

	if processor == model.BlockContentLatex_Mermaid {
		text = fmt.Sprintf(`<pre class="mermaid">%s</pre>`, text)
	}

	if processor == model.BlockContentLatex_Kroki {
		text = fmt.Sprintf(`<img src="%s" />`, text)
	}

	if processor == model.BlockContentLatex_Graphviz {

		jsObj := JsSVGString{
			Content: text,
		}
		jsObjString, err := json.Marshal(jsObj)
		if err != nil {
			log.Error("svg json marshal error", zap.Error(err))
			text = fmt.Sprintf("<script>window.svgSrc['%s'] = `digraph { graphviz -> render error }`</script>", "block-" + b.Id)
		} else {
			text = fmt.Sprintf("<script>window.svgSrc['%s'] = %s</script>", "block-" + b.Id, string(jsObjString))
		}

	}

	return &EmbedRenderParams{
		Id:      b.Id,
		Classes: strings.Join(classes, " "),
		Content: text,
	}
}
func (r *Renderer) RenderEmbed(b *model.Block) templ.Component {

	switch b.GetLatex().Processor {
	case model.BlockContentLatex_Youtube:
		fallthrough
	case model.BlockContentLatex_Vimeo:
		fallthrough
	case model.BlockContentLatex_Soundcloud:
		fallthrough
	case model.BlockContentLatex_GoogleMaps:
		fallthrough
	case model.BlockContentLatex_Miro:
		fallthrough
	case model.BlockContentLatex_Figma:
		fallthrough
	case model.BlockContentLatex_Twitter:
		fallthrough
	case model.BlockContentLatex_OpenStreetMap:
		fallthrough
	case model.BlockContentLatex_Reddit:
		fallthrough
	case model.BlockContentLatex_Facebook:
		fallthrough
	case model.BlockContentLatex_Instagram:
		fallthrough
	case model.BlockContentLatex_Telegram:
		fallthrough
	case model.BlockContentLatex_GithubGist:
		fallthrough
	case model.BlockContentLatex_Codepen:
		fallthrough
	case model.BlockContentLatex_Latex:
		fallthrough
	case model.BlockContentLatex_Mermaid:
		fallthrough
	case model.BlockContentLatex_Bilibili:
		fallthrough
	case model.BlockContentLatex_Kroki:
		fallthrough
	case model.BlockContentLatex_Sketchfab:
		fallthrough
	case model.BlockContentLatex_Graphviz:
		params := r.MakeEmbedRenderParams(b)
		return EmbedTemplate(r, params)
	case model.BlockContentLatex_Chart:
	case model.BlockContentLatex_Excalidraw:
	case model.BlockContentLatex_Image:
	default:
	}

	log.Warn("embed block is not supported",
		zap.String("processor", b.GetLatex().Processor.String()),
		zap.String("id", b.Id))
	return NoneTemplate(fmt.Sprintf("embed block is not supported: %s", b.GetLatex().Processor.String()))

}
