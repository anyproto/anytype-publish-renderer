package renderer

import (
	"fmt"
	"strconv"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"go.uber.org/zap"
)

type FileImageRenderParams struct {
	Id         string
	Src        string
	Classes    string
	ImageWidth string
}

func (r *Renderer) MakeRenderFileImageParams(b *model.Block) (params *FileImageRenderParams, err error) {
	file := b.GetFile()
	var src string
	src, err = r.AssetResolver.ByTargetObjectId(file.TargetObjectId)
	if err != nil {
		err = fmt.Errorf("file not found %s", file.TargetObjectId)
		return
	}

	align := "align" + strconv.Itoa(int(b.GetAlign()))

	width := pbtypes.GetFloat64(b.Fields, "width")
	log.Debug("image width", zap.Float64("width", width))
	imageWidth := strconv.Itoa(int(width*100)) + "%"

	params = &FileImageRenderParams{
		Id:         b.Id,
		Src:        src,
		Classes:    align,
		ImageWidth: imageWidth,
	}

	return
}

type FilePDFRenderParams struct {
	Id  string
	Src string
}

func (r *Renderer) MakeRenderFilePDFParams(b *model.Block) (params *FilePDFRenderParams, err error) {
	file := b.GetFile()
	var src string
	src, err = r.AssetResolver.ByTargetObjectId(file.TargetObjectId)
	if err != nil {
		err = fmt.Errorf("file not found %s", file.TargetObjectId)
		return
	}

	params = &FilePDFRenderParams{
		Id:  b.Id,
		Src: src,
	}

	return
}

type FileAudioRenderParams struct {
	Id  string
	Src string
}

func (r *Renderer) MakeRenderFileAudioParams(b *model.Block) (params *FileAudioRenderParams, err error) {
	file := b.GetFile()
	var src string
	src, err = r.AssetResolver.ByTargetObjectId(file.TargetObjectId)
	if err != nil {
		err = fmt.Errorf("file not found %s", file.TargetObjectId)
		return
	}

	params = &FileAudioRenderParams{
		Id:  b.Id,
		Src: src,
	}

	return
}

type FileVideoRenderParams struct {
	Id  string
	Src string
}

func (r *Renderer) MakeRenderFileVideoParams(b *model.Block) (params *FileVideoRenderParams, err error) {
	file := b.GetFile()
	var src string
	src, err = r.AssetResolver.ByTargetObjectId(file.TargetObjectId)
	if err != nil {
		err = fmt.Errorf("file not found %s", file.TargetObjectId)
		return
	}

	params = &FileVideoRenderParams{
		Id:  b.Id,
		Src: src,
	}

	return
}

func (r *Renderer) RenderFile(b *model.Block) templ.Component {
	file := b.GetFile()
	fileType := file.GetType()
	switch fileType {
	case model.BlockContentFile_Image:
		params, err := r.MakeRenderFileImageParams(b)
		if err != nil {
			return NoneTemplate(err.Error())
		}
		return FileImageTemplate(r, params)
	case model.BlockContentFile_PDF:
		params, err := r.MakeRenderFilePDFParams(b)
		if err != nil {
			return NoneTemplate(err.Error())
		}
		return FilePDFTemplate(r, params)
	case model.BlockContentFile_Audio:
		params, err := r.MakeRenderFileAudioParams(b)
		if err != nil {
			return NoneTemplate(err.Error())
		}
		return FileAudioTemplate(r, params)
	case model.BlockContentFile_Video:
		params, err := r.MakeRenderFileVideoParams(b)
		if err != nil {
			return NoneTemplate(err.Error())
		}
		return FileVideoTemplate(r, params)

	default:
		log.Warn("file type is not supported", zap.String("type", fileType.String()))
		err := fmt.Errorf("file type is not supported: %s", fileType.String())
		return NoneTemplate(err.Error())
	}
}
