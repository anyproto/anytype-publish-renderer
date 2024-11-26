package renderer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

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

	BlockNumbers map[string]int

	AssetResolver AssetResolver
}

func NewRenderer(resolver AssetResolver, writer io.Writer) (r *Renderer, err error) {
	rootPath := resolver.GetRootPagePath()
	snapshotData, err := resolver.GetSnapshotPbFile(rootPath)
	if err != nil {
		log.Error("Error reading protobuf snapshot", zap.Error(err))
		return
	}

	snapshot := pb.SnapshotWithType{}
	err = proto.Unmarshal(snapshotData, &snapshot)
	if err != nil {
		return
	}

	var snapshotJson []byte
	snapshotJson, err = json.Marshal(snapshot)
	os.WriteFile("./snapshot.json", snapshotJson, 0644)

	if snapshot.SbType != model.SmartBlockType_Page {
		log.Error("published snaphost is not Page", zap.Int("type", int(snapshot.SbType)))
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
		BlockNumbers:  make(map[string]int),
		Root:          blocks[0],
		AssetResolver: resolver,
	}

	specialBlocks := []string{"title", "description"}
	r.hydrateSpecialBlocks(specialBlocks, details)
	r.hydrateNumberBlocks()
	fmt.Printf("--num:: \n %#v\n", r.BlockNumbers)

	return
}

func (r *Renderer) Render() (err error) {
	err = r.RootComp.Render(context.Background(), r.Out)
	if err != nil {
		return
	}
	return
}

// Adds numbers as marks to text blocks with type "number"
func (r *Renderer) hydrateNumberBlocks() {
	for _, block := range r.Sp.Snapshot.Data.GetBlocks() {
		prevNumber := 0
		for _, chId := range block.ChildrenIds {
			chBlock := r.BlocksById[chId]
			if t := chBlock.GetText(); t != nil {
				if t.GetStyle() == model.BlockContentText_Marked {
					r.BlockNumbers[chId] = prevNumber + 1
				} else {
					prevNumber = 0
				}
			}
		}
	}
}

// Adds text from Details to special blocks like `title`
func (r *Renderer) hydrateSpecialBlocks(blockIds []string, details *types.Struct) {
	for _, bId := range blockIds {
		block, ok := r.BlocksById[bId]
		if !ok {
			log.Warn("hydrate: block not found, skipping", zap.String("id", bId))
			return
		}

		err := blockutils.HydrateBlock(block, details)
		if err != nil {
			log.Warn("hydrate: failed to hydrate block",
				zap.String("id", bId),
				zap.Error(err))
		}

	}

}
