package renderer

import (
	"github.com/a-h/templ"

	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

func (r *Renderer) RenderLayout(b *model.Block) templ.Component {
	blockParams := r.makeLayoutBlockParams(b)
	return BlockTemplate(r, blockParams)
}

func (r *Renderer) makeLayoutBlockParams(b *model.Block) *BlockParams {
	blockParams := makeDefaultBlockParams(b)
	blockParams.Width = GetWidth(b.GetFields())
	blockParams.Classes = append(blockParams.Classes, "layout"+b.GetLayout().GetStyle().String())
	return blockParams
}
