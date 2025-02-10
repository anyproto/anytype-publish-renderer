package renderer

import (
	"path/filepath"
	"testing"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/assert"

	"github.com/anyproto/anytype-heart/core/domain"
	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
)

func TestMakeFeaturedRelationsComponent(t *testing.T) {
	t.Run("empty details", func(t *testing.T) {
		// given
		r := &Renderer{Sp: &pb.SnapshotWithType{Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{}}}}

		// when
		result := r.MakeFeaturedRelationsComponent()

		// then
		assert.Nil(t, result)
	})

	t.Run("no featured relations", func(t *testing.T) {
		// given
		r := &Renderer{Sp: &pb.SnapshotWithType{Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{Details: &types.Struct{Fields: map[string]*types.Value{}}}}}}

		// when
		result := r.MakeFeaturedRelationsComponent()

		// then
		assert.Nil(t, result)
	})
	t.Run("no featured relations", func(t *testing.T) {
		// given
		r := &Renderer{Sp: &pb.SnapshotWithType{Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{Details: &types.Struct{Fields: map[string]*types.Value{
			bundle.RelationKeyFeaturedRelations.String(): pbtypes.StringList(nil),
		}}}}}}

		// when
		result := r.MakeFeaturedRelationsComponent()

		// then
		assert.Nil(t, result)
	})
	t.Run("skip backlinks and links", func(t *testing.T) {
		r := &Renderer{UberSp: &PublishingUberSnapshot{PbFiles: make(map[string]string)}}

		r.Sp = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{Details: &types.Struct{
				Fields: map[string]*types.Value{
					bundle.RelationKeyFeaturedRelations.String(): {
						Kind: &types.Value_ListValue{
							ListValue: &types.ListValue{
								Values: []*types.Value{
									{Kind: &types.Value_StringValue{StringValue: bundle.RelationKeyBacklinks.String()}},
									{Kind: &types.Value_StringValue{StringValue: bundle.RelationKeyLinks.String()}},
								},
							},
						},
					},
					bundle.RelationKeyBacklinks.String(): {},
					bundle.RelationKeyLinks.String():     {},
				},
			}}},
		}

		// when
		result := r.MakeFeaturedRelationsComponent()

		// then
		assert.Nil(t, result)
	})
	t.Run("featured relations", func(t *testing.T) {
		r := &Renderer{UberSp: &PublishingUberSnapshot{PbFiles: make(map[string]string)}}

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

		// when
		result := r.MakeFeaturedRelationsComponent()

		// then
		assert.NotNil(t, result)
	})
	t.Run("empty list relation", func(t *testing.T) {
		// given
		relationKey := "tag-relation"
		r := &Renderer{UberSp: &PublishingUberSnapshot{PbFiles: make(map[string]string)}}

		r.Sp = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{Details: &types.Struct{
				Fields: map[string]*types.Value{
					bundle.RelationKeyFeaturedRelations.String(): pbtypes.StringList([]string{relationKey}),
					relationKey: pbtypes.StringList([]string{}),
				},
			}}},
		}

		sn := &pb.SnapshotWithType{
			SbType: model.SmartBlockType_STRelation,
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyUniqueKey.String():      pbtypes.String(domain.RelationKey("tag-relation").URL()),
						bundle.RelationKeyName.String():           pbtypes.String("Tag Relation"),
						bundle.RelationKeyRelationFormat.String(): pbtypes.Int64(int64(model.RelationFormat_tag)),
					},
				},
			}},
		}
		marshaler := jsonpb.Marshaler{}
		json, err := marshaler.MarshalToString(sn)
		assert.NoError(t, err)
		r.UberSp.PbFiles[filepath.Join("relations", "tag-relation.pb")] = json

		// when
		result := r.MakeFeaturedRelationsComponent()

		// then
		assert.Nil(t, result)
	})
}
