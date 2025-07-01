package renderer

import (
	"path/filepath"
	"testing"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/assert"
)

// TestGetLinkObjectIds tests the GetLinkObjectIds function which:
// 1. Only processes blocks with BlockContentOfLink type (skips text, file, etc.)
// 2. Returns target object IDs only for links with model.ObjectType_basic layout
// 3. Includes missing/invalid targets (they default to basic layout)
// 4. Deduplicates results (first occurrence wins)
// 5. Processes blocks in the order they appear in Root.ChildrenIds

func TestGetLinkObjectIds(t *testing.T) {
	t.Run("empty root children returns empty slice", func(t *testing.T) {
		// given
		rootBlock := &model.Block{
			Id:          "root",
			ChildrenIds: []string{},
		}
		r := NewTestRenderer(
			WithBlocksById(map[string]*model.Block{
				"root": rootBlock,
			}),
		)
		r.Root = rootBlock

		// when
		result := r.GetLinkObjectIds()

		// then
		assert.Empty(t, result)
	})

	t.Run("non-existent child blocks are skipped", func(t *testing.T) {
		// given
		rootBlock := &model.Block{
			Id:          "root",
			ChildrenIds: []string{"missing-block"},
		}
		r := NewTestRenderer(
			WithBlocksById(map[string]*model.Block{
				"root": rootBlock,
			}),
		)
		r.Root = rootBlock

		// when
		result := r.GetLinkObjectIds()

		// then
		assert.Empty(t, result)
	})

	t.Run("nil child blocks are skipped", func(t *testing.T) {
		// given
		rootBlock := &model.Block{
			Id:          "root",
			ChildrenIds: []string{"nil-block"},
		}
		r := NewTestRenderer(
			WithBlocksById(map[string]*model.Block{
				"root":      rootBlock,
				"nil-block": nil,
			}),
		)
		r.Root = rootBlock

		// when
		result := r.GetLinkObjectIds()

		// then
		assert.Empty(t, result)
	})

	t.Run("non-link blocks are properly skipped", func(t *testing.T) {
		// given
		textBlock := &model.Block{
			Id: "text-block",
			Content: &model.BlockContentOfText{
				Text: &model.BlockContentText{
					Text: "some text",
				},
			},
		}
		rootBlock := &model.Block{
			Id:          "root",
			ChildrenIds: []string{"text-block"},
		}
		r := NewTestRenderer(
			WithBlocksById(map[string]*model.Block{
				"root":       rootBlock,
				"text-block": textBlock,
			}),
		)
		r.Root = rootBlock

		// when
		result := r.GetLinkObjectIds()

		// then
		assert.Empty(t, result)
	})

	t.Run("link blocks with no target details are included due to default layout", func(t *testing.T) {
		// given
		linkBlock := &model.Block{
			Id: "link-block",
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "missing-target",
				},
			},
		}
		rootBlock := &model.Block{
			Id:          "root",
			ChildrenIds: []string{"link-block"},
		}
		r := NewTestRenderer(
			WithBlocksById(map[string]*model.Block{
				"root":       rootBlock,
				"link-block": linkBlock,
			}),
		)
		r.Root = rootBlock

		// when
		result := r.GetLinkObjectIds()

		// then
		// Note: Missing targets default to basic layout and are included
		assert.Equal(t, []string{"missing-target"}, result)
	})

	t.Run("link blocks with non-basic layout are skipped", func(t *testing.T) {
		// given
		linkBlock := &model.Block{
			Id: "link-block",
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "target-id",
				},
			},
		}
		rootBlock := &model.Block{
			Id:          "root",
			ChildrenIds: []string{"link-block"},
		}
		r := NewTestRenderer(
			WithBlocksById(map[string]*model.Block{
				"root":       rootBlock,
				"link-block": linkBlock,
			}),
			WithCachedPbFiles(map[string]*pb.SnapshotWithType{
				filepath.Join("objects", "target-id.pb"): {
					SbType: model.SmartBlockType_Page,
					Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
						Details: &types.Struct{Fields: map[string]*types.Value{
							bundle.RelationKeyId.String():     pbtypes.String("target-id"),
							bundle.RelationKeyLayout.String(): pbtypes.Float64(float64(model.ObjectType_collection)),
						}},
					}},
				},
			}),
		)
		r.Root = rootBlock

		// when
		result := r.GetLinkObjectIds()

		// then
		assert.Empty(t, result)
	})

	t.Run("single link block with basic layout returns target id", func(t *testing.T) {
		// given
		linkBlock := &model.Block{
			Id: "link-block",
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "target-id",
				},
			},
		}
		rootBlock := &model.Block{
			Id:          "root",
			ChildrenIds: []string{"link-block"},
		}
		r := NewTestRenderer(
			WithBlocksById(map[string]*model.Block{
				"root":       rootBlock,
				"link-block": linkBlock,
			}),
			WithCachedPbFiles(map[string]*pb.SnapshotWithType{
				filepath.Join("objects", "target-id.pb"): {
					SbType: model.SmartBlockType_Page,
					Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
						Details: &types.Struct{Fields: map[string]*types.Value{
							bundle.RelationKeyId.String():     pbtypes.String("target-id"),
							bundle.RelationKeyLayout.String(): pbtypes.Float64(float64(model.ObjectType_basic)),
						}},
					}},
				},
			}),
		)
		r.Root = rootBlock

		// when
		result := r.GetLinkObjectIds()

		// then
		assert.Equal(t, []string{"target-id"}, result)
	})

	t.Run("multiple link blocks with basic layout returns all target ids", func(t *testing.T) {
		// given
		linkBlock1 := &model.Block{
			Id: "link-block-1",
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "target-id-1",
				},
			},
		}
		linkBlock2 := &model.Block{
			Id: "link-block-2",
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "target-id-2",
				},
			},
		}
		rootBlock := &model.Block{
			Id:          "root",
			ChildrenIds: []string{"link-block-1", "link-block-2"},
		}
		r := NewTestRenderer(
			WithBlocksById(map[string]*model.Block{
				"root":         rootBlock,
				"link-block-1": linkBlock1,
				"link-block-2": linkBlock2,
			}),
			WithCachedPbFiles(map[string]*pb.SnapshotWithType{
				filepath.Join("objects", "target-id-1.pb"): {
					SbType: model.SmartBlockType_Page,
					Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
						Details: &types.Struct{Fields: map[string]*types.Value{
							bundle.RelationKeyId.String():     pbtypes.String("target-id-1"),
							bundle.RelationKeyLayout.String(): pbtypes.Float64(float64(model.ObjectType_basic)),
						}},
					}},
				},
				filepath.Join("objects", "target-id-2.pb"): {
					SbType: model.SmartBlockType_Page,
					Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
						Details: &types.Struct{Fields: map[string]*types.Value{
							bundle.RelationKeyId.String():     pbtypes.String("target-id-2"),
							bundle.RelationKeyLayout.String(): pbtypes.Float64(float64(model.ObjectType_basic)),
						}},
					}},
				},
			}),
		)
		r.Root = rootBlock

		// when
		result := r.GetLinkObjectIds()

		// then
		assert.Equal(t, []string{"target-id-1", "target-id-2"}, result)
	})

	t.Run("mixed block types filters correctly", func(t *testing.T) {
		// given
		linkBlockBasic := &model.Block{
			Id: "link-block-basic",
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "basic-target",
				},
			},
		}
		linkBlockCollection := &model.Block{
			Id: "link-block-collection",
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "collection-target",
				},
			},
		}
		textBlock := &model.Block{
			Id: "text-block",
			Content: &model.BlockContentOfText{
				Text: &model.BlockContentText{Text: "text"},
			},
		}
		rootBlock := &model.Block{
			Id:          "root",
			ChildrenIds: []string{"link-block-basic", "link-block-collection", "text-block"},
		}
		r := NewTestRenderer(
			WithBlocksById(map[string]*model.Block{
				"root":                  rootBlock,
				"link-block-basic":      linkBlockBasic,
				"link-block-collection": linkBlockCollection,
				"text-block":            textBlock,
			}),
			WithCachedPbFiles(map[string]*pb.SnapshotWithType{
				filepath.Join("objects", "basic-target.pb"): {
					SbType: model.SmartBlockType_Page,
					Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
						Details: &types.Struct{Fields: map[string]*types.Value{
							bundle.RelationKeyId.String():     pbtypes.String("basic-target"),
							bundle.RelationKeyLayout.String(): pbtypes.Float64(float64(model.ObjectType_basic)),
						}},
					}},
				},
				filepath.Join("objects", "collection-target.pb"): {
					SbType: model.SmartBlockType_Page,
					Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
						Details: &types.Struct{Fields: map[string]*types.Value{
							bundle.RelationKeyId.String():     pbtypes.String("collection-target"),
							bundle.RelationKeyLayout.String(): pbtypes.Float64(float64(model.ObjectType_collection)),
						}},
					}},
				},
			}),
		)
		r.Root = rootBlock

		// when
		result := r.GetLinkObjectIds()

		// then
		// Only link blocks with basic layout are included, text blocks are properly filtered out
		assert.Equal(t, []string{"basic-target"}, result)
	})

	t.Run("link block with nil link content is skipped", func(t *testing.T) {
		// given
		linkBlock := &model.Block{
			Id: "link-block",
			Content: &model.BlockContentOfLink{
				Link: nil,
			},
		}
		rootBlock := &model.Block{
			Id:          "root",
			ChildrenIds: []string{"link-block"},
		}
		r := NewTestRenderer(
			WithBlocksById(map[string]*model.Block{
				"root":       rootBlock,
				"link-block": linkBlock,
			}),
		)
		r.Root = rootBlock

		// when
		result := r.GetLinkObjectIds()

		assert.Len(t, result, 0)
	})

	// TODO: check that page layout is actually "basic"
	t.Run("target with no layout field defaults to basic and is included", func(t *testing.T) {
		// given
		linkBlock := &model.Block{
			Id: "link-block",
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "target-id",
				},
			},
		}
		rootBlock := &model.Block{
			Id:          "root",
			ChildrenIds: []string{"link-block"},
		}
		r := NewTestRenderer(
			WithBlocksById(map[string]*model.Block{
				"root":       rootBlock,
				"link-block": linkBlock,
			}),
			WithCachedPbFiles(map[string]*pb.SnapshotWithType{
				filepath.Join("objects", "target-id.pb"): {
					SbType: model.SmartBlockType_Page,
					Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
						Details: &types.Struct{Fields: map[string]*types.Value{
							bundle.RelationKeyId.String(): pbtypes.String("target-id"),
						}},
					}},
				},
			}),
		)
		r.Root = rootBlock

		// when
		result := r.GetLinkObjectIds()

		// then
		assert.Equal(t, []string{"target-id"}, result)
	})

	t.Run("duplicate target ids are deduplicated", func(t *testing.T) {
		// given
		linkBlock1 := &model.Block{
			Id: "link-block-1",
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "same-target",
				},
			},
		}
		linkBlock2 := &model.Block{
			Id: "link-block-2",
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{
					TargetBlockId: "same-target",
				},
			},
		}
		rootBlock := &model.Block{
			Id:          "root",
			ChildrenIds: []string{"link-block-1", "link-block-2"},
		}
		r := NewTestRenderer(
			WithBlocksById(map[string]*model.Block{
				"root":         rootBlock,
				"link-block-1": linkBlock1,
				"link-block-2": linkBlock2,
			}),
			WithCachedPbFiles(map[string]*pb.SnapshotWithType{
				filepath.Join("objects", "same-target.pb"): {
					SbType: model.SmartBlockType_Page,
					Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
						Details: &types.Struct{Fields: map[string]*types.Value{
							bundle.RelationKeyId.String():     pbtypes.String("same-target"),
							bundle.RelationKeyLayout.String(): pbtypes.Float64(float64(model.ObjectType_basic)),
						}},
					}},
				},
			}),
		)
		r.Root = rootBlock

		// when
		result := r.GetLinkObjectIds()

		// then
		assert.Equal(t, []string{"same-target"}, result)
	})

	t.Run("multiple layouts mixed correctly", func(t *testing.T) {
		// given
		linkBasic := &model.Block{
			Id: "link-basic",
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{TargetBlockId: "basic-target"},
			},
		}
		linkSet := &model.Block{
			Id: "link-set",
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{TargetBlockId: "set-target"},
			},
		}
		linkTodo := &model.Block{
			Id: "link-todo",
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{TargetBlockId: "todo-target"},
			},
		}
		linkProfile := &model.Block{
			Id: "link-profile",
			Content: &model.BlockContentOfLink{
				Link: &model.BlockContentLink{TargetBlockId: "profile-target"},
			},
		}
		rootBlock := &model.Block{
			Id:          "root",
			ChildrenIds: []string{"link-basic", "link-set", "link-todo", "link-profile"},
		}

		r := NewTestRenderer(
			WithBlocksById(map[string]*model.Block{
				"root":         rootBlock,
				"link-basic":   linkBasic,
				"link-set":     linkSet,
				"link-todo":    linkTodo,
				"link-profile": linkProfile,
			}),
			WithCachedPbFiles(map[string]*pb.SnapshotWithType{
				filepath.Join("objects", "basic-target.pb"): {
					SbType: model.SmartBlockType_Page,
					Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
						Details: &types.Struct{Fields: map[string]*types.Value{
							bundle.RelationKeyId.String():     pbtypes.String("basic-target"),
							bundle.RelationKeyLayout.String(): pbtypes.Float64(float64(model.ObjectType_basic)),
						}},
					}},
				},
				filepath.Join("objects", "set-target.pb"): {
					SbType: model.SmartBlockType_Page,
					Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
						Details: &types.Struct{Fields: map[string]*types.Value{
							bundle.RelationKeyId.String():     pbtypes.String("set-target"),
							bundle.RelationKeyLayout.String(): pbtypes.Float64(float64(model.ObjectType_set)),
						}},
					}},
				},
				filepath.Join("objects", "todo-target.pb"): {
					SbType: model.SmartBlockType_Page,
					Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
						Details: &types.Struct{Fields: map[string]*types.Value{
							bundle.RelationKeyId.String():     pbtypes.String("todo-target"),
							bundle.RelationKeyLayout.String(): pbtypes.Float64(float64(model.ObjectType_todo)),
						}},
					}},
				},
				filepath.Join("objects", "profile-target.pb"): {
					SbType: model.SmartBlockType_Page,
					Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
						Details: &types.Struct{Fields: map[string]*types.Value{
							bundle.RelationKeyId.String():     pbtypes.String("profile-target"),
							bundle.RelationKeyLayout.String(): pbtypes.Float64(float64(model.ObjectType_profile)),
						}},
					}},
				},
			}),
		)
		r.Root = rootBlock

		// when
		result := r.GetLinkObjectIds()

		// then
		// Only basic layout should be included
		assert.Equal(t, []string{"basic-target"}, result)
	})
}
