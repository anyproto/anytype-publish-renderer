package renderer

import (
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

func (r *Renderer) MakeLinkRenderParams(b *model.Block) (params *LinkRenderParams) {
	targetObjectId := b.GetLink().GetTargetBlockId()
	var targetObjectIdDetails *types.Struct
	for _, detail := range r.Sp.GetDependantDetails() {
		if detail.Id == targetObjectId {
			targetObjectIdDetails = detail.Details
			break
		}
	}
	if targetObjectIdDetails == nil || len(targetObjectIdDetails.Fields) == 0 {
		return &LinkRenderParams{IsDeleted: true}
	}
	var linkTypeClass string
	if b.GetLink().GetCardStyle() == model.BlockContentLink_Card {
		linkTypeClass = "card"
	}
	if b.GetLink().GetCardStyle() == model.BlockContentLink_Text {
		linkTypeClass = "text"
	}
	var (
		description      string
		descriptionValue *types.Value
	)
	if b.GetLink().GetDescription() == model.BlockContentLink_Content {
		descriptionValue = targetObjectIdDetails.GetFields()[bundle.RelationKeySnippet.String()]
	}
	if b.GetLink().GetDescription() == model.BlockContentLink_Added {
		descriptionValue = targetObjectIdDetails.GetFields()[bundle.RelationKeyDescription.String()]
	}
	if descriptionValue != nil {
		description = descriptionValue.GetStringValue()
	}

	deletedValue := targetObjectIdDetails.GetFields()[bundle.RelationKeyIsDeleted.String()]
	if deletedValue != nil && deletedValue.GetBoolValue() {
		return &LinkRenderParams{IsDeleted: true}
	}
	name := targetObjectIdDetails.GetFields()[bundle.RelationKeyName.String()].GetStringValue()
	if name == "" {
		name = defaultName
	}
	var iconStyle, icon string
	iconClass := "c20"
	if b.GetLink().GetIconSize() != model.BlockContentLink_SizeNone && b.GetLink().GetCardStyle() == model.BlockContentLink_Text {
		iconStyle = "smileImage c20"
	}
	if b.GetLink().GetCardStyle() == model.BlockContentLink_Card {
		if b.GetLink().GetIconSize() == model.BlockContentLink_SizeMedium {
			iconStyle = "smileImage c28"
			iconClass = "c48"
		}
		if b.GetLink().GetIconSize() == model.BlockContentLink_SizeSmall {
			iconStyle = "smileImage c20"
		}
	}
	layout := model.ObjectTypeLayout(targetObjectIdDetails.GetFields()[bundle.RelationKeyLayout.String()].GetNumberValue())

	var layoutClass string
	switch layout {
	case model.ObjectType_participant:
		layoutClass = "isParticipant"
	case model.ObjectType_profile:
		layoutClass = "isHuman"
	case model.ObjectType_todo:
		layoutClass = "isTask"
	case model.ObjectType_collection:
		layoutClass = "isCollection"
	case model.ObjectType_set:
		layoutClass = "isSet"
	default:
		layoutClass = "isPage"
	}

	if b.GetLink().GetIconSize() != model.BlockContentLink_SizeNone {
		if layout == model.ObjectType_todo {
			iconStyle = "iconCheckbox c20 icon checkbox unset"
			doneValue := targetObjectIdDetails.GetFields()[bundle.RelationKeyDone.String()]
			if doneValue != nil && doneValue.GetBoolValue() {
				iconStyle = "iconCheckbox c20 icon checkbox set"
			}
		} else {
			iconValue := targetObjectIdDetails.GetFields()[bundle.RelationKeyIconEmoji.String()]
			if iconValue != nil && iconValue.GetStringValue() != "" {
				emojiRune := []rune(iconValue.GetStringValue())[0]
				icon = r.GetEmojiUrl(emojiRune)
				iconClass = iconClass + " withIcon"
			}
			if iconValue == nil || iconValue.GetStringValue() == "" {
				iconValue = targetObjectIdDetails.GetFields()[bundle.RelationKeyIconImage.String()]
				if iconValue == nil || iconValue.GetStringValue() == "" {
					iconSize := "c20"
					if b.GetLink().GetIconSize() == model.BlockContentLink_SizeMedium {
						iconSize = "c28"
					}
					if layout == model.ObjectType_collection || layout == model.ObjectType_set {
						iconStyle = "iconCommon icon collection " + iconSize
					}
					if layout != model.ObjectType_note && layout != model.ObjectType_profile && layout != model.ObjectType_participant {
						iconStyle = "iconCommon icon page " + iconSize
					}
				}
				var err error
				icon, err = r.getFileUrl(icon)
				if err != nil {
					log.Error("failed to get file url", zap.Error(err))
				} else {
					iconClass = iconClass + " withImage"
				}
			}
		}
	}

	archivedValue := targetObjectIdDetails.GetFields()[bundle.RelationKeyIsArchived.String()]

	var archiveClass string
	if archivedValue != nil && archivedValue.GetBoolValue() {
		archiveClass = "isArchived"
	}

	var (
		objectTypeName string
		coverParams    *CoverRenderParams
		err            error
		coverClass     string
	)
	for _, relation := range b.GetLink().GetRelations() {
		if relation == bundle.RelationKeyType.String() {
			objectType := targetObjectIdDetails.GetFields()[bundle.RelationKeyType.String()].GetStringValue()
			for _, detail := range r.Sp.GetDependantDetails() {
				if detail.Id == objectType {
					objectTypeName = detail.Details.GetFields()[bundle.RelationKeyName.String()].GetStringValue()
					break
				}
			}
		}
		if relation == "cover" {
			coverParams, err = r.getCoverParams(targetObjectIdDetails)
			if err != nil {
				log.Error("failed to get cover params", zap.Error(err))
			} else {
				coverClass = "withCover"
			}
		}
	}

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

func (r *Renderer) RenderLink(b *model.Block) templ.Component {
	params := r.MakeLinkRenderParams(b)
	return LinkTempl(params)
}
