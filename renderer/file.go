package renderer

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"go.uber.org/zap"
)

type ImageRenderParams struct {
	Id  string
	Src string
}

func (r *Renderer) RenderFile(b *model.Block) templ.Component {
	file := b.GetFile()
	fileType := file.GetType()
	switch fileType {
	case model.BlockContentFile_Image:
		src, err := r.AssetResolver.ByTargetObjectId(file.TargetObjectId)
		if err != nil {
			log.Warn("file type is not supported", zap.String("type", fileType.String()))
			return NoneTemplate(fmt.Sprintf("file not found %s", file.TargetObjectId))
		}

		params := &ImageRenderParams{
			Id:  b.Id,
			Src: src,
		}
		return FileImageTemplate(r, params)
	default:
		log.Warn("file type is not supported", zap.String("type", fileType.String()))
		return NoneTemplate(fmt.Sprintf("file type is not supported: %s", fileType.String()))
	}
}
