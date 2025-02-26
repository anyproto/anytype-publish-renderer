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
			Classes: []string{"block", "align0", "blockText", "textParagraph"},
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

		tag, err := blockParamsToHtmlTag(actual)
		assert.NoError(t, err)
		pathAssertions := []pathAssertion{
			{"div.flex > div.text > a.markuplink > attrs[href]", "anytype://object?objectId=anytypeId&spaceId=spaceId"},
			{"div.flex > div.text > a.markuplink > Content", "test"},
		}

		assertHtmlTag(t, tag, pathAssertions)

	})
	t.Run("object is missing", func(t *testing.T) {
		// given
		r := Renderer{CachedPbFiles: make(map[string]*pb.SnapshotWithType), UberSp: &PublishingUberSnapshot{PbFiles: make(map[string]string)}}
		expected := &BlockParams{
			Classes: []string{"block", "align0", "blockText", "textParagraph"},
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

		tag, err := blockParamsToHtmlTag(actual)
		assert.NoError(t, err)

		pathAssertions := []pathAssertion{
			{"div.flex > div.text > markupobject > Content", "test"},
		}

		assertHtmlTag(t, tag, pathAssertions)

	})
	t.Run("anytype object mention in markdown", func(t *testing.T) {
		// given
		r := Renderer{}
		expected := &BlockParams{
			Classes: []string{"block", "align0", "blockText", "textParagraph"},
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

		tag, err := blockParamsToHtmlTag(actual)
		assert.NoError(t, err)

		pathAssertions := []pathAssertion{
			{"div.flex > div.text > a.markupmention.withImage > attrs[href]", "anytype://object?objectId=anytypeId&spaceId=spaceId"},
			{"div.flex > div.text > a.markupmention.withImage > span.smile > div.iconObject.withDefault.c20 > img.iconCommon > attrs[src]", "/img/icon/default/page.svg"},
			{"div.flex > div.text > a.markupmention.withImage > img.space > attrs[src]", "./static/img/space.svg"},
			{"div.flex > div.text > a.markupmention.withImage > span.name > Content", "test"},
		}
		assertHtmlTag(t, tag, pathAssertions)
	})
}

func blockParamsToHtmlTag(actual *BlockParams) (*htmltag.Tag, error) {
	builder := strings.Builder{}
	err := actual.Content.Render(context.Background(), &builder)
	if err != nil {
		return nil, err
	}

	tag, err := htmltag.HtmlToTag(builder.String())
	if err != nil {
		return nil, err
	}

	return tag, nil

}
func assertHtmlTag(t *testing.T, tag *htmltag.Tag, pathAssertions []pathAssertion) {
	for _, assertion := range pathAssertions {
		htmltag.AssertPath(t, tag, assertion.path, assertion.expectedValue)
	}

}
