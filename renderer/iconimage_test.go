package renderer

import (
	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"github.com/gogo/protobuf/types"
	"path/filepath"
	"testing"

	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"

	"github.com/stretchr/testify/assert"
)

func TestMakeRenderPageIconImageParams(t *testing.T) {
	t.Run("icon image emoji", func(t *testing.T) {
		r := NewTestRenderer(
			WithRootSnapshot(&pb.SnapshotWithType{
				Snapshot: &pb.ChangeSnapshot{
					Data: &model.SmartBlockSnapshotBase{
						Details: &types.Struct{
							Fields: map[string]*types.Value{
								bundle.RelationKeyIconEmoji.String(): pbtypes.String("😃"),
							},
						},
					},
				},
			}),
		)
		expected := &IconImageRenderParams{
			Src: "/emojies/1f603.png",
		}

		actual := r.MakeRenderIconObjectParams(r.Sp.GetSnapshot().GetData().GetDetails(), &IconObjectProps{
			NoDefault: true,
			Size:      pageIconInitSize(model.ObjectType_basic),
		})
		assert.Equal(t, expected.Src, actual.Src)
	})

	t.Run("icon image uploaded", func(t *testing.T) {
		r := NewTestRenderer(
			WithRootSnapshot(&pb.SnapshotWithType{
				Snapshot: &pb.ChangeSnapshot{
					Data: &model.SmartBlockSnapshotBase{
						Details: &types.Struct{
							Fields: map[string]*types.Value{
								bundle.RelationKeyIconImage.String(): pbtypes.String("iconImage"),
							},
						},
					},
				},
			}),
			WithCachedPbFiles(map[string]*pb.SnapshotWithType{
				filepath.Join("filesObjects", "iconImage"+pbExt): {
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
		expected := &IconImageRenderParams{
			Src: "/test.jpg",
		}

		actual := r.MakeRenderIconObjectParams(r.Sp.GetSnapshot().GetData().GetDetails(), &IconObjectProps{
			NoDefault: true,
			Size:      pageIconInitSize(model.ObjectType_basic),
		})
		assert.Equal(t, expected.Src, actual.Src)
	})

}
