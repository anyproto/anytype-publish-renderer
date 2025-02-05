package renderer

import (
	"fmt"
	"reflect"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"go.uber.org/zap"
)

func blockContentTypeToName(b *model.Block) string {
	if b == nil {
		log.Error("blockContentTypeToName: block is nil")
		return ""
	}

	switch b.Content.(type) {
	case *model.BlockContentOfText:
		return "Text"
	case *model.BlockContentOfLayout:
		return "Layout"
	case *model.BlockContentOfFeaturedRelations:
		return "Featured"
	case *model.BlockContentOfDiv:
		return "Div"
	case *model.BlockContentOfFile:
		return "File"
	case *model.BlockContentOfTable:
		return "Table"
	case *model.BlockContentOfLatex:
		return "Latex"
	case *model.BlockContentOfBookmark:
		return "Bookmark"
	case *model.BlockContentOfLink:
		return "Link"
	case *model.BlockContentOfRelation:
		return "Relation"
	case *model.BlockContentOfTableOfContents:
		return "TableOfContents"
	default:
		log.Error("blockContentTypeToName: unkonwn block type", zap.String("type", reflect.TypeOf(b.Content).String()))
		return ""
	}

}

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
	if blockType := blockContentTypeToName(b); blockType != "" {
		classes = append(classes, fmt.Sprintf("block%s", blockType))
	}
	classes = append(classes, align)

	return &BlockParams{
		Id:          b.Id,
		Classes:     classes,
		ChildrenIds: b.ChildrenIds,
	}
}
