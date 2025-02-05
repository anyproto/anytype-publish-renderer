package renderer

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"

	"github.com/a-h/templ"
	"go.uber.org/zap"
)

type RenderPageParams struct {
	Classes     string
	Name        string
	Description string
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

	return err == nil

}

func (r *Renderer) hasPageCover() bool {
	fields := r.Sp.Snapshot.Data.GetDetails()
	coverType, err := ToCoverType(pbtypes.GetInt64(fields, "coverType"))
	if err != nil {
		return false
	}
	coverId := pbtypes.GetString(fields, "coverId")
	if coverId != "" {
		switch coverType {
		case CoverType_Image, CoverType_Source:
			_, err := r.getFileUrl(coverId)
			return err == nil
		default:
			return true
		}
	}
	return false
}

func (r *Renderer) MakeRenderPageParams() (params *RenderPageParams) {
	fields := r.Sp.Snapshot.Data.GetDetails()
	layoutAlign := pbtypes.GetInt64(fields, "layoutAlign")
	classes := []string{"blocks", fmt.Sprintf("layoutAlign%d", layoutAlign)}
	name := pbtypes.GetString(fields, "name")
	description := pbtypes.GetString(fields, "description")
	snippet := pbtypes.GetString(fields, "snippet")

	hasPageIcon := r.hasPageIcon()
	hasPageCover := r.hasPageCover()

	class := ""
	switch {
	case hasPageIcon && hasPageCover:
		class = "withIconAndCover"
	case hasPageIcon:
		class = "withIcon"
	case hasPageCover:
		class = "withCover"
	}

	classes = append(classes, class)

	descr := description
	if descr == "" {
		descr = snippet
	}

	return &RenderPageParams{
		Classes:	 strings.Join(classes, " "),
		Name:   	 name,
		Description: descr,
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
		return r.RenderFeaturedRelations(b)
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
