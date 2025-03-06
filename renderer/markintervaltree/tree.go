package markintervaltree

import (
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

type MarkIntervalTreeNode struct {
	Mark        *model.BlockContentTextMark
	MaxUpperVal int32
	Left        *MarkIntervalTreeNode
	Right       *MarkIntervalTreeNode
}

func New(marks []*model.BlockContentTextMark) *MarkIntervalTreeNode {
	root := &MarkIntervalTreeNode{
		Mark:        marks[0],
		MaxUpperVal: marks[0].Range.To,
	}

	for i := 1; i < len(marks); i++ {
		root.Insert(marks[i])
	}

	return root

}

func (r *MarkIntervalTreeNode) Insert(m *model.BlockContentTextMark) {
	node := r
	for node != nil {
		if node.MaxUpperVal < m.Range.To {
			node.MaxUpperVal = m.Range.To
		}
		if m.Range.From < node.Mark.Range.From {
			if node.Left == nil {
				node.Left = &MarkIntervalTreeNode{
					Mark:        m,
					MaxUpperVal: m.Range.To,
				}
				return
			} else {
				node = node.Left
				continue
			}
		} else {
			if node.Right == nil {
				node.Right = &MarkIntervalTreeNode{
					Mark:        m,
					MaxUpperVal: m.Range.To,
				}
				return
			} else {
				node = node.Right
				continue
			}

		}
	}
}

func rangeOverlap(a, b *model.Range) bool {
	return (a.From <
		b.To && a.To > b.From) || (a.From == b.From && a.To == b.To)
}

func (r *MarkIntervalTreeNode) SearchOverlaps(i *model.Range) []*model.BlockContentTextMark {
	marksToApply := make([]*model.BlockContentTextMark, 0)
	searchOverlaps(r, i, &marksToApply)
	return marksToApply
}

func searchOverlaps(n *MarkIntervalTreeNode, i *model.Range, result *[]*model.BlockContentTextMark) {
	if n == nil || n.MaxUpperVal < i.From {
		return
	}
	searchOverlaps(n.Left, i, result)

	if rangeOverlap(n.Mark.Range, i) {
		*result = append(*result, n.Mark)
	}
	if i.From < n.Mark.Range.From {
		return
	}

	searchOverlaps(n.Right, i, result)
}
