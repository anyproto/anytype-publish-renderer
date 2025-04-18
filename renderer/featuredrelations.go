package renderer

import (
	"fmt"

	"github.com/anyproto/anytype-heart/util/pbtypes"
	"github.com/gogo/protobuf/types"
	"github.com/ipfs/go-cid"

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
	featuredRelationsList := getRelationField(details, bundle.RelationKeyFeaturedRelations, relationToList)
	if r.isFeaturedRelationsEmpty(featuredRelationsList) {
		featuredRelationsList = getRelationField(r.ObjectTypeDetails, bundle.RelationKeyRecommendedFeaturedRelations, relationToList)
		if featuredRelationsList == nil || len(featuredRelationsList.GetValues()) == 0 {
			return nil
		}
	}
	return featuredRelationsList
}

func (r *Renderer) isFeaturedRelationsEmpty(featuredRelationsList *types.ListValue) bool {
	return featuredRelationsList == nil || len(featuredRelationsList.GetValues()) == 0 || (len(featuredRelationsList.GetValues()) == 1 &&
		featuredRelationsList.GetValues()[0].GetStringValue() == bundle.RelationKeyDescription.String())
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
	}
	_, err := cid.Decode(relationKey)
	if err != nil {
		settings.Key = relationKey
	} else {
		settings.Id = relationKey
	}
	cells = append(cells, r.buildRelationComponents(settings)...)
	return cells
}

func (r *Renderer) RenderFeaturedRelations(block *model.Block) templ.Component {
	blockParams := r.makeFeaturedRelationsBlockParams(block)
	if blockParams == nil {
		return NoneTemplate("")
	}
	return BlockTemplate(r, blockParams)
}

func (r *Renderer) makeFeaturedRelationsBlockParams(block *model.Block) *BlockParams {
	block.Align = model.BlockAlign(r.LayoutAlign)
	blockParams := makeDefaultBlockParams(block)
	color := block.GetBackgroundColor()
	if color != "" {
		blockParams.ContentClasses = append(blockParams.ContentClasses, fmt.Sprintf("bgColor bgColor-%s", color))
	}
	params := r.makeFeaturedRelationsComponent()
	if params == nil {
		return nil
	}
	blockParams.Content = params
	return blockParams
}
