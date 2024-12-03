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

func marksOverlap(a, b *model.BlockContentTextMark) bool {
	return (a.Range.From <= b.Range.To && a.Range.To >= b.Range.From)
}

func SearchOverlaps(n *MarkIntervalTreeNode, m *model.BlockContentTextMark, result *[]*model.BlockContentTextMark) {
	if n == nil || n.MaxUpperVal < m.Range.From {
		return
	}

	SearchOverlaps(n.Left, m, result)

	if marksOverlap(n.Mark, m) {
		if n.Mark.Range.From != m.Range.From || n.Mark.Range.To != m.Range.To {
			*result = append(*result, n.Mark)
		}
	}
	if m.Range.From < n.Mark.Range.From {
		return
	}

	SearchOverlaps(n.Right, m, result)
}
