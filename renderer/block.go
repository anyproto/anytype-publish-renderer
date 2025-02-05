package renderer

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

type BlockParams struct {
	Id          string
	BlockType   string
	Classes     []string
	Content     templ.Component
	Additional  templ.Component
	ChildrenIds []string
}

func makeDefaultBlockParams(b *model.Block) *BlockParams {
	classes := []string{"block"}
	a := b.GetAlign()
	align := fmt.Sprintf("align%d", a)
	classes = append(classes, align)

	return &BlockParams{
		Id:          b.Id,
		Classes:     classes,
		ChildrenIds: b.ChildrenIds,
	}
}
