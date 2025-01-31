package renderer

import (
	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeTableOfContentRenderParams(t *testing.T) {
	renderer := *getTestRenderer("Anytype.WebPublish.20241217.112212.67")

	t.Run("empty Block", func(t *testing.T) {
		// given
		block := &model.Block{Id: "block1"}

		// when
		params := renderer.MakeTableOfContentRenderParams(block)

		// then
		assert.Equal(t, "block1", params.Id)
		assert.Equal(t, "", params.BackgroundColor)
		assert.Empty(t, params.Items)
		assert.True(t, params.IsEmpty)
	})

	t.Run("with background color", func(t *testing.T) {
		// given
		block := &model.Block{Id: "block1", BackgroundColor: "red"}

		// when
		params := renderer.MakeTableOfContentRenderParams(block)

		// then
		assert.Equal(t, "bgColor bgColor-red", params.BackgroundColor)
	})

	t.Run("with headers", func(t *testing.T) {
		// given
		renderer.BlocksById = map[string]*model.Block{
			"child1": {Id: "child1", Content: &model.BlockContentOfText{Text: &model.BlockContentText{Style: model.BlockContentText_Header1, Text: "Header 1"}}},
			"child2": {Id: "child2", Content: &model.BlockContentOfText{Text: &model.BlockContentText{Style: model.BlockContentText_Header3, Text: "Header 3"}}},
		}

		renderer.Sp = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Blocks: []*model.Block{
					{ChildrenIds: []string{"child1", "child2"}},
				},
			},
			},
		}

		// when
		block := &model.Block{Id: "block1"}
		params := renderer.MakeTableOfContentRenderParams(block)

		// then
		assert.False(t, params.IsEmpty)
		assert.Len(t, params.Items, 2)
	})
}
