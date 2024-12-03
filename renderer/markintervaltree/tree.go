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

func marksOverlap(a, b *model.Range) bool {
	return (a.From < b.To && a.To > b.From)
}

func rangeEqual(a, b *model.Range) bool {
	return (a.From == b.From && a.To == b.To)
}

func SearchOverlaps(n *MarkIntervalTreeNode, i *model.Range, result *[]*model.BlockContentTextMark) {
	if n == nil || n.MaxUpperVal < i.From {
		return
	}

	SearchOverlaps(n.Left, i, result)

	if marksOverlap(n.Mark.Range, i) {
		*result = append(*result, n.Mark)
	}
	if i.From < n.Mark.Range.From {
		return
	}

	SearchOverlaps(n.Right, i, result)
}
