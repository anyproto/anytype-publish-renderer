package markintervaltree

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

func TestMarkIntervalTree(t *testing.T) {
	ranges := []*model.Range{
		&model.Range{
			From: 5,
			To:   6,
		},
		&model.Range{
			From: 27,
			To:   37,
		},
		&model.Range{
			From: 17,
			To:   30,
		},
		&model.Range{
			From: 20,
			To:   23,
		},
		&model.Range{
			From: 14,
			To:   21,
		},
		&model.Range{
			From: 21,
			To:   39,
		},
		&model.Range{
			From: 18,
			To:   31,
		},
	}

	marks := make([]*model.BlockContentTextMark, 0)
	for _, r := range ranges {
		mark := &model.BlockContentTextMark{
			Range: r,
		}
		marks = append(marks, mark)
	}

	root := New(marks)

	t.Run("simple test", func(t *testing.T) {
		results := root.SearchOverlaps(&model.Range{
			From: 17,
			To:   19,
		})

		expected := []*model.BlockContentTextMark{
			&model.BlockContentTextMark{
				Range: &model.Range{
					From: 14,
					To:   21,
				},
			},

			&model.BlockContentTextMark{
				Range: &model.Range{
					From: 17,
					To:   30,
				},
			},
			&model.BlockContentTextMark{
				Range: &model.Range{
					From: 18,
					To:   31,
				},
			},
		}
		assert.EqualValues(t, expected, results)
	})

	t.Run("single range", func(t *testing.T) {
		results := root.SearchOverlaps(&model.Range{
			From: 5,
			To:   6,
		})

		expected := []*model.BlockContentTextMark{
			&model.BlockContentTextMark{
				Range: &model.Range{
					From: 5,
					To:   6,
				},
			},
		}
		assert.EqualValues(t, expected, results)
	})

	t.Run("one item", func(t *testing.T) {
		singleRoot := &MarkIntervalTreeNode{
			Mark: &model.BlockContentTextMark{
				Range: ranges[0],
			},
			MaxUpperVal: ranges[0].To,
		}

		results := singleRoot.SearchOverlaps(&model.Range{
			From: 5,
			To:   6,
		})

		expected := []*model.BlockContentTextMark{
			&model.BlockContentTextMark{
				Range: &model.Range{
					From: 5,
					To:   6,
				},
			},
		}
		assert.EqualValues(t, expected, results)
	})

}
