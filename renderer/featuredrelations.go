package renderer

import (
	"fmt"

	"github.com/anyproto/anytype-heart/util/pbtypes"
	"github.com/gogo/protobuf/types"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

func (r *Renderer) makeFeaturedRelationsComponent() templ.Component {
	details := r.Sp.GetSnapshot().GetData().GetDetails()

	if details == nil || len(details.GetFields()) == 0 {
		return nil
	}
	featuredRelationsList := r.retrieveFeaturedRelations(details)
	if featuredRelationsList == nil {
		return nil
	}
	cells := make([]templ.Component, 0, len(featuredRelationsList.Values))
	for i, featuredRelation := range featuredRelationsList.Values {
		var lastClass string
		if i == len(featuredRelationsList.Values)-1 {
			lastClass = "last"
		}
		cells = r.processFeatureRelation(featuredRelation, details, lastClass, cells)
	}
	if len(cells) == 0 {
		return nil
	}
	wrapper := BlocksWrapper(&BlockWrapperParams{
		Classes:    []string{"wrap"},
		Components: cells,
	})
	return wrapper
}

func (r *Renderer) retrieveFeaturedRelations(details *types.Struct) *types.ListValue {
	featuredRelationsList := details.GetFields()[bundle.RelationKeyFeaturedRelations.String()].GetListValue()
	if featuredRelationsList == nil || len(featuredRelationsList.GetValues()) == 0 || (len(featuredRelationsList.GetValues()) == 1 &&
		featuredRelationsList.GetValues()[0].GetStringValue() == bundle.RelationKeyDescription.URL()) {
		featuredRelationsList = r.ObjectTypeDetails.GetFields()[bundle.RelationKeyRecommendedFeaturedRelations.String()].GetListValue()
		if featuredRelationsList == nil || len(featuredRelationsList.GetValues()) == 0 {
			return nil
		}
	}
	return featuredRelationsList
}

func (r *Renderer) processFeatureRelation(featuredRelation *types.Value, details *types.Struct, lastClass string, cells []templ.Component) []templ.Component {
	if featuredRelation == nil {
		return cells
	}
	relationKey := featuredRelation.GetStringValue()
	if relationKey == bundle.RelationKeyDescription.String() {
		return cells
	}
	if relationKey == bundle.RelationKeyBacklinks.String() || relationKey == bundle.RelationKeyLinks.String() {
		list := pbtypes.GetStringList(details, relationKey)
		if len(list) == 0 {
			return cells
		}
	}
	settings := &RelationRenderSetting{
		Featured:     true,
		LimitDisplay: true,
		Classes:      []string{lastClass},
		Key:          relationKey,
	}
	cells = append(cells, r.buildRelationComponents(settings)...)
	return cells
}

func (r *Renderer) RenderFeaturedRelations(block *model.Block) templ.Component {
	blockParams := makeDefaultBlockParams(block)
	color := block.GetBackgroundColor()
	if color != "" {
		blockParams.ContentClasses = append(blockParams.ContentClasses, fmt.Sprintf("bgColor bgColor-%s", color))
	}
	params := r.makeFeaturedRelationsComponent()
	if params == nil {
		return NoneTemplate("")
	}
	blockParams.Content = params
	return BlockTemplate(r, blockParams)
}
