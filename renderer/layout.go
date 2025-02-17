package renderer

import (
	"fmt"
	"github.com/a-h/templ"

	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
)

func (r *Renderer) RenderLayout(b *model.Block) templ.Component {
	blockParams := r.makeLayoutBlockParams(b)
	return BlockTemplate(r, blockParams)
}

func (r *Renderer) makeLayoutBlockParams(b *model.Block) *BlockParams {
	blockParams := makeDefaultBlockParams(b)
	fields := b.GetFields()
	width := fmt.Sprintf("%.2f", pbtypes.GetFloat64(fields, "width"))
	blockParams.Width = width
	blockParams.Classes = append(blockParams.Classes, "layout"+b.GetLayout().GetStyle().String())
	return blockParams
}
