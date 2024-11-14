package renderer

import (
	"fmt"
	"strconv"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"go.uber.org/zap"
)

type FileRenderParams struct {
	Type    model.BlockContentFileType
	Id      string
	Src     string
	Classes string
}

func (r *Renderer) MakeRenderFileParams(b *model.Block) (params *FileRenderParams, err error) {
	file := b.GetFile()
	fileType := file.GetType()
	switch fileType {
	case model.BlockContentFile_Image:
		var src string
		src, err = r.AssetResolver.ByTargetObjectId(file.TargetObjectId)
		if err != nil {
			log.Warn("file type is not supported", zap.String("type", fileType.String()), zap.Error(err))
			err = fmt.Errorf("file not found %s", file.TargetObjectId)
			return
		}

		align := "align" + strconv.Itoa(int(b.GetAlign()))
		params = &FileRenderParams{
			Type:    model.BlockContentFile_Image,
			Id:      b.Id,
			Src:     src,
			Classes: align,
		}
	default:
		log.Warn("file type is not supported", zap.String("type", fileType.String()))
		err = fmt.Errorf("file type is not supported: %s", fileType.String())
	}

	return
}

func (r *Renderer) RenderFile(b *model.Block) templ.Component {
	params, err := r.MakeRenderFileParams(b)
	if err != nil {
		return NoneTemplate(err.Error())
	}

	return FileImageTemplate(r, params)

}
