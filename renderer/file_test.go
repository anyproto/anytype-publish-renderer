package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeRenderFileParams(t *testing.T) {
	t.Run("image file", func(t *testing.T) {
		r := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
		id := "66c7055b7e4bcd7bc81f3f37"
		imageBlock := r.BlocksById[id]

		expected := &FileMediaRenderParams{
			Id:      id,
			Classes: []string{"align0"},
			Src:     "../test_snapshots/Anytype.WebPublish.20241217.112212.67/files/img_5296.jpeg",
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
