package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveIframeWidthHeight(t *testing.T) {
	text := `<iframe width="560" height="315" src="https://www.youtube.com/embed/vPOhoud5NIM?si=4Su6t6qt7Q_Ssd4I" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>`
	expected := `<iframe   src="https://www.youtube.com/embed/vPOhoud5NIM?si=4Su6t6qt7Q_Ssd4I" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>`
	actual := removeIframeWidthHeight(text)

	assert.Equal(t, expected, actual)
}
