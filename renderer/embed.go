package renderer

import (
	"fmt"
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

func (r *Renderer) MakeEmbedRenderParams(b *model.Block) *EmbedRenderParams {
	style := b.GetLatex().Processor.String()
	embedClass := "is" + style
	align := "align" + strconv.Itoa(int(b.GetAlign()))
	classes := []string{embedClass, align}

	if bgColor := b.GetBackgroundColor(); bgColor != "" {
		classes = append(classes, "bgColor", "bgColor-"+bgColor)
	}

	content := b.GetLatex().Text
	if b.GetLatex().Processor == model.BlockContentLatex_Mermaid {
		content = fmt.Sprintf(`<pre class="mermaid">%s</pre>`, content)
	}

	if b.GetLatex().Processor == model.BlockContentLatex_Kroki {
		content = fmt.Sprintf(`<img src="%s" />`, content)
	}

	if b.GetLatex().Processor == model.BlockContentLatex_Graphviz {
		content = fmt.Sprintf(`<pre class="graphviz-content">%s</pre>`, content)
	}

	return &EmbedRenderParams{
		Id:      b.Id,
		Classes: strings.Join(classes, " "),
		Content: content,
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
