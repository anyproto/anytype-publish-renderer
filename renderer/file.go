package renderer

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"github.com/anyproto/anytype-publish-renderer/utils"
	"github.com/gogo/protobuf/types"
	"go.uber.org/zap"
)

type FileRenderParams struct {
	Id   string
	Src  templ.SafeURL
	Name string
	Size string
}

func (params *FileRenderParams) ToFileMediaRenderParams(width string, classes []string) *FileMediaRenderParams {
	return &FileMediaRenderParams{
		Id:      params.Id,
		Src:     params.Src,
		Classes: classes,
		Width:   width,
		Name:    params.Name,
		Size:    params.Size,
	}
}

type FileMediaRenderParams struct {
	Id      string
	Src     templ.SafeURL
	Classes []string
	Width   string
	Name    string
	Size    string
}

type SizeSpanRenderParams struct {
	Size string
}

type NameLinkRenderParams struct {
	Name string
	Src  templ.SafeURL
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

func (r *Renderer) getFileBlock(id string) (block *model.Block, err error) {
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
		return bl, nil
	}
	return
}

func GetWidth(fields *types.Struct) string {
	width := pbtypes.GetFloat64(fields, "width")
	log.Debug("image width", zap.Float64("width", width))

	if int(width*100) != 0 {
		return strconv.Itoa(int(width*100)) + "%"
	}
	return ""
}

func GetAlignString(b *model.Block) string {
	align := b.GetAlign()
	return "align" + strconv.Itoa(int(align))
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

func (r *Renderer) MakeRenderFileParams(b *model.Block) (params *FileRenderParams, err error) {
	file := b.GetFile()
	var src string
	src, err = r.getFileUrl(file.TargetObjectId)
	if err != nil {
		err = fmt.Errorf("file not found %s", file.TargetObjectId)
		return
	}

	name := file.Name
	size := prettyByteSize(file.Size_)

	params = &FileRenderParams{
		Id:   b.Id,
		Src:  templ.URL(src),
		Name: name,
		Size: size,
	}

	return
}

func getFileClass(b *model.Block) string {
	file := b.GetFile()
	fileType := file.GetType()
	fileTypeName := model.BlockContentFileType_name[int32(fileType)]
	fileTypeName = utils.Capitalize(strings.ToLower(fileTypeName))
	fileClass := fmt.Sprintf("is%s", fileTypeName)

	return fileClass
}

func (r *Renderer) FileIconBlock(b *model.Block, params *FileRenderParams) templ.Component {
	file := b.GetFile()
	fileClass := getFileClass(b)
	details := r.findTargetDetails(file.TargetObjectId)

	iconParams := r.MakeRenderIconObjectParams(details, &IconObjectProps{})
	iconParams.Classes = append(iconParams.Classes, fileClass)
	iconComp := IconObjectTemplate(r, iconParams)
	return iconComp
}

func (r *Renderer) InlineFileBlock(b *model.Block, params *FileRenderParams) templ.Component {
	iconComp := r.FileIconBlock(b, params)
	nameComp := NameLinkTemplate(&NameLinkRenderParams{
		Name: params.Name,
		Src:  params.Src,
	})
	sizeComp := SizeSpanTemplate(&SizeSpanRenderParams{
		Size: params.Size,
	})

	blockInnerParams := &BlockWrapperParams{
		Classes:    []string{"inner"},
		Components: []templ.Component{iconComp, nameComp, sizeComp},
	}
	blockInner := BlocksWrapper(blockInnerParams)
	blockParams := makeDefaultBlockParams(b)
	blockParams.Content = blockInner

	return BlockTemplate(r, blockParams)
}

func isInlineLink(b *model.Block) bool {
	file := b.GetFile()
	isFile := file.GetType() == model.BlockContentFile_File
	isLink := file.Style == model.BlockContentFile_Link

	return isFile || isLink
}

func (r *Renderer) RenderFile(b *model.Block) templ.Component {
	params, err := r.MakeRenderFileParams(b)
	if err != nil {
		return NoneTemplate(err.Error())
	}

	if isInlineLink(b) {
		return r.InlineFileBlock(b, params)
	} else {
		align := GetAlignString(b)
		classes := []string{align}
		width := GetWidth(b.Fields)

		mediaParams := params.ToFileMediaRenderParams(width, classes)

		var comp templ.Component
		switch b.GetFile().GetType() {
		case model.BlockContentFile_PDF:
			blockParams := makeDefaultBlockParams(b)
			blockParams.Content = FilePDFTemplate(r, mediaParams)
			return BlockTemplate(r, blockParams)
		case model.BlockContentFile_Image:
			comp = ImageTemplate(mediaParams)
		case model.BlockContentFile_Audio:
			comp = AudioTemplate(mediaParams)
		case model.BlockContentFile_Video:
			comp = VideoTemplate(mediaParams)
		default:
			fileTypeStr := b.GetFile().GetType().String()
			log.Warn("file type is not supported", zap.String("type", fileTypeStr))
			return NoneTemplate(fmt.Sprintf("file type is not supported: %s", fileTypeStr))
		}

		var styles map[string]string
		if width != "" {
			styles = map[string]string{
				"width": width,
			}
		}

		blockInnerParams := &BlockWrapperParams{
			Classes:    []string{"wrap"},
			Styles:     styles,
			Components: []templ.Component{comp},
		}
		blockInner := BlocksWrapper(blockInnerParams)
		blockParams := makeDefaultBlockParams(b)
		blockParams.Content = blockInner

		return BlockTemplate(r, blockParams)
	}

}
