package renderer

import (
	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

type IconImageRenderParams struct {
	Id          string
	Src         string
	Classes     string
	IconClasses string
}

func isTodoLayout(layout model.ObjectTypeLayout) bool {
	return layout == model.ObjectType_todo
}

func isHumanLayout(layout model.ObjectTypeLayout) bool {
	return layout == model.ObjectType_profile || layout == model.ObjectType_participant
}

func pageIconInitSize(layout model.ObjectTypeLayout) int32 {
	if isHumanLayout(layout) {
		return 128
	} else {
		return 96
	}
}

func (r *Renderer) RenderPageIconImage() templ.Component {
	details := r.Sp.Snapshot.Data.GetDetails()
	layout := getRelationField(details, bundle.RelationKeyLayout, relationToObjectTypeLayout)
	iconEmoji := getRelationField(details, bundle.RelationKeyIconEmoji, r.relationToEmojiUrl)
	iconImage := getRelationField(details, bundle.RelationKeyIconImage, r.relationToFileUrl)

	if isTodoLayout(layout) {
		return NoneTemplate("")
	}

	if iconEmoji != "" && iconImage != "" {
		return NoneTemplate("")
	}

	params := r.MakeRenderIconObjectParams(details, &IconObjectProps{
		NoDefault: true,
		Size:      pageIconInitSize(layout),
	})
	if params.Src == "" {
		return NoneTemplate("")
	}

	content := IconObjectTemplate(r, params)

	classes := []string{}
	if isHumanLayout(layout) {
		classes = append(classes, "isHuman")
	}

	blockParams := &BlockParams{
		BlockType: "Icon",
		Classes:   classes,
		Content:   content,
	}
	return BlockTemplate(r, blockParams)
}
