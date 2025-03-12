package renderer

import (
	"path/filepath"
	"testing"

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
			WithConfig(RenderConfig{StaticFilesPath: filepath.Join("..", "static")}),
		)

		// when
		params := renderer.MakeRenderIconObjectParams(details, &IconObjectProps{})

		// then
		assert.Equal(t, "data:image/svg+xml;charset=utf-8;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHN0eWxlPSJmaWxsOiNmNTU1MjIiIHdpZHRoPSI1MTIiIGhlaWdodD0iNTEyIiB2aWV3Qm94PSIwIDAgNTEyIDUxMiI+PGNpcmNsZSBjeD0iMjU1Ljc1IiBjeT0iNTYiIHI9IjU2Ii8+PHBhdGggZD0iTTM5NC42MywyNzcuOSwzODQuMywyNDMuNDlzMC0uMDcsMC0uMTFsLTIyLjQ2LTc0Ljg2aC0uMDVsLTIuNTEtOC40NWE0NC44Nyw0NC44NywwLDAsMC00My0zMi4wOGgtMTIwYTQ0Ljg0LDQ0Ljg0LDAsMCwwLTQzLDMyLjA4bC0yLjUxLDguNDVoLS4wNmwtMjIuNDYsNzQuODZzMCwuMDcsMCwuMTFMMTE3Ljg4LDI3Ny45Yy0zLjEyLDEwLjM5LDIuMywyMS42NiwxMi41NywyNS4xNGEyMCwyMCwwLDAsMCwyNS42LTEzLjE4bDI1LjU4LTg1LjI1aDBsMi4xNy03LjIzQTgsOCwwLDAsMSwxOTkuMzMsMjAwYTcuNzgsNy43OCwwLDAsMS0uMTcsMS42MXYwTDE1NS40MywzNDcuNEExNiwxNiwwLDAsMCwxNzAuNzUsMzY4aDI5VjQ4Mi42OWMwLDE2LjQ2LDEwLjUzLDI5LjMxLDI0LDI5LjMxczI0LTEyLjg1LDI0LTI5LjMxVjM2OGgxNlY0ODIuNjljMCwxNi40NiwxMC41MywyOS4zMSwyNCwyOS4zMXMyNC0xMi44NSwyNC0yOS4zMVYzNjhoMzBhMTYsMTYsMCwwLDAsMTUuMzMtMjAuNkwzMTMuMzQsMjAxLjU5YTcuNTIsNy41MiwwLDAsMS0uMTYtMS41OSw4LDgsMCwwLDEsMTUuNTQtMi42M2wyLjE3LDcuMjNoMGwyNS41Nyw4NS4yNUEyMCwyMCwwLDAsMCwzODIuMDUsMzAzQzM5Mi4zMiwyOTkuNTYsMzk3Ljc0LDI4OC4yOSwzOTQuNjMsMjc3LjlaIi8+PC9zdmc+", params.Src)
		assert.Equal(t, []string{"iconCommon"}, params.IconClasses)
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
			WithConfig(RenderConfig{StaticFilesPath: filepath.Join("..", "static")}),
		)

		// when
		params := renderer.MakeRenderIconObjectParams(details, &IconObjectProps{})

		// then
		assert.Empty(t, params.Src)
		assert.Empty(t, params.IconClasses)
	})
}
