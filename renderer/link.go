package renderer

import (
	"fmt"
	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

const linkTemplate = "anytype://object?objectId=%s&spaceId=%s"

type LinkRenderParams struct {
	Link templ.Component
}

func (r *Renderer) MakeLinkRenderParams(b *model.Block) (params *LinkRenderParams) {
	targetObjectId := b.GetLink().GetTargetBlockId()
	spaceId := r.Sp.GetSnapshot().GetData().GetDetails().GetFields()[bundle.RelationKeySpaceId.String()].GetStringValue()
	link := fmt.Sprintf(linkTemplate, targetObjectId, spaceId)
	linkComponent := fmt.Sprintf(`<a href="%s"> %s </a>`, link, link)
	params = &LinkRenderParams{Link: templ.Raw(linkComponent)}
	return

}
func (r *Renderer) RenderLink(b *model.Block) templ.Component {
	params := r.MakeLinkRenderParams(b)
	return TextTemplate(r, &TextRenderParams{
		Id:        b.Id,
		InnerFlex: []templ.Component{PlainTextWrapTemplate(params.Link)},
	})
}
