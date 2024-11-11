package renderer

import (
	"context"
	"fmt"
	"io"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/logging"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-publish-renderer/renderer/blockutils"
	"go.uber.org/zap"

	"github.com/a-h/templ"
	"github.com/gogo/protobuf/proto"
	types "github.com/gogo/protobuf/types"
)

var log = logging.Logger("renderer").Desugar()

// Maps asset addresses from snapshot to different location
// <a href>, <img src>, emoji address, etc
type AssetResolver interface {
	GetSnapshotPbFile(string) ([]byte, error)
	GetRootPagePath() string
	ByEmojiCode(rune) string
	ByTargetObjectId(string) (string, error)
}

type Renderer struct {
	Sp       *pb.SnapshotWithType
	Out      io.Writer
	RootComp templ.Component

	Root       *model.Block
	BlocksById map[string]*model.Block

	AssetResolver AssetResolver
}

func NewRenderer(resolver AssetResolver, writer io.Writer) (r *Renderer, err error) {
	rootPath := resolver.GetRootPagePath()
	snapshotData, err := resolver.GetSnapshotPbFile(rootPath)
	if err != nil {
		fmt.Printf("Error reading protobuf snapshot: %v\n", err)
		return
	}

	snapshot := pb.SnapshotWithType{}
	err = proto.Unmarshal(snapshotData, &snapshot)
	if err != nil {
		return
	}

	if snapshot.SbType != model.SmartBlockType_Page {
		err = fmt.Errorf("published snaphost is not Page, %d", snapshot.SbType)
		return
	}

	blocks := snapshot.Snapshot.Data.GetBlocks()
	blocksById := make(map[string]*model.Block)
	for _, block := range blocks {

		blocksById[block.Id] = block
	}

	details := snapshot.Snapshot.Data.GetDetails()

	r = &Renderer{
		Sp:            &snapshot,
		Out:           writer,
		BlocksById:    blocksById,
		Root:          blocks[0],
		AssetResolver: resolver,
	}

	specialBlocks := []string{"title", "description"}
	r.hydrateSpecialBlocks(specialBlocks, details)

	return
}

func (r *Renderer) Render() (err error) {
	err = r.RootComp.Render(context.Background(), r.Out)
	if err != nil {
		return
	}
	return
}

// Adds text from Details to special blocks like `title`
func (r *Renderer) hydrateSpecialBlocks(blockIds []string, details *types.Struct) {
	for _, bId := range blockIds {
		titleBlock, ok := r.BlocksById[bId]
		if !ok {
			log.Warn("hydrate: block not found, skipping", zap.String("id", bId))
			return
		}

		err := blockutils.HydrateBlock(titleBlock, details)
		if err != nil {
			log.Warn("hydrate: failed to hydrate block",
				zap.String("id", bId),
				zap.Error(err))
		}

	}

}
