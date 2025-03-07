package renderer

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

type childBlock struct {
	*model.Block
	isChild bool
}

func (r *Renderer) makeTableOfContentBlockParams(block *model.Block) *BlockParams {
	blockParams := makeDefaultBlockParams(block)

	color := block.GetBackgroundColor()
	if color != "" {
		blockParams.ContentClasses = append(blockParams.ContentClasses, fmt.Sprintf("bgColor bgColor-%s", color))
	}

	blockParams.Content = BlocksWrapper(&BlockWrapperParams{Classes: []string{"wrap"}, Components: r.getList()})
	return blockParams
}

func (r *Renderer) getList() []templ.Component {
	blocks := r.traverseBlocks(r.BlocksById, r.Root.GetId(), false)
	var list []templ.Component
	styles := []model.BlockContentTextStyle{model.BlockContentText_Header1, model.BlockContentText_Header2, model.BlockContentText_Header3}

	var hasH1, hasH2 bool

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
	params := r.makeTableOfContentBlockParams(block)
	return BlockTemplate(r, params)
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
