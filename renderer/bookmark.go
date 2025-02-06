package renderer

import (
	"html"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"go.uber.org/zap"
)

type BookmarkRendererParams struct {
	Id          string
	Classes     string
	SideLeftClasses string
	IsEmpty     bool
	Url         string
	Favicon     string
	Name        string
	Description string
	Image       string
	SafeUrl     templ.SafeURL
}

func (r *Renderer) MakeBookmarkRendererParams(b *model.Block) (params *BookmarkRendererParams) {
	bookmark := b.GetBookmark()
	bgColor := b.GetBackgroundColor()
	classes := []string{"block", "blockBookmark"}
	sideLeftClasses := []string{"side", "left"}

	if bgColor != "" {
		sideLeftClasses = append(sideLeftClasses, "bgColor", "bgColor-" + bgColor)
	}

	if bookmark.GetUrl() == "" {
		return &BookmarkRendererParams{IsEmpty: true}
	}

	targetObjectId := bookmark.GetTargetObjectId()
	targetBookmark, err := r.ReadJsonpbSnapshot(filepath.Join("objects", targetObjectId+".pb"))
	if err != nil {
		return &BookmarkRendererParams{IsEmpty: true}
	}

	details := targetBookmark.GetSnapshot().GetData().GetDetails()
	if details == nil || len(details.GetFields()) == 0 {
		return &BookmarkRendererParams{IsEmpty: true}
	}

	var (
		favicon, image, description, name string
	)

	if icon := details.Fields[bundle.RelationKeyIconImage.String()]; icon != nil && icon.GetStringValue() != "" {
		favicon, err = r.getFileUrl(icon.GetStringValue())
		if err != nil {
			log.Error("failed to get bookmark favicon url", zap.Error(err))
		}
	}

	if picture := details.Fields[bundle.RelationKeyPicture.String()]; picture != nil && picture.GetStringValue() != "" {
		image, err = r.getFileUrl(picture.GetStringValue())
		if err != nil {
			log.Error("failed to get bookmark image url", zap.Error(err))
		}
	}

	if descriptionValue := details.Fields[bundle.RelationKeyDescription.String()]; descriptionValue != nil && descriptionValue.GetStringValue() != "" {
		description = descriptionValue.GetStringValue()
	}

	if nameValue := details.Fields[bundle.RelationKeyName.String()]; nameValue != nil && nameValue.GetStringValue() != "" {
		name = nameValue.GetStringValue()
	}

	parsedUrl, err := url.Parse(bookmark.GetUrl())
	if err != nil {
		log.Error("failed to parse bookmark url", zap.Error(err))
		return &BookmarkRendererParams{IsEmpty: true}
	}

	return &BookmarkRendererParams{
		Id:          b.Id,
		Classes:     strings.Join(classes, " "),
		SideLeftClasses: strings.Join(sideLeftClasses, " "),
		Url:         parsedUrl.Host,
		Favicon:     favicon,
		Name:        html.UnescapeString(name),
		Description: html.UnescapeString(description),
		Image:       image,
		SafeUrl:     templ.SafeURL(bookmark.GetUrl()),
	}
}

func (r *Renderer) RenderBookmark(b *model.Block) templ.Component {
	params := r.MakeBookmarkRendererParams(b)

	if params.IsEmpty {
		return NoneTemplate("")
	} else {
		return BookmarkTempl(params)
	}
}
