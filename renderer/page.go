package renderer

import (
	"fmt"
	"reflect"

	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"

	"github.com/a-h/templ"
	"go.uber.org/zap"
)

type RenderPageParams struct {
	Classes string
}

func (r *Renderer) hasPageIcon() bool {
	fields := r.Sp.Snapshot.Data.GetDetails()
	iconEmoji := pbtypes.GetString(fields, "iconEmoji")
	if iconEmoji != "" {
		return true
	}

	iconImageId := pbtypes.GetString(fields, "iconImage")
	if iconImageId == "" {
		return false
	}

	_, err := r.getFileUrl(iconImageId)

	return (err == nil)

}

func (r *Renderer) MakeRenderPageParams() (params *RenderPageParams) {
	var classes string
	if r.hasPageIcon() {
		classes = "hasPageIcon"
	}
	return &RenderPageParams{
		Classes: classes,
	}
}

func (r *Renderer) RenderPage() templ.Component {
	log.Debug("root type", zap.String("type", reflect.TypeOf(r.Root.Content).String()))
	params := r.MakeRenderPageParams()
	return PageTemplate(r, params)
}

func (r *Renderer) RenderBlock(blockId string) templ.Component {
	b, ok := r.BlocksById[blockId]
	if !ok || b == nil {
		log.Error("unexpected nil block", zap.String("blockId", blockId))
		return NoneTemplate(fmt.Sprintf("unexpected nil block: %s", blockId))
	}
	if b.Content == nil {
		log.Error("unexpected nil block.Content")
		return NoneTemplate(fmt.Sprintf("unexpected nil block.Content. block.id: %s", blockId))
	}
	log.Debug("block type",
		zap.String("type", reflect.TypeOf(b.Content).String()),
		zap.String("id", b.Id))

	switch b.Content.(type) {
	case *model.BlockContentOfText:
		return r.RenderText(b)
	case *model.BlockContentOfLayout:
		return r.RenderLayout(b)
	case *model.BlockContentOfFeaturedRelations:
		return NoneTemplate("")
	case *model.BlockContentOfDiv:
		return r.RenderDiv(b)
	case *model.BlockContentOfFile:
		return r.RenderFile(b)
	case *model.BlockContentOfTable:
		return r.RenderTable(b)
	case *model.BlockContentOfLatex:
		return r.RenderEmbed(b)
	case *model.BlockContentOfBookmark:
		return r.RenderBookmark(b)
	case *model.BlockContentOfLink:
		return r.RenderLink(b)
	case *model.BlockContentOfSmartblock:
	case *model.BlockContentOfRelation:
		return r.RenderRelations(b)
	case *model.BlockContentOfTableOfContents:
		return r.RenderTableOfContent(b)
	default:

	}

	log.Warn("block is not supported",
		zap.String("type", reflect.TypeOf(b.Content).String()),
		zap.String("id", b.Id))
	return NoneTemplate(fmt.Sprintf("not supported: %s, %s", b.Id, reflect.TypeOf(b.Content).String()))
}

func (r *Renderer) joinSpaceLink() templ.SafeURL {
	return templ.SafeURL(r.UberSp.Meta.InviteLink)
}

func (r *Renderer) titleText() string {
	titleBlock, ok := r.BlocksById["title"]
	if !ok {
		return ""
	}

	return titleBlock.GetText().Text
}
