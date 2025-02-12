package renderer

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

type TableOfContentRenderParams struct {
	Id              string
	BackgroundColor string
	Items           []templ.Component
	IsEmpty         bool
}

type childBlock struct {
	*model.Block
	isChild bool
}

func (r *Renderer) MakeTableOfContentRenderParams(block *model.Block) *TableOfContentRenderParams {
	blockId := block.GetId()
	params := &TableOfContentRenderParams{
		Id: blockId,
	}
	color := block.GetBackgroundColor()
	if color != "" {
		params.BackgroundColor = fmt.Sprintf("bgColor bgColor-%s", color)
	}

	params.Items = r.getList()
	if len(params.Items) == 0 {
		params.IsEmpty = true
	}
	return params
}

func (r *Renderer) getList() []templ.Component {
	blocks := r.traverseBlocks(r.BlocksById, r.Root.GetId(), false)
	list := []templ.Component{}
	styles := []model.BlockContentTextStyle{model.BlockContentText_Header1, model.BlockContentText_Header2, model.BlockContentText_Header3}

	hasH1 := false
	hasH2 := false

	for _, bl := range blocks {
		text := bl.GetText()

		if text == nil {
			continue
		}

		style := text.GetStyle()

		if !contains(styles, style) {
			continue
		}

		depth := 0

		if style == model.BlockContentText_Header1 {
			depth = 0
			hasH1 = true
			hasH2 = false
		}

		if style == model.BlockContentText_Header2 {
			hasH2 = true
			if hasH1 {
				depth++
			}
		}

		if style == model.BlockContentText_Header3 {
			if hasH1 {
				depth++
			}
			if hasH2 {
				depth++
			}
		}

		name := r.getHeadingName(bl.Block)
		list = append(list, HeadingTemplate(bl.GetId(), name, depth))
	}

	return list
}

func (r *Renderer) getHeadingName(b *model.Block) string {
	text := b.GetText().GetText()
	if text == "" {
		text = defaultName
	}
	return text
}

func (r *Renderer) RenderTableOfContent(block *model.Block) templ.Component {
	params := r.MakeTableOfContentRenderParams(block)
	return TableOfContentTemplate(params)
}

func (r *Renderer) traverseBlocks(blockMap map[string]*model.Block, blockID string, isChild bool) []*childBlock {
	var result []*childBlock
	if block, exists := blockMap[blockID]; exists {
		result = append(result, &childBlock{Block: block, isChild: isChild})
		for _, childID := range block.ChildrenIds {
			result = append(result, r.traverseBlocks(blockMap, childID, block.GetId() != r.Root.GetId())...)
		}
	}
	return result
}

func contains(arr []model.BlockContentTextStyle, target any) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}
