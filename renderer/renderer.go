package renderer

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/logging"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/types"
	"go.uber.org/zap"

	"github.com/anyproto/anytype-publish-renderer/renderer/blockutils"
)

var log = logging.Logger("renderer").Desugar()

type PublishingUberSnapshotMeta struct {
	SpaceId    string `json:"spaceId,omitempty"`
	RootPageId string `json:"rootPageId,omitempty"`
	InviteLink string `json:"inviteLink,omitempty"`
}

// Contains all publishing .pb files
// and publishing meta info
type PublishingUberSnapshot struct {
	Meta PublishingUberSnapshotMeta `json:"meta,omitempty"`

	// A map of "dir/filename.pb -> jsonpb snapshot"
	PbFiles map[string]string `json:"pbFiles,omitempty"`
}

type RenderConfig struct {
	// common for all pages, i.e. layout.css
	StaticFilesPath string
	// assets which belong to published page
	PublishFilesPath string

	PrismJsCdnUrl string
	// anytype cdn, only for emojies for now
	AnytypeCdnUrl string

	// analytics code to inject
	AnalyticsCode string

	// classes for <html> tag, used for debug
	HtmlClasses []string
}

type Renderer struct {
	Sp       *pb.SnapshotWithType
	UberSp   *PublishingUberSnapshot
	RootComp templ.Component
	Config   RenderConfig

	CachedPbFiles map[string]*pb.SnapshotWithType

	Root       *model.Block
	BlocksById map[string]*model.Block

	BlockNumbers      map[string]int
	ObjectTypeDetails *types.Struct
	ResolvedLayout    model.ObjectTypeLayout
	LayoutAlign       int64

	urlsRewriteMap map[string]string
}

func readJsonpbSnapshot(snapshotStr string) (snapshot pb.SnapshotWithType, err error) {
	err = jsonpb.UnmarshalString(snapshotStr, &snapshot)
	if err != nil {
		return
	}
	return
}

func readUberSnapshot(path string) (uberSnapshot PublishingUberSnapshot, err error) {
	var indexFileGz io.Reader

	if strings.HasPrefix(path, "http") {
		var resp *http.Response
		var indexPath string

		indexPath, err = url.JoinPath(path, "index.json.gz")
		if err != nil {
			err = fmt.Errorf("error making http path for index.json.gz: %s", err)
			return
		}

		resp, err = http.Get(indexPath)
		if err != nil {
			err = fmt.Errorf("error reading index.json.gz: %s", err)
			return
		}

		indexFileGz = resp.Body
		defer resp.Body.Close()

	} else {
		var file *os.File
		indexPath := filepath.Join(path, "index.json.gz")
		file, err = os.Open(indexPath)
		if err != nil {
			err = fmt.Errorf("error reading index.json.gz: %s", err)
			return
		}
		indexFileGz = file
		defer file.Close()
	}

	gzReader, err := gzip.NewReader(indexFileGz)
	if err != nil {
		err = fmt.Errorf("error creating .gz reader: %s", err)
		return
	}

	indexBytes, err := io.ReadAll(gzReader)
	if err != nil {
		errgz := gzReader.Close()
		err = fmt.Errorf("error ungzipping index.json.gz: %s", err)
		if errgz != nil {
			err = fmt.Errorf("error closing gzReader: %s", errgz)
		}
		return
	}

	errgz := gzReader.Close()
	if errgz != nil {
		err = fmt.Errorf("error closing gzReader: %s", errgz)
		return
	}

	err = json.Unmarshal(indexBytes, &uberSnapshot)
	if err != nil {
		err = fmt.Errorf("error unmarshaling index.json.gz: %s", err)
		return
	}

	return

}

//lint:ignore U1000 sometimes we want to use this for debugging
func debugJsonSnapshot(snapshot pb.SnapshotWithType) error {
	var snapshotJson []byte
	snapshotJson, err := json.Marshal(snapshot)
	if err != nil {
		log.Error("failed to render snapshot.json", zap.Error(err))
		return err
	}

	err = os.WriteFile("./snapshot.json", snapshotJson, 0644)
	if err != nil {
		log.Error("failed to write snapshot.json", zap.Error(err))
		return err
	}
	return nil
}

func (r *Renderer) maybeAddDebugCss() {
	if isDebugCss := os.Getenv("ANYTYPE_PUBLISH_CSS_DEBUG"); isDebugCss != "" {
		r.Config.HtmlClasses = append(r.Config.HtmlClasses, "debug")
	}
}

func NewRenderer(config RenderConfig) (r *Renderer, err error) {
	defer func() {
		if p := recover(); p != nil {
			stack := string(debug.Stack())
			err = fmt.Errorf("panic: %v, publishFilesPath: %s, stack: %s", p, r.Config.PublishFilesPath, stack)
			log.Error("panic recover", zap.String("where", "NewRenderer()"), zap.Error(err), zap.String("stack", stack))
			return
		}
	}()

	uberSnapshot, err := readUberSnapshot(config.PublishFilesPath)
	if err != nil {
		log.Error("Error reading config.PublishFilesPath ubersnapshot", zap.Error(err))
		return
	}

	rootFilename := fmt.Sprintf("objects/%s.pb", uberSnapshot.Meta.RootPageId)
	snapshot, err := readJsonpbSnapshot(uberSnapshot.PbFiles[rootFilename])
	if err != nil {
		log.Error("Error reading protobuf snapshot index", zap.Error(err))
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
		Sp:             &snapshot,
		UberSp:         &uberSnapshot,
		CachedPbFiles:  make(map[string]*pb.SnapshotWithType),
		BlocksById:     blocksById,
		BlockNumbers:   make(map[string]int),
		Root:           blocks[0],
		Config:         config,
		urlsRewriteMap: make(map[string]string, 0),
	}

	objectType := getRelationField(snapshot.Snapshot.Data.GetDetails(), bundle.RelationKeyType, relationToString)
	if objectType == "" {
		log.Error("no object type in snapshot")
		return
	}

	r.ObjectTypeDetails = r.findTargetDetails(objectType)
	r.ResolvedLayout = r.resolveObjectLayout(snapshot.Snapshot.GetData().GetDetails())

	r.fillLayoutAlign(snapshot.Snapshot.GetData().GetDetails())
	r.maybeAddDebugCss()
	r.hydrateSpecialBlocks()
	r.hydrateNumberBlocks()
	r.RootComp = r.RenderPage()

	return
}

