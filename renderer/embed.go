package renderer

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
)

type EmbedIframeData struct {
	AllowIframeResize bool
	InsertBeforeLoad  bool
	UseRootHeight     bool
	Align             model.BlockAlign
	Processor         model.BlockContentLatexProcessor
	ClassName         string
	BlockId           string
	Js                string
	Html              string
}

type EmbedRenderParams struct {
	Id       string
	Classes  string
	Content  string
	Data     EmbedIframeData
	IsIframe bool
	Sandbox  string
}

type JsSVGString struct {
	Content string `json:"content,omitempty"`
}

var iframeParams = `frameborder="0" scrolling="no" allowfullscreen`
var domains = map[model.BlockContentLatexProcessor][]string{
	model.BlockContentLatex_Youtube:       {`youtube\.com`, `youtu\.be`},
	model.BlockContentLatex_Vimeo:         {`vimeo\.com`},
	model.BlockContentLatex_GoogleMaps:    {`google\.[^/]+/maps`},
	model.BlockContentLatex_Miro:          {`miro\.com`},
	model.BlockContentLatex_Figma:         {`figma\.com`},
	model.BlockContentLatex_OpenStreetMap: {`openstreetmap\.org/\#map`},
	model.BlockContentLatex_Telegram:      {`t\.me`},
	model.BlockContentLatex_Codepen:       {`codepen\.io`},
	model.BlockContentLatex_Bilibili:      {`bilibili\.com`, `b23\.tv`},
	model.BlockContentLatex_Kroki:         {`kroki\.io`},
	model.BlockContentLatex_GithubGist:    {`gist\.github\.com`},
	model.BlockContentLatex_Sketchfab:     {`sketchfab\.com`},
}

func (r *Renderer) MakeEmbedRenderParams(b *model.Block) *EmbedRenderParams {
	id := b.GetId()
	latex := b.GetLatex()
	processor := latex.Processor
	text := latex.Text
	style := processor.String()
	bgColor := b.GetBackgroundColor()
	embedClass := "is" + style
	align := b.GetAlign()
	classes := []string{embedClass, fmt.Sprintf("align%d", align)}
	data := EmbedIframeData{}
	isIframe := false
	sandbox := []string{}

	if bgColor != "" {
		classes = append(classes, "bgColor", "bgColor-"+bgColor)
	}

	switch processor {
	default:
		isIframe = true
		allowIframeResize := allowIframeResize(processor)
		allowScript := false
		sandbox = append(sandbox, "allow-scripts", "allow-same-origin", "allow-popups")

		if allowPresentation(processor) {
			sandbox = append(sandbox, "allow-presentation")
		}

		if allowPopup(processor) {
			sandbox = append(sandbox, "allow-popups")
		}

		data.AllowIframeResize = allowIframeResize
		data.InsertBeforeLoad = insertBeforeLoad(processor)
		data.UseRootHeight = useRootHeight(processor)
		data.Align = align
		data.Processor = processor
		data.ClassName = embedClass
		data.BlockId = id

		// Fix Bilibili schemeless URLs and autoplay
		if processor == model.BlockContentLatex_Bilibili {
			reSrc := regexp.MustCompile(`src="(//player[^"]+)"`)
			text = reSrc.ReplaceAllString(text, `src="https:$1"`)

			reAutoplay := regexp.MustCompile(`autoplay=`)
			if !reAutoplay.MatchString(text) {
				reInsertAutoplay := regexp.MustCompile(`(src="[^"]+)`)
				text = reInsertAutoplay.ReplaceAllString(text, `$1&autoplay=0`)
			}
		}

		// Convert Kroki code into an SVG URL
		if processor == model.BlockContentLatex_Kroki && !strings.HasPrefix(text, "https://kroki.io") {
			compressed, err := compressAndEncode(text)

			if err == nil {
				typeId := pbtypes.GetString(b.GetFields(), "type")
				text = fmt.Sprintf("https://kroki.io/%s/svg/%s", typeId, compressed)
			}
		}

		// Process embedded content
		if allowEmbedUrl(processor) && !regexp.MustCompile(`<iframe|script`).MatchString(text) {
			text = getHtml(processor, getParsedUrl(text))
		}

		// Sketchfab embed handling
		if processor == model.BlockContentLatex_Sketchfab && regexp.MustCompile(`<iframe|script`).MatchString(text) {
			iframeMatch := regexp.MustCompile(`<iframe.*?</iframe>`).FindString(text)
			if iframeMatch != "" {
				text = iframeMatch
			}
		}

		// Check if script tags should be allowed
		if processor == model.BlockContentLatex_GithubGist {
			allowScript = true
		}

		// Telegram embed handling
		if processor == model.BlockContentLatex_Telegram {
			allowScript = true
		}

		if !allowScript {
			text = regexp.MustCompile(`<script`).ReplaceAllString(text, "&lt;script")
		}

		// Update sanitization parameters
		if allowJs(processor) {
			data.Js = text
		} else {
			data.Html = text
		}

	case model.BlockContentLatex_Latex:
		break

	case model.BlockContentLatex_Mermaid:
		text = fmt.Sprintf(`<div class="mermaidChart">%s</div>`, text)

	case model.BlockContentLatex_Graphviz:
		break
	}

	return &EmbedRenderParams{
		Id:       b.Id,
		Classes:  strings.Join(classes, " "),
		Content:  text,
		Data:     data,
		IsIframe: isIframe,
		Sandbox:  strings.Join(sandbox, " "),
	}
}

