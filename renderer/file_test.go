package renderer

import (
	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"github.com/gogo/protobuf/types"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeRenderFileParams(t *testing.T) {
	t.Run("image file", func(t *testing.T) {
		id := "66c7055b7e4bcd7bc81f3f37"
		targetFileId := "targetFileId"
		r := NewTestRenderer(
			WithBlocksById(map[string]*model.Block{
				id: {
					Id: id,
					Content: &model.BlockContentOfFile{File: &model.BlockContentFile{
						TargetObjectId: targetFileId,
					}},
				},
			}),
			WithCachedPbFiles(map[string]*pb.SnapshotWithType{
				filepath.Join("filesObjects", targetFileId+pbExt): {
					SbType: model.SmartBlockType_FileObject,
					Snapshot: &pb.ChangeSnapshot{
						Data: &model.SmartBlockSnapshotBase{Details: &types.Struct{
							Fields: map[string]*types.Value{
								bundle.RelationKeySource.String(): pbtypes.String("test.jpg"),
							},
						}},
					},
				},
			}),
		)

		imageBlock := r.BlocksById[id]

		expected := &FileMediaRenderParams{
			Id:      id,
			Classes: []string{"align0"},
			Src:     "/test.jpg",
			Width:   "100",
		}

		fileParams, err := r.MakeRenderFileParams(imageBlock)
		actual := fileParams.ToFileMediaRenderParams("100", []string{"align0"})
		if assert.NoError(t, err) {
			assert.Equal(t, expected.Id, actual.Id)
			assert.Equal(t, expected.Classes, actual.Classes)
			assert.Equal(t, expected.Src, actual.Src)
		}
	})
}
