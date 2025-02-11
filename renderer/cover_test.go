package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeRenderCoverParams(t *testing.T) {
	t.Run("cover params", func(t *testing.T) {
		r := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
		expected := &CoverRenderParams{
			Id:      "bafyreic35rt6o6jbpfoibui4oskwo2dqurapsyuub4k7o42uatcteucan4",
			Classes: "type1 bafyreic35rt6o6jbpfoibui4oskwo2dqurapsyuub4k7o42uatcteucan4",
			Src:     "../test_snapshots/Anytype.WebPublish.20241217.112212.67/files/640px-anatomy_of_a_sunset-2.webp",
		}

		actual, err := r.makeRenderPageCoverParams()
		if assert.NoError(t, err) {
			assert.Equal(t, expected.Id, actual.Id)
			assert.Equal(t, expected.Classes, actual.Classes)
			assert.Equal(t, expected.Src, actual.Src)
		}
	})

	t.Run("solid color cover", func(t *testing.T) {
		r := getTestRenderer("test-solid-color-cover")
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
		r := getTestRenderer("test-gradient-cover")
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
