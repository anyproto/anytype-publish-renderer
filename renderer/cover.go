package renderer

import (
	"fmt"

	"github.com/gogo/protobuf/types"

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
	return r.getCoverParams(fields)
}

func (r *Renderer) getCoverParams(fields *types.Struct) (*CoverRenderParams, error) {
	coverType, err := ToCoverType(pbtypes.GetInt64(fields, "coverType"))

	if err != nil {
		log.Warn("cover rendering failed", zap.Error(err))
		return nil, err
	}

	coverId := pbtypes.GetString(fields, "coverId")

	switch coverType {
	case CoverType_Image:
	case CoverType_Source:
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

	log.Warn("cover rendering failed: unknown cover type %+v", zap.Any("params", params))

	if err != nil {
		return EmptyCoverTemplate(bson.NewObjectId().Hex())
	}

	switch params.CoverType {
	case CoverType_Image:
	case CoverType_Source:
		return CoverImageTemplate(r, params)
	case CoverType_Color:
		return CoverColorTemplate(r, params)
	case CoverType_Gradient:
		return CoverGradientTemplate(r, params)

	}

	log.Warn("cover rendering failed: unknown cover type", zap.Int("coverType", int(params.CoverType)))
	return EmptyCoverTemplate(bson.NewObjectId().Hex())

}
