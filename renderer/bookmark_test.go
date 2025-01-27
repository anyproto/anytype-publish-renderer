package renderer

import (
	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"github.com/gogo/protobuf/types"
	"path/filepath"
	"testing"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/stretchr/testify/assert"
)

func TestMakeBookmarkRendererParams(t *testing.T) {
	tests := []struct {
		name     string
		block    *model.Block
		pbFiles  map[string]*pb.SnapshotWithType
		expected *BookmarkRendererParams
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
			expected: &BookmarkRendererParams{
				Id:          "block1",
				Url:         "example.com",
				Name:        "name1",
				Description: "description1",
				SafeUrl:     templ.SafeURL("https://example.com"),
			},
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
			expected: &BookmarkRendererParams{
				IsEmpty: true,
			},
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
			expected: &BookmarkRendererParams{
				IsEmpty: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
			r.CachedPbFiles = tt.pbFiles
			result := r.MakeBookmarkRendererParams(tt.block)
			assert.Equal(t, tt.expected, result)
		})
	}
}
