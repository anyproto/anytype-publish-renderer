package renderer

import (
	"path/filepath"
	"testing"

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
		expected := &BlockParams{Classes: []string{"block", "align0", "blockLink", "withIcon", "c20"}}
		block := &model.Block{
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "nonexistent-id",
				},
			},
		}
		actual := r.makeLinkBlockParams(block)
		assert.Equal(t, expected.Classes, actual.Classes)
		assert.Equal(t, expected.ContentClasses, actual.ContentClasses)

		tag, err := blockParamsToHtmlTag(actual)
		assert.NoError(t, err)

		pathAssertions := []pathAssertion{
			{"div.deleted > div.iconObject.withDefault.c20 > img.iconCommon.c18 > attrs[src]", "/static/img/icon/ghost.svg"},
			{"div.deleted > div.name > Content", "Non-existent object"},
		}
		assertHtmlTag(t, tag, pathAssertions)

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

		expected := &BlockParams{Classes: []string{"block", "align0", "blockLink", "withIcon", "c20"}}
		block := &model.Block{
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "deleted-id",
				},
			},
		}
		actual := r.makeLinkBlockParams(block)
		assert.Equal(t, expected.Classes, actual.Classes)
		assert.Equal(t, expected.ContentClasses, actual.ContentClasses)

		tag, err := blockParamsToHtmlTag(actual)
		assert.NoError(t, err)

		pathAssertions := []pathAssertion{
			{"div.deleted > div.iconObject.withDefault.c20 > img.iconCommon.c18 > attrs[src]", "/static/img/icon/ghost.svg"},
			{"div.deleted > div.name > Content", "Non-existent object"},
		}
		assertHtmlTag(t, tag, pathAssertions)
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

		expected := &BlockParams{Classes: []string{"block", "align0", "blockLink", "text", "isArchived"}}
		block := &model.Block{
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "archived-id",
				},
			},
		}

		actual := r.makeLinkBlockParams(block)
		assert.Equal(t, expected.Classes, actual.Classes)
		assert.Equal(t, expected.ContentClasses, actual.ContentClasses)

		tag, err := blockParamsToHtmlTag(actual)
		assert.NoError(t, err)

		pathAssertions := []pathAssertion{
			{"a.linkCard.isPage.c1 > attrs[href]", "anytype://object?objectId=archived-id&spaceId=spaceId"},
			{"a.linkCard.isPage.c1 > div.sides > div.side.left > div.cardName > div.name > Content", "Archived Block"},
			{"a.linkCard.isPage.c1 > div.sides > div.side.left > div.cardName > div.tagItem.isMultiSelect.archive > Content", "Deleted"},
		}
		assertHtmlTag(t, tag, pathAssertions)

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

		expected := &BlockParams{Classes: []string{"block", "align0", "blockLink", "text"}}
		actual := r.makeLinkBlockParams(block)
		assert.Equal(t, expected.Classes, actual.Classes)
		assert.Equal(t, expected.ContentClasses, actual.ContentClasses)

		tag, err := blockParamsToHtmlTag(actual)
		assert.NoError(t, err)

		pathAssertions := []pathAssertion{
			{"a.linkCard.isPage.withIcon.c20.c1 > attrs[href]", "anytype://object?objectId=emoji-icon-id&spaceId=spaceId"},
			{"a > div.sides > div.side.left > div.cardName > div.iconObject.c20 > img.smileImage.c20 > attrs[src]", "https://anytype-static.fra1.cdn.digitaloceanspaces.com/emojies/1f60a.png"},
			{"a > div.sides > div.side.left > div.cardName > div.name > Content", "Emoji Icon Block"},
		}
		assertHtmlTag(t, tag, pathAssertions)

	})
	t.Run("collection layout", func(t *testing.T) {
		// given
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

		// when
		block := &model.Block{
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "collection-id",
				},
			},
		}
		expected := &BlockParams{Classes: []string{"block", "align0", "blockLink", "text"}}
		actual := r1.makeLinkBlockParams(block)
		assert.Equal(t, expected.Classes, actual.Classes)
		assert.Equal(t, expected.ContentClasses, actual.ContentClasses)

		tag, err := blockParamsToHtmlTag(actual)
		assert.NoError(t, err)

		pathAssertions := []pathAssertion{
			{"a.linkCard.isCollection.c1 > attrs[href]", "anytype://object?objectId=collection-id&spaceId=spaceId"},
			{"a.linkCard.isCollection.c1 > div.sides > div.side.left > div.cardName > div.name > Content", "Collection Block"}}

		assertHtmlTag(t, tag, pathAssertions)

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

		block := &model.Block{
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "todo-id",
				},
			},
		}
		expected := &BlockParams{Classes: []string{"block", "align0", "blockLink", "text"}}
		actual := r1.makeLinkBlockParams(block)
		assert.Equal(t, expected.Classes, actual.Classes)
		assert.Equal(t, expected.ContentClasses, actual.ContentClasses)

		tag, err := blockParamsToHtmlTag(actual)
		assert.NoError(t, err)

		pathAssertions := []pathAssertion{
			{"a.linkCard.isTask.c1 > attrs[href]", "anytype://object?objectId=todo-id&spaceId=spaceId"},
			{"a.linkCard.isTask.c1 > div.sides > div.side.left > div.cardName > div.name > Content", "Todo"},
		}
		assertHtmlTag(t, tag, pathAssertions)
	})

	t.Run("todo layout, checkbox set", func(t *testing.T) {
		// given
		r1 := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
		r1.CachedPbFiles = map[string]*pb.SnapshotWithType{
			filepath.Join("objects", "todo-id.pb"): {
				SbType: model.SmartBlockType_Page,
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():     pbtypes.String("todo-id"),
						bundle.RelationKeyLayout.String(): pbtypes.Float64(float64(model.ObjectType_todo)),
						bundle.RelationKeyName.String():   pbtypes.String("Todo"),
						// TODO: same test, nothing changed with this relation enabled?
						bundle.RelationKeyDone.String():    pbtypes.Bool(true),
						bundle.RelationKeySpaceId.String(): pbtypes.String("spaceId"),
					}},
				}},
			},
		}

		block := &model.Block{
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "todo-id",
				},
			},
		}
		expected := &BlockParams{Classes: []string{"block", "align0", "blockLink", "text"}}
		actual := r1.makeLinkBlockParams(block)
		assert.Equal(t, expected.Classes, actual.Classes)
		assert.Equal(t, expected.ContentClasses, actual.ContentClasses)

		tag, err := blockParamsToHtmlTag(actual)
		assert.NoError(t, err)

		pathAssertions := []pathAssertion{
			{"a.linkCard.isTask.c1 > attrs[href]", "anytype://object?objectId=todo-id&spaceId=spaceId"},
			{"a.linkCard.isTask.c1 > div.sides > div.side.left > div.cardName > div.name > Content", "Todo"},
		}
		assertHtmlTag(t, tag, pathAssertions)

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
		result1 := r1.makeLinkBlockParams(&model.Block{
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "test-id",
					Description:   model.BlockContentLink_Added,
				},
			},
		})

		// then
		expectedHtml := `<a href="anytype://object?objectId=test-id&amp;spaceId=spaceId" class="linkCard isHuman c2"><div class="sides"><div class="side left"><div class="cardName"><div class="name">Test</div></div><div class="relationItem cardDescription"><div class="description">description</div></div></div></div></a>`
		compareLinks(t, &BlockParams{
			Classes: []string{"block", "align0", "blockLink", "text"},
		}, result1, expectedHtml)
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
		result1 := r1.makeLinkBlockParams(&model.Block{
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "test-id",
					Description:   model.BlockContentLink_Content,
					CardStyle:     model.BlockContentLink_Card,
				},
			},
		})

		// then
		expectedHtml := `<a href="anytype://object?objectId=test-id&amp;spaceId=spaceId" class="linkCard isParticipant c2"><div class="sides"><div class="side left"><div class="cardName"><div class="name">Test</div></div><div class="relationItem cardDescription"><div class="description">snippet</div></div></div></div></a>`
		compareLinks(t, &BlockParams{
			Classes: []string{"block", "align0", "blockLink", "card"},
		}, result1, expectedHtml)
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
		result1 := r1.makeLinkBlockParams(&model.Block{
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "test-id",
					Relations:     []string{"cover"},
				},
			},
		})

		// then
		expectedHtml := `<a href="anytype://object?objectId=test-id&amp;spaceId=spaceId" class="linkCard isSet withCover c1"><div class="sides"><div class="side left"><div class="cardName"><div class="name">Test</div></div></div><div class="side right"><div class="cover type2 gray" style="background-position:0% 0%;background-size:100%;"></div></div></div></a>`
		compareLinks(t, &BlockParams{
			Classes: []string{"block", "align0", "blockLink", "text"},
		}, result1, expectedHtml)
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
		result1 := r1.makeLinkBlockParams(&model.Block{
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "test-id",
					Relations:     []string{"type"},
				},
			},
		})

		// then
		expectedHtml := `<a href="anytype://object?objectId=test-id&amp;spaceId=spaceId" class="linkCard isSet c2"><div class="sides"><div class="side left"><div class="cardName"><div class="name">Test</div></div><div class="relationItem cardType"><div class="item">Type</div></div></div></div></a>`
		compareLinks(t, &BlockParams{
			Classes: []string{"block", "align0", "blockLink", "text"},
		}, result1, expectedHtml)
	})
}
func compareLinks(t *testing.T, expected *BlockParams, result *BlockParams, expectedHtml string) {
	assert.Equal(t, true, false)
}
