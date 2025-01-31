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

func (r *Renderer) MakeTableOfContentRenderParams(block *model.Block) *TableOfContentRenderParams {
	blockId := block.GetId()
	params := &TableOfContentRenderParams{
		Id: blockId,
	}
	color := block.GetBackgroundColor()
	if color != "" {
		params.BackgroundColor = fmt.Sprintf("bgColor bgColor-%s", color)
	}
	var tableOfContentItems []templ.Component
	for _, b := range r.Sp.GetSnapshot().GetData().GetBlocks() {
		for _, id := range b.ChildrenIds {
			if childBlock, ok := r.BlocksById[id]; ok {
				tableOfContentItems = r.retrieveTableOfContentItem(childBlock, tableOfContentItems)
			}
		}
	}
	params.Items = tableOfContentItems
	if len(params.Items) == 0 {
		params.IsEmpty = true
	}
	return params
}

func (r *Renderer) retrieveTableOfContentItem(childBlock *model.Block, tableOfContentItems []templ.Component) []templ.Component {
	if childBlock.GetText() != nil {
		style := childBlock.GetText().GetStyle()
		name := r.getHeadingName(childBlock)
		switch style {
		case model.BlockContentText_Header1, model.BlockContentText_Header2:
			tableOfContentItems = append(tableOfContentItems, FirstHeadingTemplate(name))
		case model.BlockContentText_Header3:
			tableOfContentItems = append(tableOfContentItems, SecondHeadingsTemplate(name))
		}
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
