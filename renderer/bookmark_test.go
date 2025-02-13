package renderer

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"github.com/gogo/protobuf/types"

	"github.com/stretchr/testify/assert"

	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

func TestMakeBookmarkRendererParams(t *testing.T) {
	tests := []struct {
		name         string
		block        *model.Block
		pbFiles      map[string]*pb.SnapshotWithType
		expected     *BlockParams
		expectedHtml string
	}{
		{
			name: "valid bookmark",
			block: &model.Block{
				Id: "block1",
				Content: &model.BlockContentOfBookmark{
					Bookmark: &model.BlockContentBookmark{
						Url:            "https://example.com",
						TargetObjectId: "object1",
					},
				},
			},
			pbFiles: map[string]*pb.SnapshotWithType{
				filepath.Join("objects", "object1.pb"): {
					SbType: model.SmartBlockType_Page,
					Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
						Details: &types.Struct{Fields: map[string]*types.Value{
							bundle.RelationKeyIconImage.String():   pbtypes.String("favicon1"),
							bundle.RelationKeyPicture.String():     pbtypes.String("image1"),
							bundle.RelationKeyDescription.String(): pbtypes.String("description1"),
							bundle.RelationKeyName.String():        pbtypes.String("name1"),
						}},
					}},
				},
			},
			expected: &BlockParams{
				Id:      "block1",
				Classes: []string{"block", "align0", "blockBookmark"},
			},
			expectedHtml: `<a href="https://example.com" target="_blank" class="inner"><div class="side left"><div class="link">example.com</div><div class="name">name1</div><div class="descr">description1</div></div><div class="side right"></div></a>`,
		},
		{
			name: "missing details",
			block: &model.Block{
				Id: "block2",
				Content: &model.BlockContentOfBookmark{
					Bookmark: &model.BlockContentBookmark{
						Url:            "https://example.com",
						TargetObjectId: "object12",
					},
				},
			},
			pbFiles: map[string]*pb.SnapshotWithType{
				filepath.Join("objects", "object12.pb"): {
					SbType:   model.SmartBlockType_Page,
					Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{}},
				},
			},
			expected: nil,
		},
		{
			name: "missing bookmark",
			block: &model.Block{
				Id: "block2",
				Content: &model.BlockContentOfBookmark{
					Bookmark: &model.BlockContentBookmark{
						Url:            "https://example.com",
						TargetObjectId: "object12",
					},
				},
			},
			expected: nil,
		},
		{
			name: "invalid URL",
			block: &model.Block{
				Id: "block3",
				Content: &model.BlockContentOfBookmark{
					Bookmark: &model.BlockContentBookmark{
						Url:            "::::",
						TargetObjectId: "object3",
					},
				},
			},
			pbFiles: map[string]*pb.SnapshotWithType{
				filepath.Join("objects", "object3.pb"): {
					SbType: model.SmartBlockType_Page,
					Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
						Details: &types.Struct{Fields: map[string]*types.Value{
							bundle.RelationKeyIconImage.String():   pbtypes.String("favicon3"),
							bundle.RelationKeyPicture.String():     pbtypes.String("image3"),
							bundle.RelationKeyDescription.String(): pbtypes.String("description3"),
							bundle.RelationKeyName.String():        pbtypes.String("name3"),
						}},
					}},
				},
			},
			expected: nil,
		},
		{
			name: "empty URL",
			block: &model.Block{
				Id: "block3",
				Content: &model.BlockContentOfBookmark{
					Bookmark: &model.BlockContentBookmark{},
				},
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Renderer{CachedPbFiles: make(map[string]*pb.SnapshotWithType), UberSp: &PublishingUberSnapshot{PbFiles: make(map[string]string)}}
			r.CachedPbFiles = tt.pbFiles
			result := r.makeBookmarkBlockParams(tt.block)
			if tt.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.Equal(t, tt.expected.Classes, result.Classes)
				assert.Equal(t, tt.expected.Id, result.Id)
				builder := strings.Builder{}
				err := result.Content.Render(context.Background(), &builder)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedHtml, builder.String())
			}
		})
	}
}
