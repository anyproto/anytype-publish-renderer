package renderer

import (
	"path/filepath"
	"strings"

	"github.com/anyproto/anytype-heart/core/domain"
	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/localstore/addr"
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

func (r *Renderer) getObjectSnapshot(objectId string) *pb.SnapshotWithType {
	if strings.HasPrefix(objectId, addr.DatePrefix) {
		return r.getDateSnapshot(objectId)
	}
	directories := []string{"objects", "relations", "types", "templates", "filesObjects"}
	var (
		snapshot *pb.SnapshotWithType
		err      error
	)
	for _, dir := range directories {
		path := filepath.Join(dir, objectId+".pb")
		snapshot, err = r.ReadJsonpbSnapshot(path)
		if err == nil {
			return snapshot
		}
	}
	log.Error("failed to get snapshot for object", zap.String("objectId", objectId), zap.Error(err))
	return nil
}

func relationToBool(boolField *types.Value) bool {
	var null bool
	if boolField != nil {
		return boolField.GetBoolValue()
	}

	return null
}

func relationToInt64(field *types.Value) int64 {
	var null int64
	if field != nil {
		return int64(field.GetNumberValue())
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

func getRelationField[V relType](targetDetails *types.Struct, relationKey domain.RelationKey, tr relTransformer[V]) V {
	var null V
	if f, ok := targetDetails.GetFields()[relationKey.String()]; ok {
		return tr(f)
	}

	return null
}
