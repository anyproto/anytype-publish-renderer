package renderer

import (
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
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

	layoutType := model.ObjectTypeLayout(pbtypes.GetInt64(fields, "layout"))
	layoutAlign := pbtypes.GetInt64(fields, "layoutAlign")
	align := "align" + strconv.Itoa(int(layoutAlign))
	classes := []string{align}

	params = &IconImageRenderParams{}

	switch layoutType {
	case model.ObjectType_basic:
		// if basic, then we get emoji or image, if emoji is not set
		iconEmoji := pbtypes.GetString(fields, "iconEmoji")
		if iconEmoji != "" {
			log.Debug("icon emoji", zap.String("id", iconEmoji))
			code := []rune(iconEmoji)[0]
			emojiSrc := r.GetEmojiUrl(code)
			params = &IconImageRenderParams{
				Classes: strings.Join(classes, " "),
				Src:     emojiSrc,
			}
			return
		} else {
			iconImageId := pbtypes.GetString(fields, "iconImage")
			src, errIcon := r.getFileUrl(iconImageId)
			if errIcon != nil {
				log.Warn("cover image rendering failed", zap.Error(err))
				// TODO: just don't show an icon here? check in client
				return nil, errIcon
			}
			params = &IconImageRenderParams{
				Classes: strings.Join(classes, " "),
				Src:     src,
			}
			return
		}
	case model.ObjectType_profile:
		// if profile, then we get image or userSvg, and add isHuman class
		iconImageId := pbtypes.GetString(fields, "iconImage")
		if iconImageId == "" {
			// TODO: no image user svg: https://github.com/anyproto/anytype-ts/blob/main/src/ts/component/util/iconObject.tsx#L269
			return params, nil
		}
		src, errIcon := r.getFileUrl(iconImageId)

		if errIcon != nil {
			log.Warn("cover image rendering failed", zap.Error(err))
			return nil, errIcon
		}
		classes = append(classes, "isHuman")
		params = &IconImageRenderParams{
			Classes: strings.Join(classes, " "),
			Src:     src,
		}

		return params, nil

	}

	return

}

func (r *Renderer) RenderPageIconImage() templ.Component {
	params, err := r.MakeRenderPageIconImageParams()
	if err != nil {
		return NoneTemplate("")
	}

	return IconImageTemplate(r, params)
}
