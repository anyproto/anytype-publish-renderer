package renderer

import (
	"strings"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

type DivRenderParams struct {
	Classes string
	Id      string
	Comp    templ.Component
}

func (r *Renderer) MakeRenderDivParams(b *model.Block) (params *DivRenderParams) {
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

	classes := []string{"block", "blockDiv", divClass}
	params = &DivRenderParams{
		Id:      "block-" + b.Id,
		Classes: strings.Join(classes, " "),
		Comp:    comp,
	}

	return
}

func (r *Renderer) RenderDiv(b *model.Block) templ.Component {
	params := r.MakeRenderDivParams(b)
	return DivTemplate(params)
}