func (r *Renderer) RenderEmbed(b *model.Block) templ.Component {

	params := r.MakeEmbedRenderParams(b)
	return EmbedTemplate(r, params)
}

func getProcessorByUrl(inputUrl string) *model.BlockContentLatexProcessor {
	parsedUrl, err := url.Parse(inputUrl)
	if err != nil {
		return nil
	}

	for processor, patterns := range domains {
		for _, pattern := range patterns {
			reg := regexp.MustCompile(`(?i)://([^.]*.)?` + pattern)

			if reg.MatchString(inputUrl) {
				// Restrict YouTube channel links
				if processor == model.BlockContentLatex_Youtube && strings.HasPrefix(parsedUrl.Path, "/@") {
					return nil
				}
				return &processor
			}
		}
	}
	return nil
}

func getParsedUrl(inputUrl string) string {
	processor := getProcessorByUrl(inputUrl)
	if processor == nil {
		return inputUrl
	}

	switch *processor {
	case model.BlockContentLatex_Youtube:
		return fmt.Sprintf("https://www.youtube.com/embed/%s", extractYoutubeId(inputUrl))

	case model.BlockContentLatex_Vimeo:
		if parsed, err := url.Parse(inputUrl); err == nil {
			return fmt.Sprintf("https://player.vimeo.com/video%s", parsed.Path)
		}
	case model.BlockContentLatex_GoogleMaps:
		return parseGoogleMapsUrl(inputUrl)
	case model.BlockContentLatex_Miro:
		return strings.Split(inputUrl, "?")[0] + "/live-embed/"
	case model.BlockContentLatex_Figma:
		return fmt.Sprintf("https://www.figma.com/embed?embed_host=share&url=%s", url.QueryEscape(inputUrl))
	case model.BlockContentLatex_OpenStreetMap:
		return parseOpenStreetMapUrl(inputUrl)
	case model.BlockContentLatex_Bilibili:
		return parseBilibiliUrl(inputUrl)
	case model.BlockContentLatex_Sketchfab:
		return fmt.Sprintf("https://sketchfab.com/models/%s/embed", extractSketchfabId(inputUrl))
	case model.BlockContentLatex_GithubGist:
		return strings.Split(inputUrl, "#")[0]
	}

	return inputUrl
}

func extractYoutubeId(url string) string {
	// Regex for detecting YouTube Shorts URLs
	shortsReg := regexp.MustCompile(`/shorts/`)
	url = shortsReg.ReplaceAllString(url, "/watch?v=")

	// Regex to extract video ID
	pm := regexp.MustCompile(`^.*(youtu\.be/|v/|u/\w/|embed/|watch\?v=|&v=)([^#&?]*).*`).FindStringSubmatch(url)
	// Regex to extract timestamp parameter
	tm := regexp.MustCompile(`(\?t=|&t=)(\d+)`).FindStringSubmatch(url)

	// Check if video ID was found
	if len(pm) < 3 || len(pm[2]) == 0 {
		return ""
	}

	// Build the final result
	id := pm[2]
	if len(tm) > 2 && len(tm[2]) > 0 {
		return id + "?start=" + tm[2]
	}

	return id
}

