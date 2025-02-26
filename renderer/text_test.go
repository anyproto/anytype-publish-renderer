package renderer

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"github.com/anyproto/anytype-publish-renderer/utils/tests/htmltag"
	"github.com/gogo/protobuf/types"

	"github.com/stretchr/testify/assert"
)

type pathAssertion struct {
	path          string
	expectedValue string
}

func TestMakeRenderText(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		r := getTestRenderer("Anytype.WebPublish.20241217.112212.67")
		id := "66c58b2a7e4bcd764b24c205"
		textBlock := r.BlocksById[id]

		expected := &BlockParams{
			Id:          id,
			Classes:     []string{"block", "align0", "blockText", "textParagraph"},
			ChildrenIds: nil,
		}

		actual := r.makeTextBlockParams(textBlock)
		assert.Equal(t, expected.Id, actual.Id)
		assert.Equal(t, expected.Classes, actual.Classes)
		assert.EqualValues(t, expected.ChildrenIds, actual.ChildrenIds)
	})
	t.Run("anytype object link in markdown", func(t *testing.T) {
		// given
		r := Renderer{}
		expected := &BlockParams{
			Classes:     []string{"block", "align0", "blockText", "textParagraph"},
			ChildrenIds: nil,
		}
		pbFiles := map[string]*pb.SnapshotWithType{
			filepath.Join("objects", "anytypeId.pb"): {
				SbType: model.SmartBlockType_Page,
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():      pbtypes.String("anytypeId"),
						bundle.RelationKeySpaceId.String(): pbtypes.String("spaceId"),
					}},
				}},
			},
		}
		r.CachedPbFiles = pbFiles

		// when
		actual := r.makeTextBlockParams(&model.Block{Content: &model.BlockContentOfText{Text: &model.BlockContentText{
			Text:  "test",
			Style: 0,
			Marks: &model.BlockContentTextMarks{
				Marks: []*model.BlockContentTextMark{
					{
						Range: &model.Range{
							From: 0,
							To:   4,
						},
						Type:  model.BlockContentTextMark_Object,
						Param: "anytypeId",
					},
				},
			},
		}}})

		// then
		assert.Equal(t, expected.Id, actual.Id)
		assert.Equal(t, expected.Classes, actual.Classes)
		assert.NotNil(t, actual.Content)
		builder := strings.Builder{}
		err := actual.Content.Render(context.Background(), &builder)
		assert.NoError(t, err)
		expectedHtml := `<div class="flex"><div class="text"><a href="anytype://object?objectId=anytypeId&spaceId=spaceId" class="markuplink" target="_blank">test</a></div></div>`
		assert.Equal(t, expectedHtml, builder.String())
	})
	t.Run("object is missing", func(t *testing.T) {
		// given
		r := Renderer{CachedPbFiles: make(map[string]*pb.SnapshotWithType), UberSp: &PublishingUberSnapshot{PbFiles: make(map[string]string)}}
		expected := &BlockParams{
			Classes:     []string{"block", "align0", "blockText", "textParagraph"},
			ChildrenIds: nil,
		}

		// when
		actual := r.makeTextBlockParams(&model.Block{Content: &model.BlockContentOfText{Text: &model.BlockContentText{
			Text:  "test",
			Style: 0,
			Marks: &model.BlockContentTextMarks{
				Marks: []*model.BlockContentTextMark{
					{
						Range: &model.Range{
							From: 0,
							To:   4,
						},
						Type:  model.BlockContentTextMark_Object,
						Param: "anytypeId",
					},
				},
			},
		}}})

		// then
		assert.Equal(t, expected.Id, actual.Id)
		assert.Equal(t, expected.Classes, actual.Classes)
		assert.NotNil(t, actual.Content)
		builder := strings.Builder{}
		err := actual.Content.Render(context.Background(), &builder)
		assert.NoError(t, err)
		expectedHtml := `<div class="flex"><div class="text"><markupobject>test</markupobject></div></div>`
		assert.Equal(t, expectedHtml, builder.String())
	})
	t.Run("anytype object mention in markdown", func(t *testing.T) {
		// given
		r := Renderer{}
		expected := &BlockParams{
			Classes:     []string{"block", "align0", "blockText", "textParagraph"},
			ChildrenIds: nil,
		}
		pbFiles := map[string]*pb.SnapshotWithType{
			filepath.Join("objects", "anytypeId.pb"): {
				SbType: model.SmartBlockType_Page,
				Snapshot: &pb.ChangeSnapshot{Data: &model.SmartBlockSnapshotBase{
					Details: &types.Struct{Fields: map[string]*types.Value{
						bundle.RelationKeyId.String():      pbtypes.String("anytypeId"),
						bundle.RelationKeySpaceId.String(): pbtypes.String("spaceId"),
					}},
				}},
			},
		}
		r.CachedPbFiles = pbFiles

		// when
		actual := r.makeTextBlockParams(&model.Block{Content: &model.BlockContentOfText{Text: &model.BlockContentText{
			Text:  "test",
			Style: 0,
			Marks: &model.BlockContentTextMarks{
				Marks: []*model.BlockContentTextMark{
					{
						Range: &model.Range{
							From: 0,
							To:   4,
						},
						Type:  model.BlockContentTextMark_Mention,
						Param: "anytypeId",
					},
				},
			},
		}}})

		// then
		assert.Equal(t, expected.Id, actual.Id)
		assert.Equal(t, expected.Classes, actual.Classes)
		assert.NotNil(t, actual.Content, 1)

		// expectedHtml := `<div class="flex"><div class="text"><a href=anytype://object?objectId=anytypeId&spaceId=spaceId target="_blank" class="markupmention withImage"><span class="smile"><div class="iconObject withDefault c20"><img src="/img/icon/default/page.svg" class="iconCommon c18"></div></span><img src="./static/img/space.svg" class="space" /><span class="name">test</span></a></div></div>`
		pathAssertions := []pathAssertion{
			{"div.flex > div.text > a.markupmention.withImage > attrs[href]", "anytype://object?objectId=anytypeId&spaceId=spaceId"},
			{"div.flex > div.text > a.markupmention.withImage > span.smile > div.iconObject.withDefault.c20 > img.iconCommon > attrs[src]", "/img/icon/default/page.svg"},
			{"div.flex > div.text > a.markupmention.withImage > img.space > attrs[src]", "./static/img/space.svg"},
		}
		assertHtmlTag(t, actual, pathAssertions)
	})
}

func assertHtmlTag(t *testing.T, actual *BlockParams, pathAssertions []pathAssertion) {
	builder := strings.Builder{}
	err := actual.Content.Render(context.Background(), &builder)
	assert.NoError(t, err)

	got, err := htmltag.HtmlToTag(builder.String())
	assert.NoError(t, err)

	for _, assertion := range pathAssertions {
		htmltag.AssertPath(t, got, assertion.path, assertion.expectedValue)
	}

}
