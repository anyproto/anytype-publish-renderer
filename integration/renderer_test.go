package integration

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anyproto/anytype-publish-renderer/renderer"
)

func TestRenderer(t *testing.T) {
	// given
	testDir := "testdata"
	testRenderer, err := makeTestRenderer(testDir)
	assert.NoError(t, err)
	buffer := bytes.NewBuffer(nil)

	// then
	err = testRenderer.Render(buffer)

	// when
	assert.NoError(t, err)
	fileContent, err := os.ReadFile("index.html")
	assert.NoError(t, err)
	assert.Equal(t, fileContent, buffer.Bytes())
}

func makeTestRenderer(dir string) (*renderer.Renderer, error) {
	config := renderer.RenderConfig{
		StaticFilesPath:  "/static",
		PublishFilesPath: dir,
		PrismJsCdnUrl:    "https://cdn.jsdelivr.net/npm/prismjs@1.29.0",
		AnytypeCdnUrl:    "https://anytype-static.fra1.cdn.digitaloceanspaces.com",
		AnalyticsCode:    `<script>console.log("sending dummy analytics...")</script>`,
	}

	r, err := renderer.NewRenderer(config)

	if err != nil {
		return nil, err
	}

	return r, nil
}
