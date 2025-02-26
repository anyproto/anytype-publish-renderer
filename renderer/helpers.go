package renderer

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/anyproto/anytype-heart/core/domain"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/gogo/protobuf/types"
)

const linkTemplate = "anytype://object?objectId=%s&spaceId=%s"

func (r *Renderer) findTargetDetails(targetObjectId string) *types.Struct {
	snapshot := r.getObjectSnapshot(targetObjectId)
	if snapshot == nil {
		return nil
	}
	return snapshot.GetSnapshot().GetData().GetDetails()
}

type relType interface {
	string | bool | int64 | model.ObjectTypeLayout | model.RelationFormat
}

type relTransformer[V relType] func(*types.Value) V

func relationToString(field *types.Value) string {
	return field.GetStringValue()
}

func (r *Renderer) relationToEmojiUrl(emojiField *types.Value) string {
	if emojiField.GetStringValue() != "" {
		emojiRune := []rune(emojiField.GetStringValue())[0]
		return r.GetEmojiUrl(emojiRune)
	}

	return ""
}

func (r *Renderer) relationToFileUrl(imageField *types.Value) string {
	if imageField != nil && imageField.GetStringValue() != "" {
		icon, err := r.getFileUrl(imageField.GetStringValue())
		if err != nil {
			log.Error("Failed to get file URL for icon", zap.Error(err))
			return ""
		}
		return icon
	}

	return ""
}

func relationToBool(boolField *types.Value) bool {
	var null bool
	if boolField != nil {
		return boolField.GetBoolValue()
	}

	return null
}

func relationToObjectTypeLayout(layout *types.Value) model.ObjectTypeLayout {
	if layout != nil {
		return model.ObjectTypeLayout(layout.GetNumberValue())
	}

	return model.ObjectType_basic
}

func relationToRelationFormat(format *types.Value) model.RelationFormat {
	if format != nil {
		return model.RelationFormat(format.GetNumberValue())
	}

	return model.RelationFormat_longtext
}

func relationToInt64(field *types.Value) int64 {
	var null int64
	if field != nil {
		return int64(field.GetNumberValue())
	}
	return null
}

func getRelationField[V relType](targetDetails *types.Struct, relationKey domain.RelationKey, tr relTransformer[V]) V {
	var null V
	if f, ok := targetDetails.GetFields()[relationKey.String()]; ok {
		return tr(f)
	}

	return null
}

func (r *Renderer) makeAnytypeLink(targetDetails *types.Struct, targetObjectId string) string {
	layout := getRelationField(targetDetails, bundle.RelationKeyLayout, relationToObjectTypeLayout)
	switch layout {
	case model.ObjectType_file, model.ObjectType_image, model.ObjectType_pdf, model.ObjectType_audio, model.ObjectType_video:
		src, err := r.getFileUrl(targetObjectId)
		if err != nil {
			log.Error("failed to get file url", zap.Error(err))
			return ""
		}
		return src
	default:
		spaceId := getRelationField(targetDetails, bundle.RelationKeySpaceId, relationToString)
		return fmt.Sprintf(linkTemplate, targetObjectId, spaceId)
	}
}

func getLayoutClass(layout model.ObjectTypeLayout) string {
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
