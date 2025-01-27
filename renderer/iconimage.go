package renderer

import (
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"go.uber.org/zap"
)

type IconImageRenderParams struct {
	Id      string
	Src     string
	Classes string
}

func (r *Renderer) MakeRenderPageIconImageParams() (params *IconImageRenderParams, err error) {
	fields := r.Sp.Snapshot.Data.GetDetails()
	iconEmoji := pbtypes.GetString(fields, "iconEmoji")

	// TODO: how to get layout align? hack around via title block now:
	titleBlock := r.BlocksById["title"]
	align := "align" + strconv.Itoa(int(titleBlock.GetAlign()))
	classes := []string{align}

	params = &IconImageRenderParams{
		Classes: strings.Join(classes, " "),
	}
	// TODO: support isHuman for profile icon
	if iconEmoji != "" {
		log.Debug("icon emoji", zap.String("id", iconEmoji))
		code := []rune(iconEmoji)[0]
		emojiSrc := r.GetEmojiUrl(code)
		params.Src = emojiSrc

		return
	}

	iconImageId := pbtypes.GetString(fields, "iconImage")
	src, err := r.getFileUrl(iconImageId)
	if err != nil {
		log.Warn("cover image rendering failed", zap.Error(err))
		return
	}

	params.Src = src
	return

}

func (r *Renderer) RenderPageIconImage() templ.Component {
	params, err := r.MakeRenderPageIconImageParams()
	if err != nil {
		return NoneTemplate("")
	}

	return IconImageTemplate(r, params)
}
