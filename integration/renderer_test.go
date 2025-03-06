package integration

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"

	"github.com/anyproto/anytype-publish-renderer/renderer"
)

func TestRenderer(t *testing.T) {
	t.Run("test snapshot", func(t *testing.T) {
		testDir := "testdata"
		testRendering(t, testDir)
	})
	t.Run("test primitives snapshots", func(t *testing.T) {
		testDir := "primitives"
		testRendering(t, testDir)
	})
}

func testRendering(t *testing.T, testDir string) {
	// given
	testRenderer, err := makeTestRenderer(testDir)
	assert.NoError(t, err)
	buffer := bytes.NewBuffer(nil)

	// then
	err = testRenderer.Render(buffer)

	// when
	assert.NoError(t, err)
	fileContent, err := os.ReadFile(filepath.Join(testDir, "index.html"))
	assert.NoError(t, err)
	fileContentStr := strings.TrimSuffix(string(fileContent), "\n")
	if fileContentStr != buffer.String() {
		assert.Fail(t, "")
		diffHTML(string(fileContent), buffer.String())
	}
}

func diffHTML(expected, actual string) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(expected, actual, false)
	dmp.DiffPrettyText(diffs)
	prettyPrintDiff(diffs)
}

func prettyPrintDiff(diffs []diffmatchpatch.Diff) {
	for _, diff := range diffs {
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			fmt.Printf("\033[32m+ %s\033[0m\n", diff.Text)
		case diffmatchpatch.DiffDelete:
			fmt.Printf("\033[31m- %s\033[0m\n", diff.Text)
		case diffmatchpatch.DiffEqual:
			fmt.Printf("  %s\n", diff.Text)
		}
	}
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
