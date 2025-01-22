package renderer

import (
	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"testing"

	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/assert"
)

func TestMakeLinkRenderParams(t *testing.T) {
	tests := []struct {
		name     string
		block    *model.Block
		details  []*pb.DependantDetail
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
			details:  nil,
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
			details: []*pb.DependantDetail{
				{
					Id: "deleted-id",
					Details: &types.Struct{
						Fields: map[string]*types.Value{
							bundle.RelationKeyId.String():        pbtypes.String("deleted-id"),
							bundle.RelationKeyIsDeleted.String(): pbtypes.Bool(true),
						},
					},
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
			details: []*pb.DependantDetail{
				{
					Id: "archived-id",
					Details: &types.Struct{
						Fields: map[string]*types.Value{
							bundle.RelationKeyId.String():         pbtypes.String("archived-id"),
							bundle.RelationKeyIsArchived.String(): pbtypes.Bool(true),
							bundle.RelationKeyName.String():       pbtypes.String("Archived Block"),
						},
					},
				},
			},
			expected: &LinkRenderParams{
				LayoutClass:   "isPage",
				IsArchived:    "isArchived",
				Name:          "Archived Block",
				IconClass:     "c20",
				LinkTypeClass: "text",
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
			details: []*pb.DependantDetail{
				{
					Id: "emoji-icon-id",
					Details: &types.Struct{Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():        pbtypes.String("emoji-icon-id"),
						bundle.RelationKeyName.String():      pbtypes.String("Emoji Icon Block"),
						bundle.RelationKeyIconEmoji.String(): pbtypes.String("😊"),
					},
					},
				}},
			expected: &LinkRenderParams{
				LayoutClass:   "isPage",
				Name:          "Emoji Icon Block",
				Icon:          "https://anytype-static.fra1.cdn.digitaloceanspaces.com/emojies/1f60a.png",
				IconClass:     "c20 withIcon",
				IconStyle:     "smileImage c20",
				LinkTypeClass: "text",
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
			details: []*pb.DependantDetail{
				{
					Id: "default-icon-id",
					Details: &types.Struct{
						Fields: map[string]*types.Value{
							bundle.RelationKeyId.String():   pbtypes.String("default-icon-id"),
							bundle.RelationKeyName.String(): pbtypes.String("Default Icon Block"),
						},
					},
				}},
			expected: &LinkRenderParams{
				Name:          "Default Icon Block",
				IconStyle:     "iconCommon icon page c28",
				IconClass:     "c48",
				LinkTypeClass: "card",
				LayoutClass:   "isPage",
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
			details: []*pb.DependantDetail{
				{
					Id: "collection-id",
					Details: &types.Struct{Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():     pbtypes.String("collection-id"),
						bundle.RelationKeyLayout.String(): pbtypes.Float64(float64(model.ObjectType_collection)),
						bundle.RelationKeyName.String():   pbtypes.String("Collection Block")},
					},
				},
			},
			expected: &LinkRenderParams{
				Name:          "Collection Block",
				LayoutClass:   "isCollection",
				IconStyle:     "iconCommon icon collection c20",
				IconClass:     "c20",
				LinkTypeClass: "text",
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
			details: []*pb.DependantDetail{
				{
					Id: "todo-id",
					Details: &types.Struct{Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():     pbtypes.String("todo-id"),
						bundle.RelationKeyLayout.String(): pbtypes.Float64(float64(model.ObjectType_todo)),
						bundle.RelationKeyName.String():   pbtypes.String("Todo"),
					}},
				},
			},
			expected: &LinkRenderParams{
				Name:          "Todo",
				LayoutClass:   "isTask",
				IconStyle:     "iconCheckbox c20 icon checkbox unset",
				IconClass:     "c20",
				LinkTypeClass: "text",
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
			details: []*pb.DependantDetail{
				{
					Id: "todo-id",
					Details: &types.Struct{Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():     pbtypes.String("todo-id"),
						bundle.RelationKeyLayout.String(): pbtypes.Float64(float64(model.ObjectType_todo)),
						bundle.RelationKeyName.String():   pbtypes.String("Todo"),
						bundle.RelationKeyDone.String():   pbtypes.Bool(true),
					}},
				},
			},
			expected: &LinkRenderParams{
				Name:          "Todo",
				LayoutClass:   "isTask",
				IconStyle:     "iconCheckbox c20 icon checkbox set",
				IconClass:     "c20",
				LinkTypeClass: "text",
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
			details: []*pb.DependantDetail{
				{
					Id: "test-id",
					Details: &types.Struct{Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():          pbtypes.String("test-id"),
						bundle.RelationKeyLayout.String():      pbtypes.Float64(float64(model.ObjectType_profile)),
						bundle.RelationKeyName.String():        pbtypes.String("Test"),
						bundle.RelationKeyDescription.String(): pbtypes.String("description"),
					}},
				},
			},
			expected: &LinkRenderParams{
				Name:          "Test",
				LayoutClass:   "isHuman",
				Description:   "description",
				LinkTypeClass: "text",
				IconClass:     "c20",
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
			details: []*pb.DependantDetail{
				{
					Id: "test-id",
					Details: &types.Struct{Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():      pbtypes.String("test-id"),
						bundle.RelationKeyLayout.String():  pbtypes.Float64(float64(model.ObjectType_participant)),
						bundle.RelationKeyName.String():    pbtypes.String("Test"),
						bundle.RelationKeySnippet.String(): pbtypes.String("snippet"),
					}},
				},
			},
			expected: &LinkRenderParams{
				Name:          "Test",
				LayoutClass:   "isParticipant",
				Description:   "snippet",
				LinkTypeClass: "card",
				IconClass:     "c20",
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
			details: []*pb.DependantDetail{
				{
					Id: "test-id",
					Details: &types.Struct{Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():        pbtypes.String("test-id"),
						bundle.RelationKeyLayout.String():    pbtypes.Float64(float64(model.ObjectType_set)),
						bundle.RelationKeyName.String():      pbtypes.String("Test"),
						bundle.RelationKeyCoverType.String(): pbtypes.Int64(2),
						bundle.RelationKeyCoverId.String():   pbtypes.String("gray"),
					}},
				},
			},
			expected: &LinkRenderParams{
				Name:          "Test",
				LayoutClass:   "isSet",
				IconClass:     "c20",
				LinkTypeClass: "text",
				CoverClass:    "withCover",
				CoverParams: &CoverRenderParams{
					Id:        "gray",
					Src:       "",
					Classes:   "gray",
					CoverType: 2,
				},
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
			details: []*pb.DependantDetail{
				{
					Id: "test-id",
					Details: &types.Struct{Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():     pbtypes.String("test-id"),
						bundle.RelationKeyLayout.String(): pbtypes.Float64(float64(model.ObjectType_set)),
						bundle.RelationKeyName.String():   pbtypes.String("Test"),
						bundle.RelationKeyType.String():   pbtypes.String("type")},
					},
				},
				{
					Id: "type",
					Details: &types.Struct{Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():   pbtypes.String("type"),
						bundle.RelationKeyName.String(): pbtypes.String("Type")},
					},
				},
			},
			expected: &LinkRenderParams{
				Name:          "Test",
				LayoutClass:   "isSet",
				IconClass:     "c20",
				LinkTypeClass: "text",
				Type:          "Type",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
			r.Sp.DependantDetails = tt.details
			result := r.MakeLinkRenderParams(tt.block)
			assert.Equal(t, tt.expected, result)
		})
	}
}
