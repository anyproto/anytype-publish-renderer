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

func (r *Renderer) makeRootRenderParams(b *model.Block) (params *RootRenderParams) {
	fields := b.Fields
	var width float64
	if fields != nil && fields.Fields != nil && fields.Fields["width"] != nil {
		width = fields.Fields["width"].GetNumberValue()
	}
	params = &RootRenderParams{}
	if width == 0 {
		return params
	}
	widthPercentage := strconv.Itoa(int(width*100)) + "%"
	params.Style = fmt.Sprintf(`
<style> 
.blocks {
	width: %s
}
</style> 
`, widthPercentage)
	return
}
func (r *Renderer) getStyle(params *RootRenderParams) templ.Component {
	return templ.Raw(params.Style)
}

func (r *Renderer) RenderRoot() templ.Component {
	params := r.makeRootRenderParams(r.Root)
	return RootTemplate(r, params)
}
