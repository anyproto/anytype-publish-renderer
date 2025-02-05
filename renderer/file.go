package renderer

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"github.com/gogo/protobuf/types"
	"go.uber.org/zap"
)

type FileFileRenderParams struct {
	Id   string
	Src  templ.SafeURL
	Name string
	Size string
}

type FileMediaRenderParams struct {
	Id         string
	Src        string
	Classes    string
	Width      string
}

type FilePdfRenderParams struct {
	Id         string
	Src        string
	Classes    string
	Width      string
	Name	   string
	Size	   string
}

func (r *Renderer) getFileUrl(id string) (url string, err error) {
	path := fmt.Sprintf("filesObjects/%s.pb", id)
	snapshot, err := r.ReadJsonpbSnapshot(path)
	if err != nil {
		return
	}

	if snapshot.SbType != model.SmartBlockType_FileObject {
		err = fmt.Errorf("snaphot %s is not FileObjects, %d", path, snapshot.SbType)
		return
	}

	fields := snapshot.Snapshot.Data.GetDetails()
	source := pbtypes.GetString(fields, "source")
	if source == "" {
		err = fmt.Errorf("FileObject %s 'source' is empty", id)
		return
	}

	// fixes GO-4975
	source = strings.ReplaceAll(source, `\`, "%5C")
	url = fmt.Sprintf("%s/%s", r.Config.PublishFilesPath, source)

	return

}

func (r *Renderer) getFileBlock(id string) (block *model.BlockContentFile, err error) {
	path := fmt.Sprintf("filesObjects/%s.pb", id)
	var (
		jsonPbSnapshot string
		ok             bool
	)
	if jsonPbSnapshot, ok = r.UberSp.PbFiles[path]; !ok {
		return nil, fmt.Errorf("file %s not exists", id)
	}
	snapshot, err := readJsonpbSnapshot(jsonPbSnapshot)
	if err != nil {
		return
	}

	if snapshot.SbType != model.SmartBlockType_FileObject {
		err = fmt.Errorf("snaphot %s is not FileObjects, %d", path, snapshot.SbType)
		return
	}

	blocks := snapshot.GetSnapshot().GetData().GetBlocks()
	for _, bl := range blocks {
		if bl.GetFile() == nil {
			continue
		}
		return bl.GetFile(), nil
	}
	return
}

func GetWidth(fields *types.Struct) string {
	width := pbtypes.GetFloat64(fields, "width")
	log.Debug("image width", zap.Float64("width", width))

	if int(width * 100) != 0 {
		return strconv.Itoa(int(width*100)) + "%"
	}
	return ""
}

func GetAlign(align model.BlockAlign) string {
	return "align" + strconv.Itoa(int(align));
}

func (r *Renderer) MakeRenderFileImageParams(b *model.Block) (params *FileMediaRenderParams, err error) {
	file := b.GetFile()
	
	var src string
	src, err = r.getFileUrl(file.TargetObjectId)
	if err != nil {
		err = fmt.Errorf("file not found %s", file.TargetObjectId)
		return
	}

	align := "align" + strconv.Itoa(int(b.GetAlign()))

	params = &FileMediaRenderParams{
		Id:         b.Id,
		Src:        src,
		Classes:    align,
		Width:		GetWidth(b.Fields),
	}

	return
}

func (r *Renderer) MakeRenderFilePDFParams(b *model.Block) (params *FilePdfRenderParams, err error) {
	file := b.GetFile()
	var src string
	src, err = r.getFileUrl(file.TargetObjectId)
	if err != nil {
		err = fmt.Errorf("file not found %s", file.TargetObjectId)
		return
	}

	name := file.Name
	size := prettyByteSize(file.Size_)

	params = &FilePdfRenderParams{
		Id:         b.Id,
		Src:        src,
		Classes:    GetAlign(b.GetAlign()),
		Width:		GetWidth(b.Fields),
		Name:		name,
		Size:		size,
	}

	return
}

func (r *Renderer) MakeRenderFileAudioParams(b *model.Block) (params *FileMediaRenderParams, err error) {
	file := b.GetFile()
	var src string
	src, err = r.getFileUrl(file.TargetObjectId)
	if err != nil {
		err = fmt.Errorf("file not found %s", file.TargetObjectId)
		return
	}

	params = &FileMediaRenderParams{
		Id:         b.Id,
		Src:        src,
		Classes:    GetAlign(b.GetAlign()),
		Width:		GetWidth(b.Fields),
	}

	return
}

func (r *Renderer) MakeRenderFileVideoParams(b *model.Block) (params *FileMediaRenderParams, err error) {
	file := b.GetFile()
	var src string
	src, err = r.getFileUrl(file.TargetObjectId)
	if err != nil {
		err = fmt.Errorf("file not found %s", file.TargetObjectId)
		return
	}

	params = &FileMediaRenderParams{
		Id:         b.Id,
		Src:        src,
		Classes:    GetAlign(b.GetAlign()),
		Width:		GetWidth(b.Fields),
	}

	return
}

func prettyByteSize(b int64) string {
	bf := float64(b)
	for _, unit := range []string{"", "K", "M", "G", "T", "P", "E", "Z"} {
		if math.Abs(bf) < 1024.0 {
			return fmt.Sprintf("%3.1f%sB", bf, unit)
		}
		bf /= 1024.0
	}
	return fmt.Sprintf("%.1fYiB", bf)
}

func (r *Renderer) MakeRenderFileFileParams(b *model.Block) (params *FileFileRenderParams, err error) {
	file := b.GetFile()
	var src string
	src, err = r.getFileUrl(file.TargetObjectId)
	if err != nil {
		err = fmt.Errorf("file not found %s", file.TargetObjectId)
		return
	}

	name := file.Name
	size := prettyByteSize(file.Size_)

	params = &FileFileRenderParams{
		Id:   b.Id,
		Src:  templ.SafeURL(src),
		Name: name,
		Size: size,
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
	case model.BlockContentFile_File:
		params, err := r.MakeRenderFileFileParams(b)
		if err != nil {
			return NoneTemplate(err.Error())
		}
		return FileFileTemplate(r, params)

	default:
		log.Warn("file type is not supported", zap.String("type", fileType.String()))
		err := fmt.Errorf("file type is not supported: %s", fileType.String())
		return NoneTemplate(err.Error())
	}
}