func parseGoogleMapsUrl(inputUrl string) string {
	if match := regexp.MustCompile(`place/([^/]+)`).FindStringSubmatch(inputUrl); match != nil {
		return fmt.Sprintf("https://www.google.com/maps/embed/v1/place?key=%s&q=%s", GOOGLE_MAPS, url.QueryEscape(match[1]))
	}
	return inputUrl
}

func parseOpenStreetMapUrl(inputUrl string) string {
	if match := regexp.MustCompile(`#map=([-0-9\./]+)`).FindStringSubmatch(inputUrl); match != nil {
		parts := strings.Split(match[1], "/")
		if len(parts) >= 3 {
			return fmt.Sprintf("https://www.openstreetmap.org/export/embed.html?bbox=%s&layer=mapnik", url.QueryEscape(strings.Join([]string{parts[2], parts[1], parts[2], parts[1]}, ",")))
		}
	}
	return inputUrl
}

func parseBilibiliUrl(inputUrl string) string {
	parsed, err := url.Parse(inputUrl)
	if err != nil {
		return inputUrl
	}
	parts := strings.Split(parsed.Path, "/")
	if len(parts) < 3 {
		return inputUrl
	}
	bvid := parts[2]
	return fmt.Sprintf("https://player.bilibili.com/player.html?bvid=%s&high_quality=1&autoplay=0", bvid)
}

func extractSketchfabId(inputUrl string) string {
	parts := strings.Split(inputUrl, "/")
	if len(parts) == 0 {
		return ""
	}
	nameParts := strings.Split(parts[len(parts)-1], "-")
	if len(nameParts) == 0 {
		return ""
	}
	return nameParts[len(nameParts)-1]
}

func getHtml(processor model.BlockContentLatexProcessor, content string) string {
	fnName := fmt.Sprintf("Get%sHtml", processor.String())
	method := reflect.ValueOf(&embedFunctions).MethodByName(fnName)

	if method.IsValid() {
		results := method.Call([]reflect.Value{reflect.ValueOf(content)})
		return results[0].String()
	}
	return content
}

type embedFuncStruct struct{}

var embedFunctions = embedFuncStruct{}

func (e *embedFuncStruct) GetYoutubeHtml(content string) string {
	parsedUrl, err := url.Parse(content)
	if err != nil {
		return content
	}
	q := parsedUrl.Query()
	q.Set("enablejsapi", "1")
	q.Set("rel", "0")
	parsedUrl.RawQuery = q.Encode()

	return fmt.Sprintf(`<iframe id="player" src="%s" %s title="YouTube video player"></iframe>`, parsedUrl.String(), iframeParams)
}

func (e *embedFuncStruct) GetVimeoHtml(content string) string {
	return fmt.Sprintf(`<iframe src="%s" %s allow="autoplay; fullscreen; picture-in-picture"></iframe>`, content, iframeParams)
}

func (e *embedFuncStruct) GetGoogleMapsHtml(content string) string {
	return fmt.Sprintf(`<iframe src="%s" %s loading="lazy"></iframe>`, content, iframeParams)
}

func (e *embedFuncStruct) GetMiroHtml(content string) string {
	return fmt.Sprintf(`<iframe src="%s" %s allow="fullscreen; clipboard-read; clipboard-write"></iframe>`, content, iframeParams)
}

func (e *embedFuncStruct) GetFigmaHtml(content string) string {
	return fmt.Sprintf(`<iframe src="%s" %s></iframe>`, content, iframeParams)
}

func (e *embedFuncStruct) GetOpenStreetMapHtml(content string) string {
	return fmt.Sprintf(`<iframe src="%s" %s></iframe>`, content, iframeParams)
}

func (e *embedFuncStruct) GetGithubGistHtml(content string) string {
	return fmt.Sprintf(`<script src="%s.js"></script>`, content)
}

func (e *embedFuncStruct) GetCodepenHtml(content string) string {
	parsedUrl, err := url.Parse(content)
	if err != nil {
		return ""
	}
	parts := strings.Split(parsedUrl.Path, "/")
	if len(parts) < 4 {
		return ""
	}
	return fmt.Sprintf(`<p class="codepen" data-height="300" data-default-tab="html,result" data-slug-hash="%s" data-user="%s"></p>`, parts[3], parts[1])
}

