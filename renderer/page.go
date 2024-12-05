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

func (r *Renderer) hasIconAndCover() bool {
	fields := r.Sp.Snapshot.Data.GetDetails()
	coverId := pbtypes.GetString(fields, "coverId")
	if coverId == "" {
		return false
	}

	_, err := r.getFileUrl(coverId)

	return (err == nil)

}

func (r *Renderer) MakeRenderPageParams() (params *RenderPageParams) {
	var classes string
	if r.hasIconAndCover() {
		classes = "withIconAndCover"
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

func (r *Renderer) RenderBlock(b *model.Block) templ.Component {
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
	case *model.BlockContentOfBookmark:
	case *model.BlockContentOfLink:
	case *model.BlockContentOfSmartblock:
	default:

	}

	log.Warn("block is not supported",
		zap.String("type", reflect.TypeOf(b.Content).String()),
		zap.String("id", b.Id))
	return NoneTemplate(fmt.Sprintf("not supported: %s, %s", b.Id, reflect.TypeOf(b.Content).String()))
}

func (r *Renderer) joinSpaceLink() templ.SafeURL {
	return templ.SafeURL(r.AssetResolver.GetJoinSpaceLink())
}
