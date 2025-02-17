package renderer

import (
	"net/url"

	"github.com/a-h/templ"
	"github.com/gogo/protobuf/types"
	"go.uber.org/zap"

	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
)

func (r *Renderer) makeBookmarkBlockParams(b *model.Block) *BlockParams {
	bookmark := b.GetBookmark()

	details := r.getBookmarkDetails(bookmark)
	return r.getBookmarkBlockParams(b, details)
}

func (r *Renderer) getBookmarkDetails(bookmark *model.BlockContentBookmark) *types.Struct {
	targetObjectId := bookmark.GetTargetObjectId()
	targetBookmark := r.getObjectSnapshot(targetObjectId)
	if targetBookmark == nil {
		// fallback to old logic with block
		return r.getDetailsFromBlock(bookmark)
	}

	details := targetBookmark.GetSnapshot().GetData().GetDetails()
	if details == nil || len(details.GetFields()) == 0 {
		// fallback to old logic with block
		return r.getDetailsFromBlock(bookmark)
	}
	return details
}

func (r *Renderer) getDetailsFromBlock(bookmark *model.BlockContentBookmark) *types.Struct {
	return &types.Struct{Fields: map[string]*types.Value{
		bundle.RelationKeyIconImage.String():   pbtypes.String(bookmark.GetFaviconHash()),
		bundle.RelationKeyPicture.String():     pbtypes.String(bookmark.GetImageHash()),
		bundle.RelationKeyDescription.String(): pbtypes.String(bookmark.GetDescription()),
		bundle.RelationKeyName.String():        pbtypes.String(bookmark.GetTitle()),
		bundle.RelationKeySource.String():      pbtypes.String(bookmark.GetUrl()),
	}}
}

func (r *Renderer) getBookmarkBlockParams(b *model.Block, details *types.Struct) *BlockParams {
	bookmarkUrl := r.getBookmarkUrl(b, details)
	if bookmarkUrl == "" {
		return nil
	}
	parsedUrl, err := url.Parse(bookmarkUrl)
	if err != nil {
		log.Error("failed to parse bookmark url", zap.Error(err))
		return nil
	}
	bgColor := b.GetBackgroundColor()
	innerClasses := []string{"inner"}

	if bgColor != "" {
		innerClasses = append(innerClasses, "bgColor", "bgColor-"+bgColor)
	}

	sideLeft := r.getSideLeftComponent(details, parsedUrl)
	sideRightComponents, innerClasses := r.getSideRightComponent(details, innerClasses)
	blockParams := makeDefaultBlockParams(b)
	blockParams.Content = BookmarkLinkTemplate(templ.URL(bookmarkUrl), innerClasses, []templ.Component{sideLeft, sideRightComponents})
	return blockParams
}

func (r *Renderer) getBookmarkUrl(b *model.Block, details *types.Struct) string {
	bookmarkUrl := getRelationField(details, bundle.RelationKeySource, relationToString)
	if bookmarkUrl == "" {
		return b.GetBookmark().GetUrl()
	}
	return bookmarkUrl
}

func (r *Renderer) getSideLeftComponent(details *types.Struct, parsedUrl *url.URL) templ.Component {
	var (
		sideLeftsComponents []templ.Component
		linkParams          = &BlockWrapperParams{Classes: []string{"link"}}
	)
	icon := getRelationField(details, bundle.RelationKeyIconImage, r.relationToFileUrl)
	if icon != "" {
		linkParams.Components = append(linkParams.Components, ImageWithSourceTemplate(icon, "fav"))
	}
	linkParams.Components = append(linkParams.Components, templ.Raw(parsedUrl.Host))
	sideLeftsComponents = append(sideLeftsComponents, BlocksWrapper(linkParams))
	description := getRelationField(details, bundle.RelationKeyDescription, relationToString)
	name := getRelationField(details, bundle.RelationKeyName, relationToString)
	sideLeftsComponents = append(sideLeftsComponents, BasicTemplate("name", name), BasicTemplate("descr", description))
	wrapper := BlocksWrapper(&BlockWrapperParams{
		Classes:    []string{"side left"},
		Components: sideLeftsComponents,
	})
	return wrapper
}

func (r *Renderer) getSideRightComponent(details *types.Struct, innerClasses []string) (templ.Component, []string) {
	sideRightComponent := &BlockWrapperParams{Classes: []string{"side right"}}
	image := getRelationField(details, bundle.RelationKeyPicture, r.relationToFileUrl)
	if image != "" {
		innerClasses = append(innerClasses, "withImage")
		sideRightComponent.Components = append(sideRightComponent.Components, ImageWithSourceTemplate(image, "img"))
	}
	return BlocksWrapper(sideRightComponent), innerClasses
}

func (r *Renderer) RenderBookmark(b *model.Block) templ.Component {
	params := r.makeBookmarkBlockParams(b)
	if params == nil {
		return NoneTemplate("")
	}
	return BlockTemplate(r, params)
}
