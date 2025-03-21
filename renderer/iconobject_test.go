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

func TestRenderer_MakeRenderIconObjectParams(t *testing.T) {
	t.Run("type icon from relation iconName", func(t *testing.T) {
		// given
		details := &types.Struct{Fields: map[string]*types.Value{
			bundle.RelationKeyIconName.String():       pbtypes.String("woman"),
			bundle.RelationKeyIconOption.String():     pbtypes.Int64(4),
			bundle.RelationKeyResolvedLayout.String(): pbtypes.Int64(int64(model.ObjectType_objectType)),
		}}
		renderer := NewTestRenderer(
			WithConfig(RenderConfig{StaticFilesPath: "static"}),
		)

		// when
		params := renderer.MakeRenderIconObjectParams(details, &IconObjectProps{})

		// then
		assert.Equal(t, "static/img/icon/type/woman.svg", params.SvgSrc)
		assert.Equal(t, []string{"iconCommon"}, params.IconClasses)
		assert.Equal(t, "#f55522", params.SvgColor)
	})
	t.Run("default type icon", func(t *testing.T) {
		// given
		details := &types.Struct{Fields: map[string]*types.Value{
			bundle.RelationKeyResolvedLayout.String(): pbtypes.Int64(int64(model.ObjectType_objectType)),
		}}
		renderer := NewTestRenderer(
			WithConfig(RenderConfig{StaticFilesPath: filepath.Join("..", "static")}),
		)

		// when
		params := renderer.MakeRenderIconObjectParams(details, &IconObjectProps{})

		// then
		assert.Equal(t, "../static/img/icon/default/type.svg", params.Src)
		assert.Equal(t, []string{"iconCommon"}, params.IconClasses)
	})
	t.Run("wrong icon option - default icon", func(t *testing.T) {
		// given
		details := &types.Struct{Fields: map[string]*types.Value{
			bundle.RelationKeyIconName.String():       pbtypes.String("woman"),
			bundle.RelationKeyIconOption.String():     pbtypes.Int64(11),
			bundle.RelationKeyResolvedLayout.String(): pbtypes.Int64(int64(model.ObjectType_objectType)),
		}}
		renderer := NewTestRenderer(
			WithConfig(RenderConfig{StaticFilesPath: "static"}),
		)

		// when
		params := renderer.MakeRenderIconObjectParams(details, &IconObjectProps{})

		// then
		assert.Equal(t, "static/img/icon/type/woman.svg", params.SvgSrc)
		assert.Equal(t, []string{"iconCommon"}, params.IconClasses)
		assert.Equal(t, "", params.SvgColor)
	})
	t.Run("with image in icon", func(t *testing.T) {
		// given
		targetFileId := "fileId"
		details := &types.Struct{Fields: map[string]*types.Value{
			bundle.RelationKeyIconImage.String():      pbtypes.String(targetFileId),
			bundle.RelationKeyResolvedLayout.String(): pbtypes.Int64(int64(model.ObjectType_objectType)),
		}}
		renderer := NewTestRenderer(
			WithConfig(RenderConfig{StaticFilesPath: "static"}),
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

		// when
		params := renderer.MakeRenderIconObjectParams(details, &IconObjectProps{})

		// then
		assert.Equal(t, "/test.jpg", params.Src)
		assert.Equal(t, []string{"iconImage"}, params.IconClasses)
		assert.Equal(t, []string{"iconObject", "withImage"}, params.Classes)
	})
}