func (r *Renderer) GetLinkObjectIds() (linkObjectIds []string) {
	seen := make(map[string]struct{})

	for _, childID := range r.Root.ChildrenIds {
		b, ok := r.BlocksById[childID]
		if !ok || b == nil {
			continue
		}

		switch b.Content.(type) {
		case *model.BlockContentOfLink:
			targetObjectID := b.GetLink().GetTargetBlockId()
			if targetObjectID == "" {
				continue
			}
			targetDetails := r.findTargetDetails(targetObjectID)
			layout := getRelationField(targetDetails, bundle.RelationKeyLayout, relationToObjectTypeLayout)
			switch layout {
			// TODO: basic = page, what else?
			case model.ObjectType_basic:
				if _, ok := seen[targetObjectID]; !ok {
					seen[targetObjectID] = struct{}{}
				}
			default:
				continue
			}
		default:
			continue
		}

	}

	linkObjectIds = make([]string, len(seen))
	i := 0
	for objID := range seen {
		linkObjectIds[i] = objID
		i++
	}
	return
}

func (r *Renderer) SetUrlRewriteMap(urls map[string]string) {
	r.urlsRewriteMap = urls
}

// asset resolver parts

func (r *Renderer) GetEmojiUrl(code rune) string {
	return fmt.Sprintf("%s/emojies/%x.png", r.Config.AnytypeCdnUrl, code)
}

func (r *Renderer) GetStaticFolderUrl(filepath string) string {
	return fmt.Sprintf("%s%s", r.Config.StaticFilesPath, filepath)
}

func (r *Renderer) GetPrismJsUrl(filepath string) string {
	return fmt.Sprintf("%s%s", r.Config.PrismJsCdnUrl, filepath)
}

func (r *Renderer) Render(writer io.Writer) (err error) {
	defer func() {
		if p := recover(); p != nil {
			stack := string(debug.Stack())
			err = fmt.Errorf("panic: %v, publishFilesPath: %s, stack: %s", p, r.Config.PublishFilesPath, stack)
			log.Error("panic recover", zap.String("where", "Render()"), zap.Error(err), zap.String("stack", stack))
		}
	}()

	err = r.RootComp.Render(context.Background(), writer)
	if err != nil {
		return
	}
	return
}

func (r *Renderer) ReadJsonpbSnapshot(path string) (*pb.SnapshotWithType, error) {
	snapshot, ok := r.CachedPbFiles[path]
	if ok {
		return snapshot, nil
	}

	snapshotStr, ok := r.UberSp.PbFiles[path]
	if !ok {
		err := fmt.Errorf("path %s is not found in snapshot", path)
		return nil, err
	}

	newSnapshot, err := readJsonpbSnapshot(snapshotStr)
	if err != nil {
		return nil, err
	}

	r.CachedPbFiles[path] = &newSnapshot

	return &newSnapshot, nil
}

func (r *Renderer) unwrapLayouts(blocks []*model.Block) []*model.Block {
	ret := make([]*model.Block, 0)
	for _, b := range blocks {
		layout := b.GetLayout()

		if layout == nil || layout.GetStyle() != model.BlockContentLayout_Div {
			ret = append(ret, b)
			continue
		}

		chBlocks := make([]*model.Block, len(b.ChildrenIds))
		for i, chId := range b.ChildrenIds {
			chBlocks[i] = r.BlocksById[chId]
		}

		ret = append(ret, chBlocks...)
	}

	return ret
}

func (r *Renderer) hydrateNumberBlocksInner(blocks []*model.Block) {
	unwrapped := r.unwrapLayouts(blocks)
	prevNumber := 0

	for _, b := range unwrapped {
		if t := b.GetText(); t != nil {
			if t.GetStyle() == model.BlockContentText_Numbered {
				if _, ok := r.BlockNumbers[b.Id]; !ok {
					prevNumber++
					r.BlockNumbers[b.Id] = prevNumber
				}
			} else {
				prevNumber = 0
			}
		}

		layout := b.GetLayout()
		if layout != nil || layout.GetStyle() == model.BlockContentLayout_Div {
			continue
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

func (r *Renderer) getSpaceData() (string, templ.Component) {
	details, err := r.findWorkspaceDetails()
	if err != nil {
		return "", nil
	}
	spaceName := getRelationField(details, bundle.RelationKeyName, relationToString)
	if spaceName == "" {
		spaceName = "Untitled"
	}
	iconTemplate := r.getIconFromDetails(details)
	return spaceName, iconTemplate
}

func Comment(text string) templ.ComponentFunc {
	return func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, "<!--comment:\n")
		if err != nil {
			log.Error("comment error", zap.Error(err))
		}
		_, err = io.WriteString(w, text)
		if err != nil {
			log.Error("comment error", zap.Error(err))
		}
		_, err = io.WriteString(w, "\n-->\n")
		if err != nil {
			log.Error("comment error", zap.Error(err))
		}
		return nil
	}
}
