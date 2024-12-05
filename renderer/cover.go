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
	src, err := r.getFileUrl(coverId)
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

type IconImageRenderParams struct {
	Id  string
	Src string
}

func (r *Renderer) MakeRenderPageIconImageParams() (params *IconImageRenderParams, err error) {
	fields := r.Sp.Snapshot.Data.GetDetails()
	iconImageId := pbtypes.GetString(fields, "iconImage")
	log.Debug("-- icon image", zap.String("id", iconImageId))
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
