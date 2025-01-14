package renderer

import (
	"context"
	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/stretchr/testify/assert"
	"html/template"
	"testing"
)

func TestRenderer_MakeLinkRenderParams(t *testing.T) {
	r := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
	id := "blockId"
	expected := "<a href=\"anytype://object?objectId=targetBlockId&spaceId=bafyreiholtkdzlvc5ahtgzgbb3ftyszrpad6swilhkfzrgnvsah2rz6zke.35ssi7ciufxuc\"> anytype://object?objectId=targetBlockId&spaceId=bafyreiholtkdzlvc5ahtgzgbb3ftyszrpad6swilhkfzrgnvsah2rz6zke.35ssi7ciufxuc </a>"
	// when
	params := r.MakeLinkRenderParams(&model.Block{
		Id: id,
		Content: &model.BlockContentOfLink{Link: &model.BlockContentLink{
			TargetBlockId: "targetBlockId",
		}},
	})

	// then
	html, err := templ.ToGoHTML(context.Background(), params.Link)
	assert.NoError(t, err)
	assert.Equal(t, template.HTML(expected), html)
}
