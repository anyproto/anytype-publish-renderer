package renderer

import (
	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"go.uber.org/zap"
)

type BookmarkRendererParams struct {
	Id          string
	IsEmpty     bool
	Url         string
	Favicon     string
	Name        string
	Description string
	Image       string
}

func (r *Renderer) MakeBookmarkRendererParams(b *model.Block) (params *BookmarkRendererParams) {
	bookmark := b.GetBookmark()
	if bookmark.GetUrl() == "" {
		return &BookmarkRendererParams{IsEmpty: true}
	}
	favicon, err := r.getFileUrl(bookmark.GetFaviconHash())
	if err != nil {
		log.Error("Failed to get bookmark favicon url", zap.Error(err))
	}
	image, err := r.getFileUrl(bookmark.GetImageHash())
	if err != nil {
		log.Error("Failed to get bookmark favicon url", zap.Error(err))
	}
	return &BookmarkRendererParams{
		Id:          b.Id,
		Url:         bookmark.GetUrl(),
		Favicon:     favicon,
		Name:        bookmark.GetTitle(),
		Description: bookmark.GetDescription(),
		Image:       image,
	}
}

func (r *Renderer) RenderBookmark(b *model.Block) templ.Component {
	params := r.MakeBookmarkRendererParams(b)
	return BookmarkTempl(params)
}
