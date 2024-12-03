package markintervaltree

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

func TestMarkIntervalTree(t *testing.T) {
	ranges := []*model.Range{
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
	root := &MarkIntervalTreeNode{
		Mark: &model.BlockContentTextMark{
			Range: ranges[0],
		},
	}

	for i := 1; i < len(ranges); i++ {
		root.Insert(&model.BlockContentTextMark{
			Range: ranges[i],
		})
	}

	t.Run("simple test", func(t *testing.T) {
		results := make([]*model.BlockContentTextMark, 0)
		SearchOverlaps(root, &model.BlockContentTextMark{
			Range: &model.Range{
				From: 17,
				To:   19,
			},
		}, &results)

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

}
