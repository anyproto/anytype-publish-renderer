package renderer

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"github.com/gogo/protobuf/types"
	"go.uber.org/zap"
)

type RenderPageParams struct {
	Classes       string
	HeaderClasses string
	Name          string
	Description   string
	SpaceLink     templ.SafeURL
	SpaceIcon     templ.Component
	OgImageUrl    string
	SpaceName     string
}

func (r *Renderer) getOgImageUrl() string {
	defaultImage := r.GetStaticFolderUrl("/img/og-image.png")
	details := r.Sp.Snapshot.Data.GetDetails()
	iconImage := getRelationField(details, bundle.RelationKeyIconImage, r.relationToFileUrl)
	if iconImage != "" {
		return iconImage
	}

	iconEmoji := getRelationField(details, bundle.RelationKeyIconEmoji, r.relationToEmojiUrl)
	if iconEmoji != "" {
		return iconEmoji
	}

	coverType, err := ToCoverType(pbtypes.GetInt64(details, "coverType"))

	if err != nil {
		return defaultImage
	}

	coverId := pbtypes.GetString(details, "coverId")
	if coverId != "" {
		switch coverType {
		case CoverType_Image, CoverType_Source:
			coverImageUrl, _ := r.getFileUrl(coverId)
			if coverImageUrl != "" {
				return coverImageUrl
			}
		}
	}
	return defaultImage

}

func (r *Renderer) hasPageIcon() bool {
	details := r.Sp.Snapshot.Data.GetDetails()
	iconEmoji := getRelationField(details, bundle.RelationKeyIconEmoji, r.relationToEmojiUrl)
	iconImage := getRelationField(details, bundle.RelationKeyIconImage, r.relationToFileUrl)

	if isTodoLayout(r.ResolvedLayout) {
		return false
	}

	if iconEmoji != "" {
		return true
	}

	if iconImage == "" {
		return false
	}

	return true
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

	classes := []string{"blocks", fmt.Sprintf("layoutAlign%d", r.LayoutAlign)}
	headerClasses := []string{"header"}
	name := pbtypes.GetString(fields, "name")
	description := pbtypes.GetString(fields, "description")
	snippet := pbtypes.GetString(fields, "snippet")
	spaceLink := r.joinSpaceLink()
	spaceName, spaceIcon := r.getSpaceData()

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

	classes = append(classes, class, getLayoutClass(r.ResolvedLayout))

	descr := description
	if descr == "" {
		descr = snippet
	}

	if spaceLink != "" {
		headerClasses = append(headerClasses, "withJoinSpace")
	}

	return &RenderPageParams{
		Classes:       strings.Join(classes, " "),
		HeaderClasses: strings.Join(headerClasses, " "),
		Name:          name,
		Description:   descr,
		SpaceLink:     spaceLink,
		SpaceIcon:     spaceIcon,
		SpaceName:     spaceName,
		OgImageUrl:    r.getOgImageUrl(),
	}
}

func (r *Renderer) fillLayoutAlign(details *types.Struct) {
	if value, ok := details.GetFields()[bundle.RelationKeyLayoutAlign.String()]; ok {
		r.LayoutAlign = int64(value.GetNumberValue())
		return
	}
	r.LayoutAlign = pbtypes.GetInt64(r.ObjectTypeDetails, bundle.RelationKeyLayoutAlign.String())
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

func (r *Renderer) supportLink() templ.SafeURL {
	supportEmail := "support@anytype.io"
	subject := "subject=Web Publishing Report"
	body := fmt.Sprintf("body=PublishFilesPath: %s", r.Config.PublishFilesPath)
	mailtoUrl := fmt.Sprintf("mailto:%s?%s&%s", supportEmail, subject, body)
	return templ.SafeURL(mailtoUrl)
}

func (r *Renderer) joinSpaceLink() templ.SafeURL {
	return templ.URL(r.UberSp.Meta.InviteLink)
}