func (e *embedFuncStruct) GetBilibiliHtml(content string) string {
	return fmt.Sprintf(`<iframe src="%s" %s></iframe>`, content, iframeParams)
}

func (e *embedFuncStruct) GetSketchfabHtml(content string) string {
	return fmt.Sprintf(`<iframe src="%s" %s></iframe>`, content, iframeParams)
}

func (e *embedFuncStruct) GetImageHtml(content string) string {
	return fmt.Sprintf(`<img src="%s" />`, content)
}

func allowPresentation(p model.BlockContentLatexProcessor) bool {
	allowed := map[model.BlockContentLatexProcessor]bool{
		model.BlockContentLatex_Youtube: true, model.BlockContentLatex_Vimeo: true, model.BlockContentLatex_Bilibili: true,
	}
	return allowed[p]
}

func allowEmbedUrl(p model.BlockContentLatexProcessor) bool {
	allowed := map[model.BlockContentLatexProcessor]bool{
		model.BlockContentLatex_Youtube: true, model.BlockContentLatex_Vimeo: true, model.BlockContentLatex_GoogleMaps: true,
		model.BlockContentLatex_Miro: true, model.BlockContentLatex_Figma: true, model.BlockContentLatex_OpenStreetMap: true,
		model.BlockContentLatex_Telegram: true, model.BlockContentLatex_GithubGist: true, model.BlockContentLatex_Codepen: true,
		model.BlockContentLatex_Bilibili: true, model.BlockContentLatex_Kroki: true, model.BlockContentLatex_Sketchfab: true,
		model.BlockContentLatex_Image: true,
	}
	return allowed[p]
}

func allowJs(p model.BlockContentLatexProcessor) bool {
	allowed := map[model.BlockContentLatexProcessor]bool{
		model.BlockContentLatex_Chart: true,
	}
	return allowed[p]
}

func allowPopup(p model.BlockContentLatexProcessor) bool {
	allowed := map[model.BlockContentLatexProcessor]bool{
		model.BlockContentLatex_Bilibili: true,
	}
	return allowed[p]
}

func allowIframeResize(p model.BlockContentLatexProcessor) bool {
	allowed := map[model.BlockContentLatexProcessor]bool{
		model.BlockContentLatex_Twitter: true, model.BlockContentLatex_Reddit: true, model.BlockContentLatex_Facebook: true,
		model.BlockContentLatex_Instagram: true, model.BlockContentLatex_Telegram: true, model.BlockContentLatex_GithubGist: true,
		model.BlockContentLatex_Codepen: true, model.BlockContentLatex_Kroki: true, model.BlockContentLatex_Chart: true,
		model.BlockContentLatex_Image: true,
	}
	return allowed[p]
}

func insertBeforeLoad(p model.BlockContentLatexProcessor) bool {
	allowed := map[model.BlockContentLatexProcessor]bool{
		model.BlockContentLatex_Twitter: true, model.BlockContentLatex_Reddit: true, model.BlockContentLatex_Instagram: true,
		model.BlockContentLatex_Codepen: true,
	}
	return allowed[p]
}

func useRootHeight(p model.BlockContentLatexProcessor) bool {
	allowed := map[model.BlockContentLatexProcessor]bool{
		model.BlockContentLatex_Twitter: true, model.BlockContentLatex_Telegram: true, model.BlockContentLatex_Instagram: true,
		model.BlockContentLatex_GithubGist: true, model.BlockContentLatex_Codepen: true, model.BlockContentLatex_Kroki: true,
		model.BlockContentLatex_Chart: true,
	}
	return allowed[p]
}

func compressAndEncode(text string) (string, error) {
	var buf bytes.Buffer

	writer := zlib.NewWriter(&buf)
	_, err := writer.Write([]byte(text))
	if err != nil {
		return "", err
	}
	writer.Close()

	// Encode to base64 and replace characters for URL safety
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	encoded = strings.ReplaceAll(encoded, "+", "-")
	encoded = strings.ReplaceAll(encoded, "/", "_")

	return encoded, nil
}
