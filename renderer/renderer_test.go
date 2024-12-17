package renderer

import (
	"go.uber.org/zap"
)

var (
	testRenderers = make(map[string]*Renderer)
)

func makeTestRenderer(dir string) *Renderer {
	config := RenderConfig{
		StaticFilesPath:  "/static",
		PublishFilesPath: "../test_snapshots/" + dir,
		PrismJsCdnUrl:    "https://cdn.jsdelivr.net/npm/prismjs@1.29.0",
		AnytypeCdnUrl:    "https://anytype-static.fra1.cdn.digitaloceanspaces.com",
		AnalyticsCode:    `<script>console.log("sending dummy analytics...")</script>`,
	}

	r, err := NewRenderer(config)

	if err != nil {
		log.Fatal("failed to make test renderer", zap.Error(err))
	}

	return r
}

func getTestRenderer(dir string) *Renderer {
	if _, ok := testRenderers[dir]; !ok {
		testRenderers[dir] = makeTestRenderer(dir)
	}

	return testRenderers[dir]
}
