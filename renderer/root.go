package renderer

import (
	"fmt"
	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"strconv"
)

type RootRenderParams struct {
	Style string
}

func (r *Renderer) MakeRootRenderParams(b *model.Block) (params *RootRenderParams) {
	fields := b.Fields
	var width float64
	if fields != nil && fields.Fields != nil && fields.Fields["width"] != nil {
		width = fields.Fields["width"].GetNumberValue()
	}
	params = &RootRenderParams{}
	params.Style = fmt.Sprintf("style={\"width:\" + %s}", strconv.Itoa(int(width*100)))
	return

}
func (r *Renderer) RenderRoot(b *model.Block) templ.Component {
	params := r.MakeRootRenderParams(b)
	return RootTemplate(r, params)
}
