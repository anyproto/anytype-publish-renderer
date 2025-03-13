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

func TestMakeRenderCoverParams(t *testing.T) {
	t.Run("cover params", func(t *testing.T) {
		coverId := "coverId"
		r := NewTestRenderer(
			WithRootSnapshot(&pb.SnapshotWithType{
				Snapshot: &pb.ChangeSnapshot{
					Data: &model.SmartBlockSnapshotBase{
						Details: &types.Struct{
							Fields: map[string]*types.Value{
								bundle.RelationKeyCoverType.String(): pbtypes.Int64(1),
								bundle.RelationKeyCoverId.String():   pbtypes.String(coverId),
							},
						},
					},
				},
			}),
			WithLinkedSnapshot(t, filepath.Join("filesObjects", coverId+pbExt), &pb.SnapshotWithType{
				SbType: model.SmartBlockType_FileObject,
				Snapshot: &pb.ChangeSnapshot{
					Data: &model.SmartBlockSnapshotBase{Details: &types.Struct{
						Fields: map[string]*types.Value{
							bundle.RelationKeySource.String(): pbtypes.String("test.jpg"),
						},
					}},
				}},
			),
		)
		expected := &CoverRenderParams{
			Id:      coverId,
			Classes: "type1 " + coverId,
			Src:     "/test.jpg",
		}

		actual, err := r.makeRenderPageCoverParams()
		if assert.NoError(t, err) {
			assert.Equal(t, expected.Id, actual.Id)
			assert.Equal(t, expected.Classes, actual.Classes)
			assert.Equal(t, expected.Src, actual.Src)
		}
	})

	t.Run("solid color cover", func(t *testing.T) {
		r := &Renderer{Sp: &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{
				Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{
						Fields: map[string]*types.Value{
							bundle.RelationKeyCoverType.String(): pbtypes.Int64(2),
							bundle.RelationKeyCoverId.String():   pbtypes.String("red"),
						},
					},
				},
			},
		}}
		expected := &CoverRenderParams{
			Id:        "red",
			Classes:   "type2 red",
			CoverType: CoverType_Color,
		}

		actual, err := r.makeRenderPageCoverParams()
		if assert.NoError(t, err) {
			assert.Equal(t, expected.Id, actual.Id)
			assert.Equal(t, expected.Classes, actual.Classes)
			assert.Equal(t, expected.CoverType, actual.CoverType)
		}
	})

	t.Run("gradient cover", func(t *testing.T) {
		r := NewTestRenderer(
			WithRootSnapshot(&pb.SnapshotWithType{
				Snapshot: &pb.ChangeSnapshot{
					Data: &model.SmartBlockSnapshotBase{
						Details: &types.Struct{
							Fields: map[string]*types.Value{
								bundle.RelationKeyCoverType.String(): pbtypes.Int64(3),
								bundle.RelationKeyCoverId.String():   pbtypes.String("blue"),
							},
						},
					},
				},
			}),
		)
		expected := &CoverRenderParams{
			Id:        "blue",
			Classes:   "type3 blue",
			CoverType: CoverType_Gradient,
		}

		actual, err := r.makeRenderPageCoverParams()
		if assert.NoError(t, err) {
			assert.Equal(t, expected.Id, actual.Id)
			assert.Equal(t, expected.Classes, actual.Classes)
			assert.Equal(t, expected.CoverType, actual.CoverType)
		}
	})
}
