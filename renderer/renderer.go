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

type Renderer struct {
	Sp       *pb.SnapshotWithType
	Out      io.Writer
	RootComp templ.Component

	Root       *model.Block
	BlocksById map[string]*model.Block
}

func NewRenderer(snapshotData []byte, writer io.Writer) (r *Renderer, err error) {
	snapshot := pb.SnapshotWithType{}
	err = proto.Unmarshal(snapshotData, &snapshot)
	if err != nil {
		return
	}

	var snapshotJson []byte
	snapshotJson, err = json.Marshal(snapshot)
	os.WriteFile("./snapshot.json", snapshotJson, 0644)

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
		Sp:         &snapshot,
		Out:        writer,
		BlocksById: blocksById,
		Root:       blocks[0],
	}
	r.hydrateSpecialBlocks(details)

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
func (r *Renderer) hydrateSpecialBlocks(details *types.Struct) {
	titleBlock, ok := r.BlocksById["title"]
	if !ok {
		log.Warn("hydrate: title block not found, skipping")
		return
	}

	err := blockutils.HydrateBlock(titleBlock, details)
	if err != nil {
		log.Warn("hydrate: failed to hydrate title block", zap.Error(err))
	}
}
