package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeRenderFileParams(t *testing.T) {
	t.Run("image file", func(t *testing.T) {
		r := getTestRenderer("snapshot_pb")
		id := "66c7055b7e4bcd7bc81f3f37"
		imageBlock := r.BlocksById[id]

		expected := &FileImageRenderParams{
			Id:      id,
			Classes: "align1",
			Src:     "/../test_snapshots/snapshot_pb/files/img_5296.jpeg",
		}

		actual, err := r.MakeRenderFileImageParams(imageBlock)
		if assert.NoError(t, err) {
			assert.Equal(t, expected.Id, actual.Id)
			assert.Equal(t, expected.Classes, actual.Classes)
			assert.Equal(t, expected.Src, actual.Src)
		}

	})
}
