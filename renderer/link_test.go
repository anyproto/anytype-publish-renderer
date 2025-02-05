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
	tests := []struct {
		name     string
		block    *model.Block
		pbFiles  map[string]*pb.SnapshotWithType
		expected *LinkRenderParams
	}{
		{
			name: "Target details not found",
			block: &model.Block{
				Content: &model.BlockContentOfLink{
					Link: &model.BlockContentLink{
						TargetBlockId: "nonexistent-id",
					},
				},
			},
			expected: &LinkRenderParams{IsDeleted: true},
		},
		{
			name: "Deleted block",
			block: &model.Block{
				Content: &model.BlockContentOfLink{
					Link: &model.BlockContentLink{
						TargetBlockId: "deleted-id",
					},
				},
			},
			pbFiles: map[string]*pb.SnapshotWithType{
				filepath.Join("objects", "deleted-id.pb"): {
					SbType: model.SmartBlockType_Page,
					Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
						Details: &types.Struct{Fields: map[string]*types.Value{
							bundle.RelationKeyId.String():        pbtypes.String("deleted-id.pb"),
							bundle.RelationKeyIsDeleted.String(): pbtypes.Bool(true),
						}},
					}},
				},
			},
			expected: &LinkRenderParams{IsDeleted: true},
		},
		{
			name: "Archived block",
			block: &model.Block{
				Content: &model.BlockContentOfLink{
					Link: &model.BlockContentLink{
						TargetBlockId: "archived-id",
					},
				},
			},
			pbFiles: map[string]*pb.SnapshotWithType{
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
			},
			expected: &LinkRenderParams{
				LayoutClass:   "isPage",
				IsArchived:    "isArchived",
				Name:          "Archived Block",
				IconClass:     "c20",
				Url:           templ.SafeURL("anytype://object?objectId=archived-id&spaceId=spaceId"),
			},
		},
		{
			name: "Block with icon emoji",
			block: &model.Block{
				Content: &model.BlockContentOfLink{
					Link: &model.BlockContentLink{
						IconSize:      model.BlockContentLink_SizeMedium,
						TargetBlockId: "emoji-icon-id",
					},
				},
			},
			pbFiles: map[string]*pb.SnapshotWithType{
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
			},
			expected: &LinkRenderParams{
				LayoutClass:   "isPage",
				Name:          "Emoji Icon Block",
				Icon:          "https://anytype-static.fra1.cdn.digitaloceanspaces.com/emojies/1f60a.png",
				IconClass:     "c20 withIcon",
				IconStyle:     "smileImage c20",
				Url:           templ.SafeURL("anytype://object?objectId=emoji-icon-id&spaceId=spaceId"),
			},
		},
		{
			name: "Block with default icon style",
			block: &model.Block{
				Content: &model.BlockContentOfLink{
					Link: &model.BlockContentLink{
						TargetBlockId: "default-icon-id",
						IconSize:      model.BlockContentLink_SizeMedium,
						CardStyle:     model.BlockContentLink_Card,
					},
				},
			},
			pbFiles: map[string]*pb.SnapshotWithType{
				filepath.Join("objects", "default-icon-id.pb"): {
					SbType: model.SmartBlockType_Page,
					Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
						Details: &types.Struct{Fields: map[string]*types.Value{
							bundle.RelationKeyId.String():      pbtypes.String("default-icon-id"),
							bundle.RelationKeyName.String():    pbtypes.String("Default Icon Block"),
							bundle.RelationKeySpaceId.String(): pbtypes.String("spaceId"),
						}},
					}},
				},
			},
			expected: &LinkRenderParams{
				Name:          "Default Icon Block",
				IconStyle:     "iconCommon icon page c28",
				IconClass:     "c48",
				LayoutClass:   "isPage",
				Url:           templ.SafeURL("anytype://object?objectId=default-icon-id&spaceId=spaceId"),
			},
		},
		{
			name: "Block with collection layout",
			block: &model.Block{
				Content: &model.BlockContentOfLink{
					Link: &model.BlockContentLink{
						TargetBlockId: "collection-id",
						IconSize:      model.BlockContentLink_SizeSmall,
					},
				},
			},
			pbFiles: map[string]*pb.SnapshotWithType{
				filepath.Join("objects", "collection-id.pb"): {
					SbType: model.SmartBlockType_Page,
					Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
						Details: &types.Struct{Fields: map[string]*types.Value{
							bundle.RelationKeyId.String():      pbtypes.String("collection-id"),
							bundle.RelationKeyLayout.String():  pbtypes.Float64(float64(model.ObjectType_collection)),
							bundle.RelationKeyName.String():    pbtypes.String("Collection Block"),
							bundle.RelationKeySpaceId.String(): pbtypes.String("spaceId")},
						}},
					}},
			},
			expected: &LinkRenderParams{
				Name:          "Collection Block",
				LayoutClass:   "isCollection",
				IconStyle:     "iconCommon icon collection c20",
				IconClass:     "c20",
				Url:           templ.SafeURL("anytype://object?objectId=collection-id&spaceId=spaceId"),
			},
		},
		{
			name: "Block with todo layout",
			block: &model.Block{
				Content: &model.BlockContentOfLink{
					Link: &model.BlockContentLink{
						TargetBlockId: "todo-id",
						IconSize:      model.BlockContentLink_SizeMedium,
					},
				},
			},
			pbFiles: map[string]*pb.SnapshotWithType{
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
			},
			expected: &LinkRenderParams{
				Name:          "Todo",
				LayoutClass:   "isTask",
				IconStyle:     "iconCheckbox c20 icon checkbox unset",
				IconClass:     "c20",
				Url:           templ.SafeURL("anytype://object?objectId=todo-id&spaceId=spaceId"),
			},
		},
		{
			name: "Block with todo layout, checkbox set",
			block: &model.Block{
				Content: &model.BlockContentOfLink{
					Link: &model.BlockContentLink{
						TargetBlockId: "todo-id",
						IconSize:      model.BlockContentLink_SizeSmall,
					},
				},
			},
			pbFiles: map[string]*pb.SnapshotWithType{
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
			},
			expected: &LinkRenderParams{
				Name:          "Todo",
				LayoutClass:   "isTask",
				IconStyle:     "iconCheckbox c20 icon checkbox set",
				IconClass:     "c20",
				Url:           templ.SafeURL("anytype://object?objectId=todo-id&spaceId=spaceId"),
			},
		},
		{
			name: "Block with description",
			block: &model.Block{
				Content: &model.BlockContentOfLink{
					Link: &model.BlockContentLink{
						TargetBlockId: "test-id",
						Description:   model.BlockContentLink_Added,
					},
				},
			},
			pbFiles: map[string]*pb.SnapshotWithType{
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
			},
			expected: &LinkRenderParams{
				Name:          "Test",
				LayoutClass:   "isHuman",
				Description:   "description",
				IconClass:     "c20",
				Url:           templ.SafeURL("anytype://object?objectId=test-id&spaceId=spaceId"),
			},
		},
		{
			name: "Block with description from snippet",
			block: &model.Block{
				Content: &model.BlockContentOfLink{
					Link: &model.BlockContentLink{
						TargetBlockId: "test-id",
						Description:   model.BlockContentLink_Content,
						CardStyle:     model.BlockContentLink_Card,
					},
				},
			},
			pbFiles: map[string]*pb.SnapshotWithType{
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
			},
			expected: &LinkRenderParams{
				Name:          "Test",
				LayoutClass:   "isParticipant",
				Description:   "snippet",
				IconClass:     "c20",
				Url:           templ.SafeURL("anytype://object?objectId=test-id&spaceId=spaceId"),
			},
		},
		{
			name: "Block with cover relation",
			block: &model.Block{
				Content: &model.BlockContentOfLink{
					Link: &model.BlockContentLink{
						TargetBlockId: "test-id",
						Relations:     []string{"cover"},
					},
				},
			},
			pbFiles: map[string]*pb.SnapshotWithType{
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
			},
			expected: &LinkRenderParams{
				Name:          "Test",
				LayoutClass:   "isSet",
				IconClass:     "c20",
				CoverClass:    "withCover",
				CoverParams: &CoverRenderParams{
					Id:        "gray",
					Src:       "",
					Classes:   "gray",
					CoverType: 2,
				},
				Url: templ.SafeURL("anytype://object?objectId=test-id&spaceId=spaceId"),
			},
		},
		{
			name: "Block with type relation",
			block: &model.Block{
				Content: &model.BlockContentOfLink{
					Link: &model.BlockContentLink{
						TargetBlockId: "test-id",
						Relations:     []string{"type"},
					},
				},
			},
			pbFiles: map[string]*pb.SnapshotWithType{
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
			},
			expected: &LinkRenderParams{
				Name:          "Test",
				LayoutClass:   "isSet",
				IconClass:     "c20",
				Type:          "Type",
				Url:           templ.SafeURL("anytype://object?objectId=test-id&spaceId=spaceId"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
			r.CachedPbFiles = tt.pbFiles
			result := r.MakeLinkRenderParams(tt.block)
			assert.Equal(t, tt.expected, result)
		})
	}
}
