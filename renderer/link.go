package renderer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/gogo/protobuf/types"
	"go.uber.org/zap"
)

const linkTemplate = "anytype://object?objectId=%s&spaceId=%s"

type LinkRenderParams struct {
	Id            string
	LayoutClass   string
	Classes       string
	ContentClasses string
	SidesClasses  string
	CardClasses   string
	IsDeleted     bool
	IsArchived    string
	Name          string
	Description   string
	Type          string
	Icon          string
	IconClass     string
	IconStyle     string
	CoverClass    string
	CoverParams   *CoverRenderParams
	Url           templ.SafeURL
}

func (r *Renderer) MakeLinkRenderParams(b *model.Block) *LinkRenderParams {
	targetObjectId := b.GetLink().GetTargetBlockId()
	targetDetails := r.findTargetDetails(targetObjectId)
	if targetDetails == nil || len(targetDetails.Fields) == 0 {
		return &LinkRenderParams{IsDeleted: true}
	}

	linkTypeClass := getLinkTypeClass(b)
	description := getDescription(b, targetDetails)
	if isDeleted(targetDetails) {
		return &LinkRenderParams{IsDeleted: true}
	}

	bgColor := b.GetBackgroundColor()
	name := getFieldValue(targetDetails, bundle.RelationKeyName.String(), defaultName)
	icon, iconClass, iconStyle := r.getIconParams(b, targetDetails)
	layoutClass := getLayoutClass(targetDetails)
	archiveClass := getArchiveClass(targetDetails)
	objectTypeName, coverParams, coverClass := r.getAdditionalParams(b, targetDetails)
	spaceId := targetDetails.GetFields()[bundle.RelationKeySpaceId.String()].GetStringValue()
	link := fmt.Sprintf(linkTemplate, targetObjectId, spaceId)
	classes := []string{linkTypeClass}
	contentClasses := []string{"content"}
	sidesClasses := []string{"sides"}
	cardClasses := []string{"linkCard", iconClass, layoutClass, coverClass}

	if bgColor != "" {
		sidesClasses = append(sidesClasses, "withBgColor")
		contentClasses = append(contentClasses, "bgColor", "bgColor-" + bgColor)
	}

	return &LinkRenderParams{
		Id:            b.GetId(),
		Classes:       strings.Join(classes, " "),
		ContentClasses: strings.Join(contentClasses, " "),
		SidesClasses:  strings.Join(sidesClasses, " "),
		CardClasses:   strings.Join(cardClasses, " "),
		LayoutClass:   layoutClass,
		IsArchived:    archiveClass,
		Name:          name,
		Description:   description,
		Type:          objectTypeName,
		Icon:          icon,
		IconClass:     iconClass,
		IconStyle:     iconStyle,
		CoverClass:    coverClass,
		CoverParams:   coverParams,
		Url:           templ.SafeURL(link),
	}
}

func (r *Renderer) findTargetDetails(targetObjectId string) *types.Struct {
	snapshot := r.getObjectSnapshot(targetObjectId)
	if snapshot == nil {
		return nil
	}
	return snapshot.GetSnapshot().GetData().GetDetails()
}

func getLinkTypeClass(b *model.Block) string {
	switch b.GetLink().GetCardStyle() {
		case model.BlockContentLink_Card:
			return "card"
		default:
			return "text"
	}
}

func getDescription(b *model.Block, details *types.Struct) string {
	var key string

	switch b.GetLink().GetDescription() {
		case model.BlockContentLink_Content:
			key = bundle.RelationKeySnippet.String()
		case model.BlockContentLink_Added:
			key = bundle.RelationKeyDescription.String()
		default:
			return ""
	}

	descriptionValue := details.GetFields()[key]
	if descriptionValue != nil {
		return descriptionValue.GetStringValue()
	}
	return ""
}

func isDeleted(details *types.Struct) bool {
	deletedValue := details.GetFields()[bundle.RelationKeyIsDeleted.String()]
	return deletedValue != nil && deletedValue.GetBoolValue()
}

func getFieldValue(details *types.Struct, key, defaultValue string) string {
	value := details.GetFields()[key]
	if value == nil || value.GetStringValue() == "" {
		return defaultValue
	}
	return value.GetStringValue()
}

func getLayoutClass(details *types.Struct) string {
	layout := model.ObjectTypeLayout(details.GetFields()[bundle.RelationKeyLayout.String()].GetNumberValue())

	switch layout {
		case model.ObjectType_participant:
			return "isParticipant"
		case model.ObjectType_profile:
			return "isHuman"
		case model.ObjectType_todo:
			return "isTask"
		case model.ObjectType_collection:
			return "isCollection"
		case model.ObjectType_set:
			return "isSet"
		default:
			return "isPage"
	}
}

