package renderer

import (
	"encoding/base64"
	"fmt"
	"slices"
	"strings"
	"unicode"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-publish-renderer/utils"
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

var FontSize = map[int]int{
	14:  10,
	16:  10,
	18:  11,
	20:  13,
	22:  14,
	24:  16,
	26:  16,
	30:  20,
	32:  20,
	36:  24,
	40:  24,
	42:  24,
	44:  24,
	48:  28,
	56:  40,
	64:  40,
	80:  64,
	96:  64,
	108: 64,
	128: 64,
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

type UserSvgProps struct {
	Size       string
	ViewBox    string
	FontWeight string
	FontSize   string
	Letter     string
}

func FirstAlnumChar(s string, defaultLetter string) string {
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return string(unicode.ToUpper(r))
		}
	}

	return defaultLetter
}

// EncodeSVGToDataURL converts an SVG string to a Base64 data URL
func EncodeSVGToDataURL(svg string) string {
	// Base64 encode the SVG string directly
	base64SVG := base64.StdEncoding.EncodeToString([]byte(svg))

	// Construct the final data URL

	return "data:image/svg+xml;charset=utf-8;base64," + base64SVG
}

func makeUserSvgProps(size int, username string) *UserSvgProps {
	sizeStr := fmt.Sprintf("%dpx", size)
	viewBox := fmt.Sprintf("0 0 %d %d", size, size)

	fontWeight := "500"
	if size > 18 {
		fontWeight = "600"
	}

	fontSize := 72
	if fs, ok := FontSize[size]; ok {
		fontSize = min(fontSize, fs)
	}
	fontSizeStr := fmt.Sprintf("%dpx", fontSize)

	// "U" stands for "Untitled"
	letter := FirstAlnumChar(username, "U")

	return &UserSvgProps{
		Size:       sizeStr,
		ViewBox:    viewBox,
		FontWeight: fontWeight,
		FontSize:   fontSizeStr,
		Letter:     letter,
	}
}

// for some reason, templ create </circle> closing tag
// which breaks the image rendering.
func makeSvgString(props *UserSvgProps) string {
	return fmt.Sprintf(`
<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" id="Layer_1" x="0px" y="0px" viewBox="%s" xml:space="preserve" height="%s" width="%s">
	<circle cx="50%%" cy="50%%" r="50%%" fill="#f2f2f2" />
	<text x="50%%" y="50%%" text-anchor="middle" dominant-baseline="central" fill="#b6b6b6" font-family="Inter, Helvetica" font-weight="%s" font-size="%s">%s</text>
</svg>`, props.ViewBox, props.Size, props.Size, props.FontWeight, props.FontSize, props.Letter)
}

func getIconSize(props *IconObjectProps, layout model.ObjectTypeLayout, gsProps *GetSizeProps) int32 {
	s, ok := IconSize[props.Size]
	if !ok {
		s = props.Size
	}

	if gsProps.IsDeleted {
		return s
	}
	if props.Size == 18 && isTodoLayout(layout) {
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
			if isTodoLayout(layout) && layout == model.ObjectType_relation && props.ForceLetter {
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

func (r *Renderer) getDefaultIconPath(name string) (path string) {
	path = r.GetStaticFolderUrl(fmt.Sprintf("/img/icon/default/%s.svg", name))
	return
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
	hasIconEmoji := iconEmoji != ""
	hasIconImage := iconImage != ""
	defaultIcon := "page"

	if hasIconImage {
		src = iconImage
	}

	switch layout {
	default:
		fallthrough
	case model.ObjectType_collection, model.ObjectType_set:
		defaultIcon = "set"
		fallthrough
	case model.ObjectType_basic:
		if hasIconEmoji {
			iconClasses = append(iconClasses, "smileImage")
			src = iconEmoji
		} else if hasIconImage {
			classes = append(classes, "withImage")
			iconClasses = append(iconClasses, "iconImage")
		} else {
			classes = append(classes, "withDefault")
			iconClasses = append(iconClasses, "iconCommon")
			src = r.getDefaultIconPath(defaultIcon)
		}

		if props.ForceLetter {
			classes = append(classes, "withLetter")
			// todo: commonsvg
		}
	case model.ObjectType_participant, model.ObjectType_profile:
		classes = append(classes, "isHuman")
		iconClasses = append(iconClasses, "iconImage")
		if hasIconImage {
			classes = append(classes, "withImage")
		} else {
			name := getRelationField(targetDetails, bundle.RelationKeyName, relationToString)
			if name == "" {
				name = "Untitled"
			}
			props := makeUserSvgProps(128, name)
			svg := makeSvgString(props)
			src = EncodeSVGToDataURL(svg)
		}

	case model.ObjectType_date:
		defaultIcon = "date"
		classes = append(classes, "withDefault")
		iconClasses = append(iconClasses, "iconCommon")
		src = r.getDefaultIconPath(defaultIcon)
	case model.ObjectType_todo:
		done := getRelationField(targetDetails, bundle.RelationKeyDone, relationToBool)
		checkIconNum := 0
		if done {
			checkIconNum = 2
		}
		src = r.GetStaticFolderUrl(fmt.Sprintf("/img/icon/object/checkbox%d.svg", checkIconNum))
		iconClasses = append(iconClasses, "iconCheckbox")
	case model.ObjectType_note:
		defaultIcon = "page"
		classes = append(classes, "withDefault")
		iconClasses = append(iconClasses, "iconCommon")
		src = r.getDefaultIconPath(defaultIcon)
	case model.ObjectType_objectType:
		if hasIconEmoji {
			iconClasses = append(iconClasses, "smileImage")
			src = iconEmoji
		} else {
			defaultIcon = "type"
			classes = append(classes, "withDefault")
			iconClasses = append(iconClasses, "iconCommon")
			src = r.getDefaultIconPath(defaultIcon)
		}
	case model.ObjectType_relation:
		format := getRelationField(targetDetails, bundle.RelationKeyRelationFormat, relationToRelationFormat)
		if format != model.RelationFormat_relations && format != model.RelationFormat_emoji {
			iconClasses = append(iconClasses, "iconCommon")
			typeName := utils.Capitalize(model.RelationFormat_name[int32(format)])
			src = r.GetStaticFolderUrl(fmt.Sprintf("/img/icon/relation/%s.svg", typeName))
		}
	case model.ObjectType_bookmark:
		// TODO: should show image preview when we will have cropped images in snapshot
		iconClasses = append(iconClasses, "iconFile")
		iconName := "image"
		src = r.GetStaticFolderUrl(fmt.Sprintf("/img/icon/file/%s.svg", iconName))
	case model.ObjectType_image:
		// TODO: should show image preview when we will have cropped images in snapshot
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

	case model.ObjectType_spaceView, model.ObjectType_dashboard:
		break
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
