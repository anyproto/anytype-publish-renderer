package renderer

import (
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"go.uber.org/zap"
)

type IconImageRenderParams struct {
	Id          string
	Src         string
	Classes     string
	IconClasses string
}

func (r *Renderer) MakeRenderPageIconImageParams() (params *IconImageRenderParams, err error) {
	fields := r.Sp.Snapshot.Data.GetDetails()

	layoutType := model.ObjectTypeLayout(pbtypes.GetInt64(fields, "layout"))
	layoutAlign := pbtypes.GetInt64(fields, "layoutAlign")
	align := "align" + strconv.Itoa(int(layoutAlign))
	classes := []string{align}

	// TODO: refactoring GO-4950/support-all-layouts
	switch layoutType {
	case model.ObjectType_basic:
		// if basic, then we get emoji or image, if emoji is not set
		iconEmoji := pbtypes.GetString(fields, "iconEmoji")
		iconClasses := []string{"c96"}
		if iconEmoji != "" {
			log.Debug("icon emoji", zap.String("id", iconEmoji))
			code := []rune(iconEmoji)[0]
			emojiSrc := r.GetEmojiUrl(code)

			params = &IconImageRenderParams{
				IconClasses: strings.Join(iconClasses, " "),
				Classes:     strings.Join(classes, " "),
				Src:         emojiSrc,
			}
			return params, nil
		} else {
			iconImageId := pbtypes.GetString(fields, "iconImage")
			src, errIcon := r.getFileUrl(iconImageId)
			if errIcon != nil {
				log.Warn("cover image rendering failed", zap.Error(err))
				params = &IconImageRenderParams{
					Classes: strings.Join(classes, " "),
				}
				return params, nil
			}
			params = &IconImageRenderParams{
				IconClasses: strings.Join(iconClasses, " "),
				Classes:     strings.Join(classes, " "),
				Src:         src,
			}
			return
		}
	case model.ObjectType_profile:
		// if profile, then we get image or userSvg, and add isHuman class
		iconImageId := pbtypes.GetString(fields, "iconImage")
		iconClasses := []string{"c128", "isHuman"}
		params = &IconImageRenderParams{
			IconClasses: strings.Join(iconClasses, " "),
			Classes:     strings.Join(classes, " "),
		}

		if iconImageId == "" {
			// TODO: no image user svg: https://github.com/anyproto/anytype-ts/blob/main/src/ts/component/util/iconObject.tsx#L269
			return params, nil
		}
		src, errIcon := r.getFileUrl(iconImageId)

		if errIcon != nil {
			log.Warn("profile image rendering failed", zap.String("iconImageId", iconImageId), zap.Error(err))
			return nil, errIcon
		}

		params.Src = src
		return params, nil
	default:
		// TODО: GO-4950/support-all-layouts
		params = &IconImageRenderParams{
			Classes: strings.Join(classes, " "),
		}

		return

	}

}

func isHumanLayout(layout model.ObjectTypeLayout) bool {
	return layout == model.ObjectType_profile || layout == model.ObjectType_participant
}

func pageIconInitSize(layout model.ObjectTypeLayout) int32 {
	if isHumanLayout(layout) {
		return 128
	} else {
		return 96
	}
}
func (r *Renderer) RenderPageIconImage() templ.Component {
	details := r.Sp.Snapshot.Data.GetDetails()
	layout := getRelationField(details, bundle.RelationKeyLayout, relationToObjectTypeLayout)

	props := &IconObjectProps{
		Size: pageIconInitSize(layout),
	}
	params := r.MakeRenderIconObjectParams(details, props)
	content := IconObjectTemplate(r, params)

	classes := []string{""}
	var blockType string
	if layout == model.ObjectType_profile {
		blockType = "IconUser"
		classes = append(classes, "isHuman")
	} else if layout == model.ObjectType_basic {
		blockType = "IconPage"
	}
	blockParams := &BlockParams{
		BlockType: blockType,
		Classes:   classes,
		Content:   content,
	}
	return BlockTemplate(r, blockParams)
}
