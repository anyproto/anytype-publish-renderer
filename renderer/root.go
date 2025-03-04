package renderer

import (
	"fmt"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"

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

	if width == 0 {
		width = getRelationField(r.ObjectTypeDetails, bundle.RelationKeyLayoutWidth, relationToFloat64)
	}

	min := "60%"
	w := fmt.Sprintf("%f", width)

	str := "max(" + min + ", min(calc(100% - 96px), calc(" + min + " + (100% - " + min + " - 96px) * " + w + ")))"
	style := fmt.Sprintf(`
		<style> 
			.blocks {
				width: %s;
			}
		</style> 
	`, str)

	return &RootRenderParams{
		Style: style,
	}
}
func (r *Renderer) getStyle(params *RootRenderParams) templ.Component {
	return templ.Raw(params.Style)
}

func (r *Renderer) RenderRoot() templ.Component {
	params := r.makeRootRenderParams(r.Root)
	return RootTemplate(r, params)
}
