package renderer

import (
	"github.com/anyproto/anytype-heart/core/domain"
	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"github.com/gogo/protobuf/jsonpb"
	"path/filepath"
	"testing"
	"time"

	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/assert"
)

func TestMakeRelationRenderParams(t *testing.T) {
	renderer := getTestRenderer("Anytype.WebPublish.20241217.112212.67")

	t.Run("Default case with no relation and no color", func(t *testing.T) {
		// given
		block := &model.Block{
			Id: "test-block-id",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: "nonexistent-key",
				},
			},
		}

		// when
		params := renderer.MakeRelationRenderParams(block)

		// then
		assert.Equal(t, "test-block-id", params.Id)
		assert.Equal(t, "", params.BackgroundColor)
		assert.Equal(t, "isDeleted", params.IsDeleted)
		assert.Equal(t, "isEmpty", params.IsEmpty)
		assert.Equal(t, defaultName, params.Name)
		assert.Equal(t, "", params.Format)
		assert.Nil(t, params.Value)
	})

	t.Run("Block with background color", func(t *testing.T) {
		block := &model.Block{
			Id:              "test-block-id",
			BackgroundColor: "red",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: "nonexistent-key",
				},
			},
		}

		params := renderer.MakeRelationRenderParams(block)

		assert.Equal(t, "bgColor bgColor-red", params.BackgroundColor)
	})

	t.Run("Block with existing relation metadata", func(t *testing.T) {
		// given
		block := &model.Block{
			Id: "test-block-id",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: bundle.RelationKeyName.String(),
				},
			},
		}

		// when
		params := renderer.MakeRelationRenderParams(block)

		// then
		assert.Equal(t, "Name", params.Name)
		assert.Equal(t, "", params.IsDeleted)
		assert.Equal(t, "c-shortText", params.Format)
	})

	t.Run("Block with relation value", func(t *testing.T) {
		// given
		relationKey := "relation-with-value"
		block := &model.Block{
			Id: "test-block-id",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}
		renderer.Sp = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						relationKey: pbtypes.Int64(1),
					},
				},
			}},
		}
		sn := &pb.SnapshotWithType{
			SbType: model.SmartBlockType_STRelation,
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyUniqueKey.String():      pbtypes.String(domain.RelationKey("relation-with-value").URL()),
						bundle.RelationKeyName.String():           pbtypes.String("Value Relation"),
						bundle.RelationKeyRelationFormat.String(): pbtypes.Int64(int64(model.RelationFormat_number)),
					},
				},
			}},
		}
		marshaler := jsonpb.Marshaler{}
		json, err := marshaler.MarshalToString(sn)
		assert.NoError(t, err)
		renderer.UberSp.PbFiles[filepath.Join("relations", "relation-with-value.pb")] = json

		// when
		params := renderer.MakeRelationRenderParams(block)

		// then
		assert.Equal(t, "Value Relation", params.Name)
		assert.Equal(t, "c-number", params.Format)
		assert.NotNil(t, params.Value)
	})

	t.Run("Block with date relation value", func(t *testing.T) {
		// given
		relationKey := "date-relation"

		block := &model.Block{
			Id: "test-block-id",
			Content: &model.BlockContentOfRelation{
				Relation: &model.BlockContentRelation{
					Key: relationKey,
				},
			},
		}

		timestamp := float64(time.Now().Unix())
		renderer.Sp = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						relationKey: pbtypes.Float64(timestamp),
					},
				},
			}},
		}
		sn := &pb.SnapshotWithType{
			SbType: model.SmartBlockType_STRelation,
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyUniqueKey.String():      pbtypes.String(domain.RelationKey("date-relation").URL()),
						bundle.RelationKeyName.String():           pbtypes.String("Date Relation"),
						bundle.RelationKeyRelationFormat.String(): pbtypes.Int64(int64(model.RelationFormat_date)),
					},
				},
			}},
		}
		marshaler := jsonpb.Marshaler{}
		json, err := marshaler.MarshalToString(sn)
		assert.NoError(t, err)
		renderer.UberSp.PbFiles[filepath.Join("relations", "date-relation.pb")] = json

		// when
		params := renderer.MakeRelationRenderParams(block)

		// then
		assert.Equal(t, "Date Relation", params.Name)
		assert.Equal(t, "c-date", params.Format)
		assert.NotNil(t, params.Value)
	})
}
