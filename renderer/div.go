package renderer

import (
	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

func (r *Renderer) MakeRenderDivParams(b *model.Block) (params *BlockParams) {
	var divClass string
	var comp templ.Component

	switch b.GetDiv().Style {
	case model.BlockContentDiv_Line:
		divClass = "divLine"
		comp = DivLineTemplate()
	case model.BlockContentDiv_Dots:
		divClass = "divDot"
		comp = DivDotTemplate()
	}

	params = makeDefaultBlockParams(b)
	params.Classes = append(params.Classes, divClass)
	params.Content = comp

	return
}

func (r *Renderer) RenderDiv(b *model.Block) templ.Component {
	params := r.MakeRenderDivParams(b)
	return BlockTemplate(r, params)
}
