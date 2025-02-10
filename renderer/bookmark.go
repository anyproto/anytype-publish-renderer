package renderer

import (
	"net/url"

	"github.com/a-h/templ"
	"github.com/gogo/protobuf/types"
	"go.uber.org/zap"

	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

func (r *Renderer) MakeBookmarkRendererParams(b *model.Block) *BlockParams {
	bookmark := b.GetBookmark()
	if bookmark.GetUrl() == "" {
		return nil
	}

	targetObjectId := bookmark.GetTargetObjectId()
	targetBookmark := r.getObjectSnapshot(targetObjectId)
	if targetBookmark == nil {
		return nil
	}

	details := targetBookmark.GetSnapshot().GetData().GetDetails()
	if details == nil || len(details.GetFields()) == 0 {
		return nil
	}

	parsedUrl, err := url.Parse(bookmark.GetUrl())
	if err != nil {
		log.Error("failed to parse bookmark url", zap.Error(err))
		return nil
	}
	return r.getBookmarkBlockParams(b, details, parsedUrl)
}

func (r *Renderer) getBookmarkBlockParams(b *model.Block, details *types.Struct, parsedUrl *url.URL) *BlockParams {
	bgColor := b.GetBackgroundColor()
	innerClasses := []string{"inner"}

	if bgColor != "" {
		innerClasses = append(innerClasses, "bgColor", "bgColor-"+bgColor)
	}

	sideLeft := r.getSideLeftComponent(details, parsedUrl)
	sideRightComponents, innerClasses := r.getSideRightComponent(details, innerClasses)
	blockParams := makeDefaultBlockParams(b)
	blockParams.Content = BookmarkLinkTemplate(templ.URL(b.GetBookmark().GetUrl()), innerClasses, []templ.Component{sideLeft, sideRightComponents})
	return blockParams
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
	params := r.MakeBookmarkRendererParams(b)
	if params == nil {
		return NoneTemplate("")
	}
	return BlockTemplate(r, params)
}
