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
		result := r.makeFeaturedRelationsComponent()

		// then
		assert.Nil(t, result)
	})

	t.Run("no featured relations", func(t *testing.T) {
		// given
		r := &Renderer{Sp: &pb.SnapshotWithType{Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{Details: &types.Struct{Fields: map[string]*types.Value{}}}}}}

		// when
		result := r.makeFeaturedRelationsComponent()

		// then
		assert.Nil(t, result)
	})
	t.Run("no featured relations", func(t *testing.T) {
		// given
		r := &Renderer{Sp: &pb.SnapshotWithType{Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{Details: &types.Struct{Fields: map[string]*types.Value{
			bundle.RelationKeyFeaturedRelations.String(): pbtypes.StringList(nil),
		}}}}}}

		// when
		result := r.makeFeaturedRelationsComponent()

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
		result := r.makeFeaturedRelationsComponent()

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
		result := r.makeFeaturedRelationsComponent()

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
		result := r.makeFeaturedRelationsComponent()

		// then
		assert.Nil(t, result)
	})
	t.Run("featured relations are from object type", func(t *testing.T) {
		// given
		r := NewTestRenderer(
			WithRootSnapshot(&pb.SnapshotWithType{
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyType.String(): pbtypes.String("objectType"),
						"relation":                      pbtypes.String("test"),
					},
				}}}}),
			WithLinkedSnapshot(t, filepath.Join("relations", "bafyreihja7bgkhxjhcan26ts44qqbjoxl4sr5ckqoxdlty4edlrueoylj4.pb"), &pb.SnapshotWithType{
				SbType: model.SmartBlockType_STRelation,
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{
						Fields: map[string]*types.Value{
							bundle.RelationKeyUniqueKey.String():      pbtypes.String(domain.RelationKey("relation").URL()),
							bundle.RelationKeyName.String():           pbtypes.String("Relation"),
							bundle.RelationKeyRelationFormat.String(): pbtypes.Int64(int64(model.RelationFormat_longtext)),
							bundle.RelationKeyId.String():             pbtypes.String("bafyreihja7bgkhxjhcan26ts44qqbjoxl4sr5ckqoxdlty4edlrueoylj4"),
						},
					},
				}},
			}),
			WithObjectTypeDetails(&types.Struct{
				Fields: map[string]*types.Value{
					bundle.RelationKeyName.String():                         pbtypes.String("Type"),
					bundle.RelationKeyId.String():                           pbtypes.String("objectType"),
					bundle.RelationKeyResolvedLayout.String():               pbtypes.Int64(int64(model.ObjectType_objectType)),
					bundle.RelationKeyRecommendedFeaturedRelations.String(): pbtypes.StringList([]string{"bafyreihja7bgkhxjhcan26ts44qqbjoxl4sr5ckqoxdlty4edlrueoylj4"}),
				},
			}),
		)

		// when
		result := r.makeFeaturedRelationsBlockParams(&model.Block{
			Content: &model.BlockContentOfFeaturedRelations{FeaturedRelations: &model.BlockContentFeaturedRelations{}},
		})

		// then
		assert.NotNil(t, result)
		tag, err := blockParamsToHtmlTag(result)
		assert.NoError(t, err)
		pathAssertions := []pathAssertion{
			{"div.wrap > div > attrs[class]", "cell last c-longText"},
			{"div.wrap > div.cell.last.c-longText > div.cellContent.last.c-longText > div.name > Content", "test"}}
		assertHtmlTag(t, tag, pathAssertions)
	})
	t.Run("featured relations have only description", func(t *testing.T) {
		// given
		r := NewTestRenderer(
			WithRootSnapshot(&pb.SnapshotWithType{
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{Details: &types.Struct{
					Fields: map[string]*types.Value{
						bundle.RelationKeyFeaturedRelations.String(): pbtypes.StringList([]string{bundle.RelationKeyDescription.String()}),
					},
				}}}}),
			WithLinkedSnapshot(t, filepath.Join("relations", "bafyreihja7bgkhxjhcan26ts44qqbjoxl4sr5ckqoxdlty4edlrueoylj4.pb"), &pb.SnapshotWithType{
				SbType: model.SmartBlockType_STRelation,
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{
						Fields: map[string]*types.Value{
							bundle.RelationKeyUniqueKey.String():      pbtypes.String(domain.RelationKey("tag-relation").URL()),
							bundle.RelationKeyName.String():           pbtypes.String("Tag Relation"),
							bundle.RelationKeyRelationFormat.String(): pbtypes.Int64(int64(model.RelationFormat_tag)),
							bundle.RelationKeyId.String():             pbtypes.String("bafyreihja7bgkhxjhcan26ts44qqbjoxl4sr5ckqoxdlty4edlrueoylj4"),
						},
					},
				}},
			}),
			WithObjectTypeDetails(&types.Struct{
				Fields: map[string]*types.Value{
					bundle.RelationKeyName.String():                         pbtypes.String("Type"),
					bundle.RelationKeyId.String():                           pbtypes.String("objectType"),
					bundle.RelationKeyResolvedLayout.String():               pbtypes.Int64(int64(model.ObjectType_objectType)),
					bundle.RelationKeyRecommendedFeaturedRelations.String(): pbtypes.StringList([]string{"bafyreihja7bgkhxjhcan26ts44qqbjoxl4sr5ckqoxdlty4edlrueoylj4"}),
				},
			}),
		)

		// when
		result := r.makeFeaturedRelationsBlockParams(&model.Block{
			Content: &model.BlockContentOfFeaturedRelations{FeaturedRelations: &model.BlockContentFeaturedRelations{}},
		})

		// then
		assert.NotNil(t, result)
		tag, err := blockParamsToHtmlTag(result)
		assert.NoError(t, err)
		pathAssertions := []pathAssertion{
			{"div.wrap > div > attrs[class]", "cell last isEmpty"},
			{"div.wrap > div.cell.last.isEmpty > div.cellContent.last.isEmpty > div > attrs[class]", "empty"}}
		assertHtmlTag(t, tag, pathAssertions)
	})
}
