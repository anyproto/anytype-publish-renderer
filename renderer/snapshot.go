package renderer

import (
	"github.com/anyproto/anytype-heart/core/domain"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/gogo/protobuf/types"
	"go.uber.org/zap"
)

func (r *Renderer) findTargetDetails(targetObjectId string) *types.Struct {
	snapshot := r.getObjectSnapshot(targetObjectId)
	if snapshot == nil {
		return nil
	}
	return snapshot.GetSnapshot().GetData().GetDetails()
}

type relType interface {
	string | bool | model.ObjectTypeLayout
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

func getRelationField[V relType](targetDetails *types.Struct, relationKey domain.RelationKey, tr relTransformer[V]) V {
	var null V
	if f, ok := targetDetails.GetFields()[relationKey.String()]; ok {
		return tr(f)
	}

	return null
}
