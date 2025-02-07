package renderer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/types"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
)

type FeaturedRelationsParams struct {
	Id             string
	Classes        string
	ContentClasses string
	Cells          []templ.Component
}

func (r *Renderer) MakeFeaturedRelationsParams(b *model.Block) *FeaturedRelationsParams {
	id := b.GetId()
	details := r.Sp.GetSnapshot().GetData().GetDetails()
	align := "align" + strconv.Itoa(int(b.GetAlign()))
	bgColor := b.GetBackgroundColor()
	classes := []string{"block", "blockFeatured", align}
	contentClasses := []string{"content"}

	if bgColor != "" {
		contentClasses = append(contentClasses, "bgColor", "bgColor-"+bgColor)
	}

	param := &FeaturedRelationsParams{
		Id:             id,
		Classes:        strings.Join(classes, " "),
		ContentClasses: strings.Join(contentClasses, " "),
	}

	if details == nil || len(details.GetFields()) == 0 {
		return param
	}

	featuredRelationsList := details.GetFields()[bundle.RelationKeyFeaturedRelations.String()].GetListValue()
	if featuredRelationsList == nil {
		return param
	}

	cells := make([]templ.Component, 0, len(featuredRelationsList.Values))
	for i, featuredRelation := range featuredRelationsList.Values {
		var lastClass string

		if i == len(featuredRelationsList.Values)-1 {
			lastClass = "last"
		}

		cells = r.processFeatureRelation(featuredRelation, lastClass, details, cells)
	}

	param.Cells = cells
	return param
}

func (r *Renderer) processFeatureRelation(featuredRelation *types.Value, lastClass string, details *types.Struct, cells []templ.Component) []templ.Component {
	if featuredRelation == nil {
		return cells
	}

	relationKey := featuredRelation.GetStringValue()

	if relationKey == bundle.RelationKeyDescription.String() {
		return cells
	}

	name, format, found := r.retrieveRelationInfo(relationKey)

	if !found {
		return cells
	}
	relationValue, exists := details.GetFields()[relationKey]
	formatClass := r.getFormatClass(format)

	if !exists || relationValue == nil {
		return append(cells, EmptyCellTemplate(name, formatClass, lastClass))
	}

	switch format {
	case
		model.RelationFormat_object,
		model.RelationFormat_file,
		model.RelationFormat_tag,
		model.RelationFormat_status:
		return r.processObjectList(relationKey, relationValue, format, cells, name, formatClass, lastClass)
	default:
		return r.processOneObject(relationValue, format, cells, name, lastClass, formatClass)
	}
}

func (r *Renderer) processObjectList(key string, relationValue *types.Value, format model.RelationFormat, cells []templ.Component, name, formatClass, lastClass string) []templ.Component {
	objectsList := r.populateRelationListValue(format, relationValue)

	if len(objectsList) == 0 {
		if key == bundle.RelationKeyBacklinks.String() || key == bundle.RelationKeyLinks.String() {
			return cells
		}
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
