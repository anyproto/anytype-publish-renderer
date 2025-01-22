package renderer

import (
	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/gogo/protobuf/types"
)

const linkTemplate = "anytype://object?objectId=%s&spaceId=%s"

type LinkRenderParams struct {
	Id            string
	LayoutClass   string
	IsDeleted     bool
	IsArchived    string
	Name          string
	Description   string
	Type          string
	Icon          string
	IconClass     string
	IconStyle     string
	LinkTypeClass string
	CoverClass    string
	CoverParams   *CoverRenderParams
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

	name := getFieldValue(targetDetails, bundle.RelationKeyName.String(), defaultName)
	icon, iconClass, iconStyle := r.getIconParams(b, targetDetails)
	layoutClass := getLayoutClass(targetDetails)
	archiveClass := getArchiveClass(targetDetails)

	objectTypeName, coverParams, coverClass := r.getAdditionalParams(b, targetDetails)

	return &LinkRenderParams{
		Id:            b.GetId(),
		LayoutClass:   layoutClass,
		IsArchived:    archiveClass,
		Name:          name,
		Description:   description,
		Type:          objectTypeName,
		Icon:          icon,
		IconClass:     iconClass,
		IconStyle:     iconStyle,
		LinkTypeClass: linkTypeClass,
		CoverClass:    coverClass,
		CoverParams:   coverParams,
	}
}

func (r *Renderer) findTargetDetails(targetObjectId string) *types.Struct {
	for _, detail := range r.Sp.GetDependantDetails() {
		if detail.Id == targetObjectId {
			return detail.Details
		}
	}
	return nil
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
	if layout == model.ObjectType_todo && b.GetLink().GetIconSize() != model.BlockContentLink_SizeNone {
		iconStyle = "iconCheckbox c20 icon checkbox unset"
		if doneValue := details.GetFields()[bundle.RelationKeyDone.String()]; doneValue != nil && doneValue.GetBoolValue() {
			iconStyle = "iconCheckbox c20 icon checkbox set"
		}
		return
	}
	if b.GetLink().GetIconSize() != model.BlockContentLink_SizeNone {
		iconStyle = "smileImage c20"
	}
	if b.GetLink().GetCardStyle() == model.BlockContentLink_Card {
		if b.GetLink().GetIconSize() == model.BlockContentLink_SizeMedium {
			iconStyle = "smileImage c28"
			iconClass = "c48"
		} else if b.GetLink().GetIconSize() == model.BlockContentLink_SizeSmall {
			iconStyle = "smileImage c20"
		}
	}
	iconValue := details.GetFields()[bundle.RelationKeyIconEmoji.String()]
	if iconValue != nil && iconValue.GetStringValue() != "" {
		emojiRune := []rune(iconValue.GetStringValue())[0]
		icon = r.GetEmojiUrl(emojiRune)
		iconClass += " withIcon"
	}
	if iconValue == nil || iconValue.GetStringValue() == "" {
		iconValue = details.GetFields()[bundle.RelationKeyIconImage.String()]
		if iconValue != nil && iconValue.GetStringValue() != "" {
			icon, _ = r.getFileUrl(iconValue.GetStringValue())
			iconClass += " withImage"
		}
	}
	if icon == "" && b.GetLink().GetIconSize() != model.BlockContentLink_SizeNone {
		iconSize := "c20"
		if b.GetLink().GetIconSize() == model.BlockContentLink_SizeMedium {
			iconSize = "c28"
		}
		if layout == model.ObjectType_collection || layout == model.ObjectType_set {
			iconStyle = "iconCommon icon collection " + iconSize
		} else if layout != model.ObjectType_note && layout != model.ObjectType_profile && layout != model.ObjectType_participant {
			iconStyle = "iconCommon icon page " + iconSize
		}
	}
	return
}

func (r *Renderer) getAdditionalParams(b *model.Block, details *types.Struct) (objectTypeName string, coverParams *CoverRenderParams, coverClass string) {
	for _, relation := range b.GetLink().GetRelations() {
		if relation == bundle.RelationKeyType.String() {
			objectType := details.GetFields()[bundle.RelationKeyType.String()].GetStringValue()
			for _, detail := range r.Sp.GetDependantDetails() {
				if detail.Id == objectType {
					objectTypeName = detail.Details.GetFields()[bundle.RelationKeyName.String()].GetStringValue()
					break
				}
			}
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
