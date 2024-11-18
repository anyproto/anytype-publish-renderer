package renderer

import (
	"os"

	"github.com/anyproto/anytype-publish-renderer/cmd/resolver"
	"go.uber.org/zap"
)

var testRenderRoots map[string]string
var testRenderers map[string]*Renderer

func makeTestRenderer(id, rootId string) *Renderer {
	resolver := resolver.SimpleAssetResolver{
		CdnUrl:      "http://test-cdn",
		SnapshotDir: "../test_snapshots/" + id,
		RootPageId:  rootId,
	}

	r, err := NewRenderer(resolver, os.Stdout)
	if err != nil {
		log.Fatal("failed to make test renderer", zap.Error(err))
	}

	return r
}
func makeTestRenderRoots() {
	if testRenderRoots == nil {
		testRenderRoots = map[string]string{
			"snapshot_pb":  "bafyreiftcdiken5kayp3x4ix7tm4okmizhyshev3jjl5r2jjenz2d5uwui",
			"snapshot_pb2": "bafyreiecs2ivic6ne2lvkohrf3ojqizngmsk5ywilexo6xhtzlqfgngp64",
			"snapshot_pb3": "bafyreiecs2ivic6ne2lvkohrf3ojqizngmsk5ywilexo6xhtzlqfgngp64",
		}
	}
}
func getTestRenderer(id string) *Renderer {
	makeTestRenderRoots()
	if testRenderers == nil {
		testRenderers = make(map[string]*Renderer, 0)
	}

	if _, ok := testRenderers[id]; !ok {
		rootId := testRenderRoots[id]
		testRenderers[id] = makeTestRenderer(id, rootId)
	}
	return testRenderers[id]
}
