package renderer

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
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

	w := "max(60%, min(calc(100% - 96px), 60% + (40% - 96px) * " + fmt.Sprintf("%f", width) + "))";

	params.Style = fmt.Sprintf(`
		<style> 
			.blocks {
				width: %s;
			}
		</style> 
	`, w)
	return
}
func (r *Renderer) getStyle(params *RootRenderParams) templ.Component {
	return templ.Raw(params.Style)
}

func (r *Renderer) RenderRoot() templ.Component {
	params := r.makeRootRenderParams(r.Root)
	return RootTemplate(r, params)
}
