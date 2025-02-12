package renderer

import (
	"testing"

	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/assert"
)

func TestRenderer_MakeRootRenderParams(t *testing.T) {
	t.Run("non empty width", func(t *testing.T) {
		// given
		r := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
		id := "blockId"
		width := 100
		expected := &RootRenderParams{
			Style: `
		<style> 
			.blocks {
				width: max(60%, min(calc(100% - 96px), calc(60% + (100% - 60% - 96px) * 100.000000)));
			}
		</style> 
	`,
		}

		// when
		actual := r.makeRootRenderParams(&model.Block{
			Id: id,
			Fields: &types.Struct{Fields: map[string]*types.Value{
				"width": pbtypes.Float64(float64(width)),
			}},
			Content: &model.BlockContentOfSmartblock{Smartblock: &model.BlockContentSmartblock{}},
		})

		// then
		assert.Equal(t, expected.Style, actual.Style)
	})
	t.Run("empty width", func(t *testing.T) {
		r := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
		id := "blockId"
		expected := &RootRenderParams{
			Style: `
		<style> 
			.blocks {
				width: max(60%, min(calc(100% - 96px), calc(60% + (100% - 60% - 96px) * 0.000000)));
			}
		</style> 
	`,
		}

		// when
		actual := r.makeRootRenderParams(&model.Block{
			Id:      id,
			Content: &model.BlockContentOfSmartblock{Smartblock: &model.BlockContentSmartblock{}},
		})

		// then
		assert.Equal(t, expected.Style, actual.Style)
	})
}
