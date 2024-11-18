package renderer

import (
	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"go.uber.org/zap"
)

type CoverRenderParams struct {
	Id      string
	Src     string
	Classes string
}

func (r *Renderer) MakeRenderPageCoverParams() (params *CoverRenderParams, err error) {
	fields := r.Sp.Snapshot.Data.GetDetails()
	coverId := pbtypes.GetString(fields, "coverId")
	src, err := r.AssetResolver.ByTargetObjectId(coverId)
	if err != nil {
		log.Warn("cover rendering failed", zap.Error(err))
		return
	}

	params = &CoverRenderParams{
		Id:      coverId,
		Src:     src,
		Classes: "type1",
	}

	return

}

func (r *Renderer) RenderPageCover() templ.Component {
	params, err := r.MakeRenderPageCoverParams()
	if err != nil {
		return NoneTemplate("")
	}

	return CoverTemplate(r, params)

}
