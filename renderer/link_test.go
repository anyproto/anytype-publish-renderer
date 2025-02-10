package renderer

import (
	"path/filepath"
	"testing"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/util/pbtypes"

	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/assert"
)

func TestMakeLinkRenderParams(t *testing.T) {
	t.Run("target details not found", func(t *testing.T) {
		// given
		r := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
		block := &model.Block{
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "nonexistent-id",
				},
			},
		}
		expected := &LinkRenderParams{IsDeleted: true}

		// when
		result := r.MakeLinkRenderParams(block)

		// then
		assert.Equal(t, expected, result)
	})
	t.Run("deleted block", func(t *testing.T) {
		// given
		r := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
		pbFiles := map[string]*pb.SnapshotWithType{
			filepath.Join("objects", "deleted-id.pb"): {
				SbType: model.SmartBlockType_Page,
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():        pbtypes.String("deleted-id.pb"),
						bundle.RelationKeyIsDeleted.String(): pbtypes.Bool(true),
					}},
				}},
			},
		}
		r.CachedPbFiles = pbFiles
		block := &model.Block{
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "deleted-id",
				},
			},
		}
		expected := &LinkRenderParams{IsDeleted: true}

		// when
		result := r.MakeLinkRenderParams(block)

		// then
		assert.Equal(t, expected, result)
	})
	t.Run("archived block", func(t *testing.T) {
		// given
		r := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
		pbFiles := map[string]*pb.SnapshotWithType{
			filepath.Join("objects", "archived-id.pb"): {
				SbType: model.SmartBlockType_Page,
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():         pbtypes.String("archived-id"),
						bundle.RelationKeyIsArchived.String(): pbtypes.Bool(true),
						bundle.RelationKeyName.String():       pbtypes.String("Archived Block"),
						bundle.RelationKeySpaceId.String():    pbtypes.String("spaceId"),
					}},
				}},
			},
		}
		r.CachedPbFiles = pbFiles
		block := &model.Block{
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "archived-id",
				},
			},
		}
		expected := &LinkRenderParams{
			Classes:        "text isArchived",
			ContentClasses: "content",
			SidesClasses:   "sides",
			CardClasses:    "linkCard isPage c1",
			IsArchived:     "isArchived",
			Name:           "Archived Block",
			Url:            templ.SafeURL("anytype://object?objectId=archived-id&spaceId=spaceId"),
			CoverTemplate:  templ.Component(nil),
		}

		// when
		result := r.MakeLinkRenderParams(block)

		// then
		compareLinks(t, expected, result)
	})
	t.Run("block with icon emoji", func(t *testing.T) {
		// given
		r := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
		pbFiles := map[string]*pb.SnapshotWithType{
			filepath.Join("objects", "emoji-icon-id.pb"): {
				SbType: model.SmartBlockType_Page,
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():        pbtypes.String("emoji-icon-id"),
						bundle.RelationKeyName.String():      pbtypes.String("Emoji Icon Block"),
						bundle.RelationKeyIconEmoji.String(): pbtypes.String("😊"),
						bundle.RelationKeySpaceId.String():   pbtypes.String("spaceId"),
					}},
				}},
			},
		}
		r.CachedPbFiles = pbFiles
		block := &model.Block{
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					IconSize:      model.BlockContentLink_SizeMedium,
					TargetBlockId: "emoji-icon-id",
				},
			},
		}

		// when
		result := r.MakeLinkRenderParams(block)

		// then
		assert.NotNil(t, result.IconTemplate)
		assert.Equal(t, "linkCard isPage withIcon c20 c1", result.CardClasses)
	})
	t.Run("collection layout", func(t *testing.T) {
		r1 := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
		r1.CachedPbFiles = map[string]*pb.SnapshotWithType{
			filepath.Join("objects", "collection-id.pb"): {
				SbType: model.SmartBlockType_Page,
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():      pbtypes.String("collection-id"),
						bundle.RelationKeyLayout.String():  pbtypes.Float64(float64(model.ObjectType_collection)),
						bundle.RelationKeyName.String():    pbtypes.String("Collection Block"),
						bundle.RelationKeySpaceId.String(): pbtypes.String("spaceId"),
					}},
				}},
			},
		}
		result1 := r1.MakeLinkRenderParams(&model.Block{
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "collection-id",
				},
			},
		})
		compareLinks(t, &LinkRenderParams{
			Name:           "Collection Block",
			Url:            "anytype://object?objectId=collection-id&spaceId=spaceId",
			ContentClasses: "content",
			SidesClasses:   "sides",
			CardClasses:    "linkCard isCollection c1",
			Classes:        "text ",
		}, result1)
	})
	t.Run("todo layout", func(t *testing.T) {
		// given
		r1 := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
		r1.CachedPbFiles = map[string]*pb.SnapshotWithType{
			filepath.Join("objects", "todo-id.pb"): {
				SbType: model.SmartBlockType_Page,
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():      pbtypes.String("todo-id"),
						bundle.RelationKeyLayout.String():  pbtypes.Float64(float64(model.ObjectType_todo)),
						bundle.RelationKeyName.String():    pbtypes.String("Todo"),
						bundle.RelationKeySpaceId.String(): pbtypes.String("spaceId"),
					}},
				}},
			},
		}

		// when
		result1 := r1.MakeLinkRenderParams(&model.Block{
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "todo-id",
				},
			},
		})

		// then
		compareLinks(t, &LinkRenderParams{
			Name:           "Todo",
			Url:            "anytype://object?objectId=todo-id&spaceId=spaceId",
			ContentClasses: "content",
			SidesClasses:   "sides",
			CardClasses:    "linkCard isTask c1",
			Classes:        "text ",
			IconTemplate:   NoneTemplate(""),
		}, result1)
	})
	t.Run("todo layout, checkbox set", func(t *testing.T) {
		// given
		r1 := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
		r1.CachedPbFiles = map[string]*pb.SnapshotWithType{
			filepath.Join("objects", "todo-id.pb"): {
				SbType: model.SmartBlockType_Page,
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():      pbtypes.String("todo-id"),
						bundle.RelationKeyLayout.String():  pbtypes.Float64(float64(model.ObjectType_todo)),
						bundle.RelationKeyName.String():    pbtypes.String("Todo"),
						bundle.RelationKeyDone.String():    pbtypes.Bool(true),
						bundle.RelationKeySpaceId.String(): pbtypes.String("spaceId"),
					}},
				}},
			},
		}

		// when
		result1 := r1.MakeLinkRenderParams(&model.Block{
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "todo-id",
				},
			},
		})

		// then
		compareLinks(t, &LinkRenderParams{
			Name:           "Todo",
			Url:            "anytype://object?objectId=todo-id&spaceId=spaceId",
			ContentClasses: "content",
			SidesClasses:   "sides",
			CardClasses:    "linkCard isTask c1",
			Classes:        "text ",
			IconTemplate:   NoneTemplate(""),
		}, result1)
	})
	t.Run("block with description", func(t *testing.T) {
		// given
		r1 := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
		r1.CachedPbFiles = map[string]*pb.SnapshotWithType{
			filepath.Join("objects", "test-id.pb"): {
				SbType: model.SmartBlockType_Page,
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():          pbtypes.String("test-id"),
						bundle.RelationKeyLayout.String():      pbtypes.Float64(float64(model.ObjectType_profile)),
						bundle.RelationKeyName.String():        pbtypes.String("Test"),
						bundle.RelationKeyDescription.String(): pbtypes.String("description"),
						bundle.RelationKeySpaceId.String():     pbtypes.String("spaceId"),
					}},
				}},
			},
		}

		// when
		result1 := r1.MakeLinkRenderParams(&model.Block{
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "test-id",
					Description:   model.BlockContentLink_Added,
				},
			},
		})

		// then
		compareLinks(t, &LinkRenderParams{
			Name:           "Test",
			Description:    "description",
			Url:            "anytype://object?objectId=test-id&spaceId=spaceId",
			Classes:        "text ",
			ContentClasses: "content",
			SidesClasses:   "sides",
			CardClasses:    "linkCard isHuman c2",
			IconTemplate:   NoneTemplate(""),
		}, result1)
	})
	t.Run("block with snippet", func(t *testing.T) {
		// given
		r1 := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
		r1.CachedPbFiles = map[string]*pb.SnapshotWithType{
			filepath.Join("objects", "test-id.pb"): {
				SbType: model.SmartBlockType_Page,
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():      pbtypes.String("test-id"),
						bundle.RelationKeyLayout.String():  pbtypes.Float64(float64(model.ObjectType_participant)),
						bundle.RelationKeyName.String():    pbtypes.String("Test"),
						bundle.RelationKeySnippet.String(): pbtypes.String("snippet"),
						bundle.RelationKeySpaceId.String(): pbtypes.String("spaceId"),
					}},
				}},
			},
		}

		// when
		result1 := r1.MakeLinkRenderParams(&model.Block{
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "test-id",
					Description:   model.BlockContentLink_Content,
					CardStyle:     model.BlockContentLink_Card,
				},
			},
		})

		// then
		compareLinks(t, &LinkRenderParams{
			Name:           "Test",
			Description:    "snippet",
			Url:            "anytype://object?objectId=test-id&spaceId=spaceId",
			Classes:        "card ",
			ContentClasses: "content",
			SidesClasses:   "sides",
			CardClasses:    "linkCard isParticipant c2",
			IconTemplate:   NoneTemplate(""),
		}, result1)
	})
	t.Run("block with cover", func(t *testing.T) {
		// given
		r1 := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
		r1.CachedPbFiles = map[string]*pb.SnapshotWithType{
			filepath.Join("objects", "test-id.pb"): {
				SbType: model.SmartBlockType_Page,
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():        pbtypes.String("test-id"),
						bundle.RelationKeyLayout.String():    pbtypes.Float64(float64(model.ObjectType_set)),
						bundle.RelationKeyName.String():      pbtypes.String("Test"),
						bundle.RelationKeyCoverType.String(): pbtypes.Int64(2),
						bundle.RelationKeyCoverId.String():   pbtypes.String("gray"),
						bundle.RelationKeySpaceId.String():   pbtypes.String("spaceId"),
					}},
				}},
			},
		}

		// when
		result1 := r1.MakeLinkRenderParams(&model.Block{
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "test-id",
					Relations:     []string{"cover"},
				},
			},
		})

		// then
		assert.NotNil(t, result1.CoverTemplate)
	})
	t.Run("block with type", func(t *testing.T) {
		// given
		r1 := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
		r1.CachedPbFiles = map[string]*pb.SnapshotWithType{
			filepath.Join("objects", "test-id.pb"): {
				SbType: model.SmartBlockType_Page,
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():      pbtypes.String("test-id"),
						bundle.RelationKeyLayout.String():  pbtypes.Float64(float64(model.ObjectType_set)),
						bundle.RelationKeyName.String():    pbtypes.String("Test"),
						bundle.RelationKeyType.String():    pbtypes.String("type"),
						bundle.RelationKeySpaceId.String(): pbtypes.String("spaceId")},
					}},
				}},
			filepath.Join("types", "type.pb"): {
				SbType: model.SmartBlockType_Page,
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():   pbtypes.String("type"),
						bundle.RelationKeyName.String(): pbtypes.String("Type")},
					}},
				}},
		}

		// when
		result1 := r1.MakeLinkRenderParams(&model.Block{
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "test-id",
					Relations:     []string{"type"},
				},
			},
		})

		// then
		compareLinks(t, &LinkRenderParams{
			Name:           "Test",
			Type:           "Type",
			Url:            "anytype://object?objectId=test-id&spaceId=spaceId",
			Classes:        "text ",
			ContentClasses: "content",
			SidesClasses:   "sides",
			CardClasses:    "linkCard isSet c2",
			IconTemplate:   NoneTemplate(""),
		}, result1)
	})
}

func compareLinks(t *testing.T, expected *LinkRenderParams, result *LinkRenderParams) bool {
	return assert.Equal(t, expected.Classes, result.Classes) &&
		assert.Equal(t, expected.ContentClasses, result.ContentClasses) &&
		assert.Equal(t, expected.SidesClasses, result.SidesClasses) &&
		assert.Equal(t, expected.Name, result.Name) &&
		assert.Equal(t, expected.Url, result.Url) &&
		assert.Equal(t, expected.Description, result.Description) &&
		assert.Equal(t, expected.Type, result.Type) &&
		assert.Equal(t, expected.IsDeleted, result.IsDeleted) &&
		assert.Equal(t, expected.IsArchived, result.IsArchived)
}
