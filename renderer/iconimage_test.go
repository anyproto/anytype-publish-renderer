package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeRenderPageIconImageParams(t *testing.T) {
	t.Run("icon image emoji", func(t *testing.T) {
		r := getTestRenderer("test-emoji-icon")
		expected := &IconImageRenderParams{
			Src: "https://anytype-static.fra1.cdn.digitaloceanspaces.com/emojies/1f972.png",
		}

		actual, err := r.MakeRenderPageIconImageParams()
		if assert.NoError(t, err) {
			assert.Equal(t, expected.Src, actual.Src)
		}
	})

	t.Run("icon image uploaded", func(t *testing.T) {
		r := getTestRenderer("test-uploaded-image-icon")
		expected := &IconImageRenderParams{
			Src: "../test_snapshots/test-uploaded-image-icon/files/1737028923-16-01-25_13-02-03.png",
		}

		actual, err := r.MakeRenderPageIconImageParams()
		if assert.NoError(t, err) {
			assert.Equal(t, expected.Src, actual.Src)
		}
	})

}
