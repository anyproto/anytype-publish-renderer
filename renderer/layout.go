package renderer

import (
	"fmt"
	"github.com/a-h/templ"

	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
)

func (r *Renderer) RenderLayout(b *model.Block) templ.Component {
	blockParams := makeDefaultBlockParams(b)
	fields := b.GetFields()
	width := fmt.Sprintf("%.2f", pbtypes.GetFloat64(fields, "width"))
	blockParams.Width = width
	return BlockTemplate(r, blockParams)
}
