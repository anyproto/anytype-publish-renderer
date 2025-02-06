package renderer

import (
	"fmt"
	"strings"

	"github.com/gogo/protobuf/types"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"go.uber.org/zap"
)

type CoverResizeParams struct {
	CoverX            float64
	CoverY            float64
	CoverScale        float64
}

type CoverRenderParams struct {
	Id                string
	Src               string
	Classes           string
	CoverType         CoverType
	ResizeParams      CoverResizeParams
	UnsplashComponent templ.Component
	CoverTemplate     templ.Component
}

type CoverType int32

const (
	CoverType_Image         CoverType = 1
	CoverType_Color         CoverType = 2
	CoverType_Gradient      CoverType = 3
	CoverType_PrebuiltImage CoverType = 4
	CoverType_Source        CoverType = 5
)

func ToCoverType(val int64) (CoverType, error) {
	// TODO: cover type 0, no cover
	if val < 1 || val > 5 {
		return -1, fmt.Errorf("unknown cover type: %d", val)
	}

	return CoverType(val), nil
}

func (r *Renderer) MakeRenderPageCoverParams() (*CoverRenderParams, error) {
	fields := r.Sp.Snapshot.Data.GetDetails()
	return r.getCoverParams(fields, true, true)
}

func (r *Renderer) getCoverParams(fields *types.Struct, asImage bool, withAuthor bool) (*CoverRenderParams, error) {
	coverType, err := ToCoverType(pbtypes.GetInt64(fields, "coverType"))

	if err != nil {
		log.Warn("cover rendering failed", zap.Error(err))
		return nil, err
	}

	coverId := pbtypes.GetString(fields, "coverId")
	coverX := pbtypes.GetFloat64(fields, "coverX")
	coverY := pbtypes.GetFloat64(fields, "coverY")
	coverScale := pbtypes.GetFloat64(fields, "coverScale")
	class := fmt.Sprintf("type%d", coverType)

	params := &CoverRenderParams{
		Id:         coverId,
		CoverType:  coverType,
		Classes:    strings.Join([]string{class, coverId}, " "),
		ResizeParams: CoverResizeParams{
			CoverX:     coverX,
			CoverY:     coverY,
			CoverScale: coverScale,
		},
	}

	switch coverType {
	case CoverType_Image, CoverType_Source:
		src, err := r.getFileUrl(coverId)
		if err != nil {
			log.Warn("cover rendering failed", zap.Error(err))
			return nil, err
		}

		if withAuthor && coverType == CoverType_Source {
			author, authorUrl := r.getUnsplashDetails(coverId)
			if author != "" || authorUrl != "" {
				params.UnsplashComponent = UnsplashReferral(author, templ.SafeURL(authorUrl))
			}
		}

		params.Src = src
		if asImage {
			params.CoverTemplate = CoverImageTemplate(params)
		} else {
			params.CoverTemplate = CoverDefaultTemplate(params)
		}

		return params, nil

	case CoverType_Color, CoverType_Gradient:
		params.CoverTemplate = CoverDefaultTemplate(params)

		return params, nil
	}

	err = fmt.Errorf("unknown cover type: %d", int(coverType))
	log.Warn("cover rendering failed", zap.Error(err))
	return nil, err
}

func (r *Renderer) getUnsplashDetails(coverId string) (string, string) {
	fileSnapshot := r.getObjectSnapshot(coverId)
	fileDetails := fileSnapshot.GetSnapshot().GetData().GetDetails()
	author := pbtypes.GetString(fileDetails, "mediaArtistName")
	authorUrl := pbtypes.GetString(fileDetails, "mediaArtistURL")
	return author, authorUrl
}

func (r *Renderer) RenderPageCover() templ.Component {
	params, err := r.MakeRenderPageCoverParams()

	log.Warn("cover rendering failed: unknown cover type %+v", zap.Any("params", params))

	if err != nil {
		return NoneTemplate("")
	}

	switch params.CoverType {
	case CoverType_Image, CoverType_Source, CoverType_Color, CoverType_Gradient:
		return CoverBlockTemplate(r, params)
	default:
		log.Warn("cover rendering failed: unknown cover type", zap.Int("coverType", int(params.CoverType)))
		return NoneTemplate("")
	}
}
