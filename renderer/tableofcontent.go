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

	blocks := r.traverseBlocks(r.BlocksById, r.Root.GetId(), false)
	var tableOfContentItems []templ.Component
	for _, bl := range blocks {
		tableOfContentItems = r.retrieveTableOfContentItem(bl, tableOfContentItems)
	}
	params.Items = tableOfContentItems
	if len(params.Items) == 0 {
		params.IsEmpty = true
	}
	return params
}

func (r *Renderer) retrieveTableOfContentItem(childBlock *childBlock, tableOfContentItems []templ.Component) []templ.Component {
	if childBlock.GetText() != nil {
		style := childBlock.GetText().GetStyle()
		name := r.getHeadingName(childBlock.Block)
		switch style {
		case model.BlockContentText_Header1:
			tableOfContentItems = append(tableOfContentItems, HeadingTemplate(childBlock.GetId(), name, 1))
		case model.BlockContentText_Header2:
			tableOfContentItems = r.processSecondHeading(childBlock, tableOfContentItems, name)
		case model.BlockContentText_Header3:
			tableOfContentItems = r.processThirdHeading(childBlock, tableOfContentItems, name)
		}
	}
	return tableOfContentItems
}

func (r *Renderer) processThirdHeading(childBlock *childBlock, tableOfContentItems []templ.Component, name string) []templ.Component {
	if childBlock.isChild {
		tableOfContentItems = append(tableOfContentItems, HeadingTemplate(childBlock.GetId(), name, 3))
	} else {
		tableOfContentItems = append(tableOfContentItems, HeadingTemplate(childBlock.GetId(), name, 2))
	}
	return tableOfContentItems
}

func (r *Renderer) processSecondHeading(childBlock *childBlock, tableOfContentItems []templ.Component, name string) []templ.Component {
	if childBlock.isChild {
		tableOfContentItems = append(tableOfContentItems, HeadingTemplate(childBlock.GetId(), name, 2))
	} else {
		tableOfContentItems = append(tableOfContentItems, HeadingTemplate(childBlock.GetId(), name, 1))
	}
	return tableOfContentItems
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
