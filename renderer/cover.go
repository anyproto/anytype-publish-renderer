package renderer

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"github.com/globalsign/mgo/bson"
	"go.uber.org/zap"
)

type CoverRenderParams struct {
	Id        string
	Src       string
	Classes   string
	CoverType CoverType
}

type CoverType int32

const (
	CoverType_Image         CoverType = 1
	CoverType_Color         CoverType = 2
	CoverType_Gradient      CoverType = 3
	CoverType_PrebuiltImage CoverType = 4
)

func ToCoverType(val int64) (CoverType, error) {
	if val < 1 || val > 4 {
		return -1, fmt.Errorf("Unknown cover type: %d", val)
	}

	return CoverType(val), nil
}

func (r *Renderer) MakeRenderPageCoverParams() (*CoverRenderParams, error) {
	fields := r.Sp.Snapshot.Data.GetDetails()
	coverType, err := ToCoverType(pbtypes.GetInt64(fields, "coverType"))
	if err != nil {
		log.Warn("cover rendering failed", zap.Error(err))
		return nil, err
	}

	coverId := pbtypes.GetString(fields, "coverId")

	switch coverType {
	case CoverType_Image:
		src, err := r.getFileUrl(coverId)
		if err != nil {
			log.Warn("cover rendering failed", zap.Error(err))
			return nil, err
		}

		params := &CoverRenderParams{
			Id:        coverId,
			Src:       src,
			Classes:   "type1",
			CoverType: coverType,
		}

		return params, nil

	case CoverType_Color:
		color := pbtypes.GetString(fields, "coverId")
		params := &CoverRenderParams{
			Id:        coverId,
			Classes:   color,
			CoverType: coverType,
		}
		return params, nil

	case CoverType_Gradient:
		gradient := pbtypes.GetString(fields, "coverId")
		params := &CoverRenderParams{
			Id:        coverId,
			Classes:   gradient,
			CoverType: coverType,
		}
		return params, nil
	}

	err = fmt.Errorf("unknown cover type: %d", int(coverType))
	log.Warn("cover rendering failed", zap.Error(err))
	return nil, err
}

func (r *Renderer) RenderPageCover() templ.Component {
	params, err := r.MakeRenderPageCoverParams()
	if err != nil {
		return EmptyCoverTemplate(bson.NewObjectId().Hex())
	}

	switch params.CoverType {
	case CoverType_Image:
		return CoverImageTemplate(r, params)
	case CoverType_Color:
		return CoverColorTemplate(r, params)
	case CoverType_Gradient:
		return CoverGradientTemplate(r, params)

	}

	log.Warn("cover rendering failed: unknown cover type", zap.Int("coverType", int(params.CoverType)))
	return EmptyCoverTemplate(bson.NewObjectId().Hex())

}

type IconImageRenderParams struct {
	Id  string
	Src string
}

func (r *Renderer) MakeRenderPageIconImageParams() (params *IconImageRenderParams, err error) {
	fields := r.Sp.Snapshot.Data.GetDetails()
	iconEmoji := pbtypes.GetString(fields, "iconEmoji")
	if iconEmoji != "" {
		log.Debug("icon emoji", zap.String("id", iconEmoji))
		code := []rune(iconEmoji)[0]
		emojiSrc := r.GetEmojiUrl(code)
		params = &IconImageRenderParams{
			Id:  "emoji",
			Src: emojiSrc,
		}

		return
	}

	iconImageId := pbtypes.GetString(fields, "iconImage")
	src, err := r.getFileUrl(iconImageId)
	if err != nil {
		log.Warn("cover image rendering failed", zap.Error(err))
		return
	}

	params = &IconImageRenderParams{
		Id:  iconImageId,
		Src: src,
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
