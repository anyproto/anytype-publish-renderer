package renderer

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/assert"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"github.com/anyproto/anytype-publish-renderer/utils/tests/htmltag"
)

func TestMakeBookmarkRendererParams(t *testing.T) {
	tests := []struct {
		name           string
		block          *model.Block
		pbFiles        map[string]*pb.SnapshotWithType
		expected       *BlockParams
		wantErr        bool
		pathAssertions []struct {
			path          string
			expectedValue string
		}
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
							bundle.RelationKeySource.String():      pbtypes.String("https://example.com"),
						}},
					}},
				},
			},
			expected: &BlockParams{
				Id:      "block1",
				Classes: []string{"block", "align0", "blockBookmark"},
			},
			pathAssertions: []struct {
				path          string
				expectedValue string
			}{
				{"a > attrs[href]", "https://example.com"},
				{"a > div > attrs[class]", "side left"},
				{"a > div > div.link > Content", "example.com"},
				{"a > div > div.name > Content", "name1"},
				{"a > div > div.descr > Content", "description1"},
			},
		},
		{
			name: "missing details",
			block: &model.Block{
				Id: "block2",
				Content: &model.BlockContentOfBookmark{
					Bookmark: &model.BlockContentBookmark{
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
							bundle.RelationKeySource.String():      pbtypes.String("::::"),
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
					Bookmark: &model.BlockContentBookmark{TargetObjectId: "object3"},
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
							bundle.RelationKeySource.String():      pbtypes.String(""),
						}},
					}},
				},
			},
			expected: nil,
		},
		{
			name: "empty URL in source, try to get from block",
			block: &model.Block{
				Id: "block3",
				Content: &model.BlockContentOfBookmark{
					Bookmark: &model.BlockContentBookmark{
						TargetObjectId: "object3",
						Url:            "https://example.com"},
				},
			},
			pbFiles: map[string]*pb.SnapshotWithType{
				filepath.Join("objects", "object3.pb"): {
					SbType: model.SmartBlockType_Page,
					Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
						Details: &types.Struct{Fields: map[string]*types.Value{
							bundle.RelationKeySource.String(): pbtypes.String(""),
						}},
					}},
				},
			},
			expected: &BlockParams{
				Id:      "block3",
				Classes: []string{"block", "align0", "blockBookmark"},
			},
			pathAssertions: []struct {
				path          string
				expectedValue string
			}{
				{"a > attrs[href]", "https://example.com"},
				{"a > div > div.link > Content", "example.com"},
			},
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

				got, err := htmltag.HtmlToTag(builder.String())
				if (err != nil) != tt.wantErr {
					t.Errorf("HtmlToTag() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				for _, assertion := range tt.pathAssertions {
					htmltag.AssertPath(t, got, assertion.path, assertion.expectedValue)
				}

			}
		})
	}
}
