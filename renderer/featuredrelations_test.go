package renderer

import (
	"github.com/anyproto/anytype-heart/core/domain"
	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestMakeFeaturedRelationsParams(t *testing.T) {
	r := getTestRenderer("Anytype.WebPublish.20241217.112212.67")

	t.Run("empty details should return empty params", func(t *testing.T) {
		// given
		block := &model.Block{Id: "block1"}
		r.Sp = &pb.SnapshotWithType{}

		// when
		params := r.MakeFeaturedRelationsComponent(block)

		// then
		assert.Equal(t, "block1", params.Id)
		assert.Empty(t, params.Cells)
	})

	t.Run("no featured relations in details should return empty params", func(t *testing.T) {
		// given
		r.Sp = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{Details: &types.Struct{
				Fields: map[string]*types.Value{},
			}}},
		}
		block := &model.Block{Id: "block2"}

		// when
		params := r.MakeFeaturedRelationsComponent(block)

		// theb
		assert.Equal(t, "block2", params.Id)
		assert.Empty(t, params.Cells)
	})

	t.Run("empty featured relations list should return empty params", func(t *testing.T) {
		// given
		r.Sp = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{Details: &types.Struct{
				Fields: map[string]*types.Value{
					bundle.RelationKeyFeaturedRelations.String(): {
						Kind: &types.Value_ListValue{ListValue: &types.ListValue{Values: []*types.Value{}}},
					},
				},
			},
			}}}
		block := &model.Block{Id: "block3"}

		// when
		params := r.MakeFeaturedRelationsComponent(block)

		// then
		assert.Equal(t, "block3", params.Id)
		assert.Empty(t, params.Cells)
	})

	t.Run("valid featured relations should return populated params", func(t *testing.T) {
		// given
		r.Sp = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{Details: &types.Struct{
				Fields: map[string]*types.Value{
					bundle.RelationKeyFeaturedRelations.String(): {
						Kind: &types.Value_ListValue{
							ListValue: &types.ListValue{
								Values: []*types.Value{
									{Kind: &types.Value_StringValue{StringValue: "relation1"}},
									{Kind: &types.Value_StringValue{StringValue: "relation2"}},
								},
							},
						},
					},
					"relation1": {Kind: &types.Value_StringValue{StringValue: "Value1"}},
					"relation2": {},
				},
			}}},
		}
		sn := &pb.SnapshotWithType{
			SbType: model.SmartBlockType_STRelation,
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyUniqueKey.String(): pbtypes.String(domain.RelationKey("relation1").URL()),
						bundle.RelationKeyName.String():      pbtypes.String("Relation1"),
					},
				},
			}},
		}
		marshaler := jsonpb.Marshaler{}
		json, err := marshaler.MarshalToString(sn)
		assert.NoError(t, err)
		r.UberSp.PbFiles[filepath.Join("relations", "relation1.pb")] = json

		sn = &pb.SnapshotWithType{
			SbType: model.SmartBlockType_STRelation,
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyUniqueKey.String(): pbtypes.String(domain.RelationKey("relation2").URL()),
						bundle.RelationKeyName.String():      pbtypes.String("Relation2"),
					},
				},
			}},
		}
		json, err = marshaler.MarshalToString(sn)
		assert.NoError(t, err)
		r.UberSp.PbFiles[filepath.Join("relations", "relation2.pb")] = json

		block := &model.Block{Id: "block4"}

		// when
		params := r.MakeFeaturedRelationsComponent(block)

		// then
		assert.Equal(t, "block4", params.Id)
		assert.Len(t, params.Cells, 2)
	})
}
