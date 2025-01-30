package renderer

import (
	"fmt"
	"github.com/gogo/protobuf/types"
	"strconv"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

type FeaturedRelationsParams struct {
	Id    string
	Cells []templ.Component
}

func (r *Renderer) MakeFeaturedRelationsParams(block *model.Block) *FeaturedRelationsParams {
	id := block.GetId()
	details := r.Sp.GetSnapshot().GetData().GetDetails()
	if details == nil || len(details.GetFields()) == 0 {
		return &FeaturedRelationsParams{Id: id}
	}
	featuredRelationsList := details.GetFields()[bundle.RelationKeyFeaturedRelations.String()].GetListValue()
	if featuredRelationsList == nil {
		return &FeaturedRelationsParams{Id: id}
	}
	cells := make([]templ.Component, 0, len(featuredRelationsList.Values))
	for i, featuredRelation := range featuredRelationsList.Values {
		var lastClass string
		if i == len(featuredRelationsList.Values)-1 {
			lastClass = "last"
		}
		cells = r.processFeatureRelation(featuredRelation, lastClass, details, cells)
	}
	return &FeaturedRelationsParams{Id: id, Cells: cells}
}

func (r *Renderer) processFeatureRelation(featuredRelation *types.Value, lastClass string, details *types.Struct, cells []templ.Component) []templ.Component {
	if featuredRelation == nil {
		return cells
	}
	name, format, found := r.retrieveRelationInfo(featuredRelation.GetStringValue())
	if !found {
		return cells
	}
	relationValue := details.GetFields()[featuredRelation.GetStringValue()]
	formatClass := r.getFormatClass(format)
	if relationValue == nil {
		return append(cells, EmptyCellTemplate(name, formatClass, lastClass))
	}
	if formatClass == "c-object" || formatClass == "c-file" || formatClass == "c-select" {
		cells = r.processObjectList(relationValue, format, cells, name, formatClass, lastClass)
	} else {
		cells = r.processOneObject(relationValue, format, cells, name, lastClass, formatClass)
	}
	return cells
}

func (r *Renderer) processObjectList(relationValue *types.Value, format model.RelationFormat, cells []templ.Component, name, formatClass, lastClass string) []templ.Component {
	objectsList := r.populateRelationListValue(format, relationValue)
	if len(objectsList) == 0 {
		return append(cells, EmptyCellTemplate(name, formatClass, lastClass))
	}
	var more string
	if len(objectsList) > 1 {
		more = fmt.Sprintf("+%s", strconv.FormatInt(int64(len(objectsList)-1), 10))
	}
	cells = append(cells, ListCellTemplate(formatClass, lastClass, more, objectsList[0]))
	return cells
}

func (r *Renderer) processOneObject(relationValue *types.Value, format model.RelationFormat, cells []templ.Component, name, lastClass, formatClass string) []templ.Component {
	cell := r.populateRelationValue(format, relationValue)
	if cell != nil {
		if format == model.RelationFormat_checkbox {
			cells = append(cells, CheckBoxCellTemplate(name, lastClass, cell))
		} else {
			cells = append(cells, CellTemplate(formatClass, lastClass, cell))
		}
	} else {
		cells = append(cells, EmptyCellTemplate(name, formatClass, lastClass))
	}
	return cells
}

func (r *Renderer) RenderFeaturedRelations(block *model.Block) templ.Component {
	params := r.MakeFeaturedRelationsParams(block)
	return FeaturedRelationTemplate(params)
}
