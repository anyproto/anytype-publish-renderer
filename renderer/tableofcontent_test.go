package renderer

import (
	"context"
	"strings"
	"testing"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/stretchr/testify/assert"
)

func TestMakeTableOfContentRenderParams(t *testing.T) {
	renderer := &Renderer{}

	t.Run("empty Block", func(t *testing.T) {
		// given
		block := &model.Block{Id: "block1", Content: &model.BlockContentOfTableOfContents{}}

		// when
		params := renderer.makeTableOfContentBlockParams(block)

		// then
		assert.Equal(t, "block1", params.Id)
		assert.Equal(t, []string{"block", "align0", "blockTableOfContents"}, params.Classes)
		builder := strings.Builder{}
		err := params.Content.Render(context.Background(), &builder)
		assert.NoError(t, err)
		assert.Equal(t, `<div class="wrap"></div>`, builder.String())
	})

	t.Run("with background color", func(t *testing.T) {
		// given
		block := &model.Block{Id: "block1", BackgroundColor: "red", Content: &model.BlockContentOfTableOfContents{}}

		// when
		params := renderer.makeTableOfContentBlockParams(block)

		// then
		assert.Equal(t, []string{"block", "align0", "blockTableOfContents"}, params.Classes)
		assert.Equal(t, []string{"bgColor bgColor-red"}, params.ContentClasses)
	})

	t.Run("with headers", func(t *testing.T) {
		// given
		renderer.BlocksById = map[string]*model.Block{
			"root":   {Id: "root", ChildrenIds: []string{"child1", "child2"}, Content: &model.BlockContentOfSmartblock{Smartblock: &model.BlockContentSmartblock{}}},
			"child1": {Id: "child1", Content: &model.BlockContentOfText{Text: &model.BlockContentText{Style: model.BlockContentText_Header1, Text: "Header 1"}}},
			"child2": {Id: "child2", Content: &model.BlockContentOfText{Text: &model.BlockContentText{Style: model.BlockContentText_Header3, Text: "Header 3"}}},
		}

		renderer.Sp = &pb.SnapshotWithType{
			Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
				Blocks: []*model.Block{
					{Id: "root", ChildrenIds: []string{"child1", "child2"}},
				},
			},
			},
		}
		renderer.Root = &model.Block{Id: "root", ChildrenIds: []string{"child1", "child2"}, Content: &model.BlockContentOfSmartblock{Smartblock: &model.BlockContentSmartblock{}}}

		// when
		block := &model.Block{Id: "block1", Content: &model.BlockContentOfTableOfContents{}}
		params := renderer.makeTableOfContentBlockParams(block)

		// then
		assert.NotEmpty(t, params.Content)
	})
}
