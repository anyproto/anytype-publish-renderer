package renderer

import (
	"fmt"
	"slices"
	"strings"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/gogo/protobuf/types"
)

var IconSize = map[int32]int32{
	14:  14,
	16:  16,
	18:  16,
	20:  18,
	22:  18,
	24:  20,
	26:  22,
	28:  22,
	32:  28,
	36:  24,
	40:  24,
	42:  24,
	44:  24,
	48:  24,
	56:  32,
	64:  32,
	80:  56,
	96:  56,
	108: 64,
	112: 64,
	128: 64,
	160: 160,
	360: 360,
}

type IconObjectParams struct {
	Classes     []string
	IconClasses []string
	Src         string
}

type IconObjectProps struct {
	ClassName   string
	IconClass   string
	Size        int32
	IconSize    int32
	ForceLetter bool
	Src         templ.SafeURL
}

type GetSizeProps struct {
	HasIconImage bool
	HasIconEmoji bool
	IsDeleted    bool
}

func getIconSize(props *IconObjectProps, layout model.ObjectTypeLayout, gsProps *GetSizeProps) int32 {
	s, ok := IconSize[props.Size]
	if !ok {
		s = props.Size
	}

	if gsProps.IsDeleted {
		return s
	}
	if props.Size == 18 && layout == model.ObjectType_todo {
		s = 16
	} else if props.Size == 48 && layout == model.ObjectType_relation {
		s = 28
	} else if props.Size >= 40 {
		if isHumanLayout(layout) {
			s = props.Size
		}
		if (layout == model.ObjectType_set || layout == model.ObjectType_spaceView) && gsProps.HasIconImage {
			s = props.Size
		}
		if !gsProps.HasIconImage && !gsProps.HasIconEmoji {
			if layout == model.ObjectType_set || layout == model.ObjectType_objectType {
				s = props.Size
			}
			if layout == model.ObjectType_todo && layout == model.ObjectType_relation && props.ForceLetter {
				s = props.Size
			}
		}
	}

	if props.IconSize != 0 {
		s = props.IconSize
	}

	return s
}

var fileExtensions = map[string][]string{
	"image": {"jpg", "jpeg", "png", "gif", "svg", "webp"},
	"video": {"mp4", "m4v", "mov"},
	"audio": {"mp3", "m4a", "flac", "ogg", "wav"},
	"pdf":   {"pdf"},
}

func fileIconName(details *types.Struct) string {
	name := getRelationField(details, bundle.RelationKeyName, relationToString)
	mime := getRelationField(details, bundle.RelationKeyFileMimeType, relationToString)
	fileExt := getRelationField(details, bundle.RelationKeyFileExt, relationToString)
	n := strings.Split(name, ".")
	e := ""

	if fileExt != "" {
		e = strings.ToLower(fileExt)
	} else if len(n) > 1 {
		e = strings.ToLower(n[len(n)-1])
	}

	icon := "other"
	var t []string
	if mime != "" {
		splitMime := strings.Split(mime, ";")
		if len(splitMime) > 0 {
			t = strings.Split(splitMime[0], "/")
		}
	}

	if len(t) > 0 {
		switch t[0] {
		case "image", "video", "text", "audio":
			icon = t[0]
		}
		switch t[1] {
		case "pdf":
			icon = "pdf"
		case "zip", "gzip", "tar", "gz", "rar":
			icon = "archive"
		case "vnd.ms-powerpoint":
			icon = "presentation"
		case "vnd.openxmlformats-officedocument.spreadsheetml.sheet":
			icon = "table"
		}
	}

	switch e {
	case "m4v":
		icon = "video"
	case "csv", "json", "txt", "doc", "docx", "md", "tsx", "scss", "html", "yml", "rtf":
		icon = "text"
	case "zip", "gzip", "tar", "gz", "rar":
		icon = "archive"
	case "xls", "xlsx", "sqlite":
		icon = "table"
	case "ppt", "pptx", "key":
		icon = "presentation"
	case "aif":
		icon = "audio"
	case "ai":
		icon = "image"
	case "dwg":
		icon = "other"

	}

	for k, v := range fileExtensions {
		if slices.Contains(v, e) {
			icon = k
			break
		}
	}

	return icon
}

func (r *Renderer) MakeRenderIconObjectParams(targetDetails *types.Struct, props *IconObjectProps) (params *IconObjectParams) {
	var src string
	classes := []string{"iconObject"}
	var iconClasses []string
	var isDeleted bool
	if targetDetails == nil || len(targetDetails.Fields) == 0 {
		isDeleted = true
	}

	layout := getRelationField(targetDetails, bundle.RelationKeyLayout, relationToObjectTypeLayout)
	iconEmoji := getRelationField(targetDetails, bundle.RelationKeyIconEmoji, r.relationToEmojiUrl)
	iconImage := getRelationField(targetDetails, bundle.RelationKeyIconImage, r.relationToFileUrl)
	// done := getRelationField(targetDetails, bundle.RelationKeyDone, relationToBool)
	hasIconEmoji := iconEmoji != ""
	hasIconImage := iconImage != ""

	//iconClass
	//done
	//relationFormat
	if hasIconImage {
		src = iconImage
	}

	switch layout {
	default:
		fallthrough
	case model.ObjectType_basic:
		if hasIconEmoji {
			iconClasses = append(iconClasses, "smileImage")
			src = iconEmoji
		} else if hasIconImage {
			classes = append(classes, "withImage")
			iconClasses = append(iconClasses, "iconImage")
		}

		if props.ForceLetter {
			classes = append(classes, "withLetter")
			// todo: commonsvg
		}
	case model.ObjectType_participant:
		fallthrough
	case model.ObjectType_profile:
		classes = append(classes, "isHuman")
		if hasIconImage {
			classes = append(classes, "withImage")
			iconClasses = append(iconClasses, "iconImage")
		}

	// case model.ObjectType_set:

	// case model.ObjectType_todo:
	// case model.ObjectType_dashboard:
	// case model.ObjectType_note:
	// case model.ObjectType_objectType:
	// case model.ObjectType_relation:
	// case model.ObjectType_bookmark:
	// case model.ObjectType_spaceView:
	case model.ObjectType_image:
		iconClasses = append(iconClasses, "iconImage")
		// TODO: should show image preview
		// getFileUrl
		fallthrough
	case model.ObjectType_video:
		fallthrough
	case model.ObjectType_audio:
		fallthrough
	case model.ObjectType_pdf:
		fallthrough
	case model.ObjectType_file:
		iconClasses = append(iconClasses, "iconFile")
		iconName := fileIconName(targetDetails)
		src = r.GetStaticFolderUrl(fmt.Sprintf("/img/icon/file/%s.svg", iconName))
	}

	if props.Size != 0 {
		classes = append(classes, fmt.Sprintf("c%d", props.Size))
	}

	gsProps := &GetSizeProps{
		HasIconEmoji: hasIconEmoji,
		HasIconImage: hasIconImage,
		IsDeleted:    isDeleted,
	}
	iconSize := getIconSize(props, layout, gsProps)
	if iconSize != 0 {
		iconClasses = append(iconClasses, fmt.Sprintf("c%d", iconSize))
	}

	if isDeleted {
		src = r.GetStaticFolderUrl("/img/icon/ghost.svg")
		iconClasses = []string{"iconCommon"}
		if iconSize != 0 {
			iconClasses = append(iconClasses, fmt.Sprintf("c%d", iconSize))
		}
	}

	return &IconObjectParams{
		Classes:     classes,
		IconClasses: iconClasses,
		Src:         src,
	}
}
