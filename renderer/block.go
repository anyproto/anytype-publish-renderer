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
		if isInlineLink(b) {
			return "File"
		} else {
			fileClass := getFileClass(b)
			return "Media " + fileClass
		}

	case *model.BlockContentOfTable:
		return "Table"
	case *model.BlockContentOfLatex:
		return "Embed"
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
	Id                string
	BlockType         string
	Classes           []string
	ContentClasses    []string
	AdditionalClasses []string
	Content           templ.Component
	Additional        templ.Component
	ChildrenIds       []string
	Width             string
}

type BlockWrapperParams struct {
	Classes    []string
	Width      string
	Components []templ.Component
}

func makeDefaultBlockParams(b *model.Block) *BlockParams {
	a := b.GetAlign()
	align := fmt.Sprintf("align%d", a)
	classes := []string{"block", align}

	if blockType := blockContentTypeToName(b); blockType != "" {
		classes = append(classes, fmt.Sprintf("block%s", blockType))
	}
	return &BlockParams{
		Id:          b.Id,
		Classes:     classes,
		ChildrenIds: b.ChildrenIds,
	}
}

func makeWrappedBlockParams(b *model.Block, t templ.Component) *BlockParams {
	blockParams := makeDefaultBlockParams(b)

	bgColor := b.GetBackgroundColor()
	if bgColor != "" {
		blockParams.Classes = append(blockParams.Classes, "bgColor", "bgColor-"+bgColor)
	}

	blockParams.Content = t
	blockParams.Width = GetWidth(b.Fields)

	return blockParams

}
