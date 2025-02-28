package integration

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"

	"github.com/anyproto/anytype-publish-renderer/renderer"
)

func TestRenderer(t *testing.T) {
	testDir := filepath.Join("testdata", "Anytype.WebPublish.20241217.112212.67")
	testRenderer, err := makeTestRenderer(testDir)
	assert.NoError(t, err)

	file, err := os.Create("index.html")
	assert.NoError(t, err)
	defer file.Close()
	err = testRenderer.Render(file)
	assert.NoError(t, err)
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
