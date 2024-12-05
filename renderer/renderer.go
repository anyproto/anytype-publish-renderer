package renderer

import (
	"context"
	"io"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/logging"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-publish-renderer/renderer/blockutils"
	"go.uber.org/zap"

	"github.com/a-h/templ"
	"github.com/gogo/protobuf/proto"
)

var log = logging.Logger("renderer").Desugar()

// Maps asset addresses from snapshot to different location
// <a href>, <img src>, emoji address, etc
type AssetResolver interface {
	GetSnapshotPbFile(string) ([]byte, error)
	GetRootPagePath() string
	GetJoinSpaceLink() string
	GetStaticFolderUrl(string) string
	GetAssetUrl(string) string
	GetPrismJsUrl(string) string
	GetEmojiUrl(rune) string
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
		log.Error("Error reading protobuf snapshot index", zap.Error(err))
		return
	}

	snapshot := pb.SnapshotWithType{}
	err = proto.Unmarshal(snapshotData, &snapshot)
	if err != nil {
		return
	}

	if snapshot.SbType != model.SmartBlockType_Page {
		log.Error("published snaphost is not Page", zap.Int("type", int(snapshot.SbType)))
		return
	}

	blocks := snapshot.Snapshot.Data.GetBlocks()
	blocksById := make(map[string]*model.Block)
	for _, block := range blocks {
		blocksById[block.Id] = block
	}

	r = &Renderer{
		Sp:            &snapshot,
		Out:           writer,
		BlocksById:    blocksById,
		BlockNumbers:  make(map[string]int),
		Root:          blocks[0],
		AssetResolver: resolver,
	}

	r.hydrateSpecialBlocks()
	r.hydrateNumberBlocks()

	return
}

func (r *Renderer) Render() (err error) {
	err = r.RootComp.Render(context.Background(), r.Out)
	if err != nil {
		return
	}
	return
}

func (r *Renderer) unwrapLayouts(blocks []*model.Block) []*model.Block {
	ret := make([]*model.Block, 0)
	for _, b := range blocks {
		switch b.Content.(type) {
		case *model.BlockContentOfLayout:
			break
		default:
			ret = append(ret, b)
		}
		chBlocks := make([]*model.Block, len(b.ChildrenIds))
		for i, chId := range b.ChildrenIds {
			chBlocks[i] = r.BlocksById[chId]
		}
		unwrapped := r.unwrapLayouts(chBlocks)
		ret = append(ret, unwrapped...)

	}

	return ret
}

func (r *Renderer) hydrateNumberBlocksInner(blocks []*model.Block) {
	unwrapped := r.unwrapLayouts(blocks)
	prevNumber := 1
	for _, b := range unwrapped {
		if t := b.GetText(); t != nil {
			if t.GetStyle() == model.BlockContentText_Numbered {
				if _, ok := r.BlockNumbers[b.Id]; !ok {
					r.BlockNumbers[b.Id] = prevNumber
					prevNumber++
				}
			} else {
				prevNumber = 1
			}
		}

		if len(b.ChildrenIds) > 0 {
			chBlocks := make([]*model.Block, len(b.ChildrenIds))
			for i, chId := range b.ChildrenIds {
				chBlocks[i] = r.BlocksById[chId]
			}
			r.hydrateNumberBlocksInner(chBlocks)
		}
	}
}

// Adds numbers as marks to text blocks with type "number"
func (r *Renderer) hydrateNumberBlocks() {
	r.hydrateNumberBlocksInner(r.Sp.Snapshot.Data.GetBlocks())
}

// Adds text from Details to special blocks like `title`
func (r *Renderer) hydrateSpecialBlocks() {
	specialBlocks := []string{"title", "description"}
	details := r.Sp.Snapshot.Data.GetDetails()

	for _, bId := range specialBlocks {
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

func Comment(text string) templ.ComponentFunc {

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		io.WriteString(w, "<!--comment:\n")
		io.WriteString(w, text)
		io.WriteString(w, "\n-->\n")
		return nil
	})
}
