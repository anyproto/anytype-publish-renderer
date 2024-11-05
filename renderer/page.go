package renderer

import (
	"reflect"

	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"

	"github.com/a-h/templ"
	"go.uber.org/zap"
)

func (r *Renderer) RenderPage() templ.Component {
	return PageTemplate(r)
}

func (r *Renderer) RenderBlock(b *model.Block) templ.Component {
	log.Debug("block type", zap.String("type", reflect.TypeOf(b.Content).String()))
	switch b.Content.(type) {
	case *model.BlockContentOfText:
		return r.RenderText(b)
	case *model.BlockContentOfFile:
	case *model.BlockContentOfBookmark:
	case *model.BlockContentOfDiv:
	case *model.BlockContentOfLayout:
		return r.RenderLayout(b)
	case *model.BlockContentOfLink:
	case *model.BlockContentOfTable:
	case *model.BlockContentOfSmartblock:

	default:

	}

	log.Warn("block is not supported", zap.String("type", reflect.TypeOf(b.Content).String()))
	return TextTemplate(r, "block is not supported")
}
