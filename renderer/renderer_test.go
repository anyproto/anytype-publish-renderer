package renderer

import (
	"testing"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/assert"
)

type TestRenderer struct {
	*Renderer
}

type Option func(*TestRenderer)

func NewTestRenderer(opts ...Option) *TestRenderer {
	renderer := &TestRenderer{Renderer: &Renderer{
		Sp: &pb.SnapshotWithType{Snapshot: &pb.ChangeSnapshot{
			Data: &model.SmartBlockSnapshotBase{
				Details: &types.Struct{Fields: make(map[string]*types.Value)},
			},
		}},
		UberSp:        &PublishingUberSnapshot{PbFiles: make(map[string]string)},
		CachedPbFiles: make(map[string]*pb.SnapshotWithType),
		BlocksById:    make(map[string]*model.Block),
	}}
	for _, opt := range opts {
		opt(renderer)
	}
	return renderer
}

func (t *TestRenderer) ensureUberSpInitialized() {
	if t.UberSp == nil {
		t.UberSp = &PublishingUberSnapshot{PbFiles: make(map[string]string)}
	}
	if t.UberSp.PbFiles == nil {
		t.UberSp.PbFiles = make(map[string]string)
	}
}

func WithPbFiles(pbFiles map[string]string) Option {
	return func(t *TestRenderer) {
		t.ensureUberSpInitialized()
		t.UberSp.PbFiles = pbFiles
	}
}

func WithCachedPbFiles(cachedPbFiles map[string]*pb.SnapshotWithType) Option {
	return func(t *TestRenderer) {
		t.CachedPbFiles = cachedPbFiles
	}
}

func WithBlocksById(blocksById map[string]*model.Block) Option {
	return func(t *TestRenderer) {
		t.BlocksById = blocksById
	}
}

func WithRootSnapshot(rootSnapshot *pb.SnapshotWithType) Option {
	return func(t *TestRenderer) {
		t.Sp = rootSnapshot
	}
}

func WithConfig(config RenderConfig) Option {
	return func(t *TestRenderer) {
		t.Config = config
	}
}

func WithLinkedSnapshot(test *testing.T, fileName string, sn *pb.SnapshotWithType) Option {
	return func(t *TestRenderer) {
		t.ensureUberSpInitialized()
		marshaler := jsonpb.Marshaler{}
		json, err := marshaler.MarshalToString(sn)
		if assert.NoError(test, err) {
			t.UberSp.PbFiles[fileName] = json
		}
	}
}

func WithObjectTypeDetails(details *types.Struct) Option {
	return func(t *TestRenderer) {
		t.ObjectTypeDetails = details
	}
}
