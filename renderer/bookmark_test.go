package renderer

import (
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testId = "testId"

func TestMakeBookmarkRendererParams(t *testing.T) {
	t.Run("empty bookmark", func(t *testing.T) {
		// given
		r := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
		bookmark := &model.Block{Id: testId, Content: &model.BlockContentOfBookmark{Bookmark: &model.BlockContentBookmark{}}}

		// when
		params := r.MakeBookmarkRendererParams(bookmark)

		// then
		assert.NotNil(t, params)
		assert.True(t, params.IsEmpty)
		assert.Empty(t, params.Url)
	})

	t.Run("non empty bookmark", func(t *testing.T) {
		r := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
		bookmark := &model.Block{Id: testId, Content: &model.BlockContentOfBookmark{Bookmark: &model.BlockContentBookmark{
			Url:         "https://example.com",
			FaviconHash: "bafyreighh5qn3qr4wcpq4n6k7imawkaytrxksg7te4gpxqrlruwmezhjii",
			ImageHash:   "bafyreighh5qn3qr4wcpq4n6k7imawkaytrxksg7te4gpxqrlruwmezhjii",
			Title:       "Example",
			Description: "An example bookmark",
		}}}

		// when
		params := r.MakeBookmarkRendererParams(bookmark)

		// then
		assert.NotNil(t, params)
		assert.False(t, params.IsEmpty)
		assert.Equal(t, "https://example.com", params.Url)
		assert.Equal(t, "../test_snapshots/Anytype.WebPublish.20241217.112212.67/files/img_5296.jpeg", params.Favicon)
		assert.Equal(t, "Example", params.Name)
		assert.Equal(t, "An example bookmark", params.Description)
		assert.Equal(t, "../test_snapshots/Anytype.WebPublish.20241217.112212.67/files/img_5296.jpeg", params.Image)
	})

	t.Run("Favicon and Image Errors", func(t *testing.T) {
		r := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
		bookmark := &model.Block{Id: testId, Content: &model.BlockContentOfBookmark{Bookmark: &model.BlockContentBookmark{
			Url:         "https://example.com",
			FaviconHash: "favicon_hash",
			ImageHash:   "image_hash",
			Title:       "Example",
			Description: "An example bookmark",
		}}}

		// when
		params := r.MakeBookmarkRendererParams(bookmark)

		// then
		assert.NotNil(t, params)
		assert.False(t, params.IsEmpty)
		assert.Equal(t, "https://example.com", params.Url)
		assert.Equal(t, "", params.Favicon)
		assert.Equal(t, "Example", params.Name)
		assert.Equal(t, "An example bookmark", params.Description)
		assert.Equal(t, "", params.Image)
	})
}