func getArchiveClass(details *types.Struct) string {
	archivedValue := details.GetFields()[bundle.RelationKeyIsArchived.String()]
	if archivedValue != nil && archivedValue.GetBoolValue() {
		return "isArchived"
	}
	return ""
}

func (r *Renderer) getIconParams(b *model.Block, details *types.Struct) (icon, iconClass, iconStyle string) {
	iconClass = "c20"

	if b.GetLink().GetIconSize() == model.BlockContentLink_SizeNone {
		return
	}

	layout := model.ObjectTypeLayout(details.GetFields()[bundle.RelationKeyLayout.String()].GetNumberValue())

	if layout == model.ObjectType_todo {
		iconStyle = r.getTodoIconStyle(details)
		return
	}
	iconStyle, iconClass = r.getDefaultIconStyle(b, iconClass)

	icon, iconClass = r.getIconFromDetails(details, iconClass)

	if icon == "" {
		iconStyle = r.getFallbackIconStyle(b, layout)
	}
	return
}

func (r *Renderer) getTodoIconStyle(details *types.Struct) string {
	iconStyle := "iconCheckbox c20 icon checkbox unset"
	if doneValue := details.GetFields()[bundle.RelationKeyDone.String()]; doneValue != nil && doneValue.GetBoolValue() {
		iconStyle = "iconCheckbox c20 icon checkbox set"
	}
	return iconStyle
}

func (r *Renderer) getDefaultIconStyle(b *model.Block, iconClass string) (iconStyle, updatedIconClass string) {
	iconStyle = "smileImage c20"
	updatedIconClass = iconClass

	if b.GetLink().GetCardStyle() == model.BlockContentLink_Card {
		switch b.GetLink().GetIconSize() {
		case model.BlockContentLink_SizeMedium:
			iconStyle = "smileImage c28"
			updatedIconClass = "c48"
		case model.BlockContentLink_SizeSmall:
			iconStyle = "smileImage c20"
		}
	}
	return iconStyle, updatedIconClass
}

func (r *Renderer) getIconFromDetails(details *types.Struct, iconClass string) (icon, updatedIconClass string) {
	emojiField := details.GetFields()[bundle.RelationKeyIconEmoji.String()]
	if emojiField != nil && emojiField.GetStringValue() != "" {
		emojiRune := []rune(emojiField.GetStringValue())[0]
		icon = r.GetEmojiUrl(emojiRune)
		return icon, iconClass + " withIcon"
	}

	imageField := details.GetFields()[bundle.RelationKeyIconImage.String()]
	if imageField != nil && imageField.GetStringValue() != "" {
		icon, err := r.getFileUrl(imageField.GetStringValue())
		if err != nil {
			log.Error("Failed to get file URL for icon", zap.Error(err))
			return "", iconClass
		}
		return icon, iconClass + " withImage"
	}

	return "", iconClass
}

func (r *Renderer) getFallbackIconStyle(b *model.Block, layout model.ObjectTypeLayout) string {
	iconSize := "c20"
	if b.GetLink().GetIconSize() == model.BlockContentLink_SizeMedium {
		iconSize = "c28"
	}

	switch layout {
	case model.ObjectType_collection, model.ObjectType_set:
		return "iconCommon icon collection " + iconSize
	case model.ObjectType_profile, model.ObjectType_participant:
		return "iconImage " + iconSize
	case model.ObjectType_note:
		return ""
	default:
		return "iconCommon icon page " + iconSize
	}
}

func (r *Renderer) getAdditionalParams(b *model.Block, details *types.Struct) (objectTypeName string, coverParams *CoverRenderParams, coverClass string) {
	for _, relation := range b.GetLink().GetRelations() {
		if relation == bundle.RelationKeyType.String() {
			objectType := details.GetFields()[bundle.RelationKeyType.String()].GetStringValue()
			snapshot, err := r.ReadJsonpbSnapshot(filepath.Join("types", objectType+".pb"))
			if err != nil {
				log.Error("failed to read jsonpb snapshot", zap.Error(err))
				continue
			}
			objectTypeName = snapshot.GetSnapshot().GetData().GetDetails().GetFields()[bundle.RelationKeyName.String()].GetStringValue()
		}
		if relation == "cover" {
			var err error
			coverParams, err = r.getCoverParams(details)
			if err == nil {
				coverClass = "withCover"
			}
		}
	}
	return
}

func (r *Renderer) RenderLink(b *model.Block) templ.Component {
	params := r.MakeLinkRenderParams(b)
	return LinkTempl(params)
}
