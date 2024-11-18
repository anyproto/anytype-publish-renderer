package renderer

import (
	"os"
	"testing"

	"github.com/anyproto/anytype-publish-renderer/cmd/resolver"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var testRenderer *Renderer

func makeTestRenderer() *Renderer {
	resolver := resolver.SimpleAssetResolver{
		CdnUrl:      "http://test-cdn",
		SnapshotDir: "../test_snapshots/snapshot_pb",
		RootPageId:  "bafyreiftcdiken5kayp3x4ix7tm4okmizhyshev3jjl5r2jjenz2d5uwui",
	}

	r, err := NewRenderer(resolver, os.Stdout)
	if err != nil {
		log.Fatal("failed to make test renderer", zap.Error(err))
	}

	return r
}
func getTestRenderer() *Renderer {
	if testRenderer == nil {
		testRenderer = makeTestRenderer()
	}
	return testRenderer
}

func TestMakeRenderDivParams(t *testing.T) {
	r := getTestRenderer()
	divBlock := r.BlocksById["66c5b61a7e4bcd764b24c213"]

	expected := &DivRenderParams{
		Id:      "66c5b61a7e4bcd764b24c213",
		Classes: "divDot",
	}

	actual := r.MakeRenderDivParams(divBlock)

	assert.Equal(t, expected.Id, actual.Id)
	assert.Equal(t, expected.Classes, actual.Classes)
}
