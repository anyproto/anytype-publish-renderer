package renderer

import (
	"fmt"
	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/core/domain"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/gogo/protobuf/types"
	"strings"
	"time"
)

const defaultName = "Untitled"

type RelationRenderParams struct {
	Id              string
	BackgroundColor string
	IsDeleted       string
	Name            string
	IsEmpty         string
	Format          string
	Value           templ.Component
}

func (r *Renderer) MakeRelationRenderParams(b *model.Block) (params *RelationRenderParams) {
	relationBlock := b.GetRelation()
	color := b.GetBackgroundColor()
	params = &RelationRenderParams{
		Id: b.Id,
	}
	params.BackgroundColor = strings.Join([]string{"bgColor", "bgColor-" + color}, " ")
	key := relationBlock.GetKey()
	relation, _ := bundle.GetRelation(domain.RelationKey(key))
	var (
		name    string
		format  model.RelationFormat
		founded bool
	)
	if relation == nil {
		for _, sn := range r.CachedPbFiles {
			if sn.SbType != model.SmartBlockType_STRelation {
				continue
			}
			uniqueKey := sn.GetSnapshot().GetData().GetDetails().GetFields()[bundle.RelationKeyUniqueKey.String()]
			if uniqueKey != nil && uniqueKey.GetStringValue() == key {
				name = sn.GetSnapshot().GetData().GetDetails().GetFields()[bundle.RelationKeyName.String()].GetStringValue()
				format = model.RelationFormat(int32(sn.GetSnapshot().GetData().GetDetails().GetFields()[bundle.RelationKeyRelationFormat.String()].GetNumberValue()))
				founded = true
				break
			}
		}
	} else {
		name = relation.Name
		format = relation.Format
		founded = true
	}
	if name == "" {
		name = defaultName
	}
	if founded {
		params.IsDeleted = "isDeleted"
		params.IsEmpty = "isEmpty"
		return
	}
	relationValue := r.Sp.GetSnapshot().GetData().GetDetails().GetFields()[key]
	if relationValue == nil {
		params.IsEmpty = "isEmpty"
		return
	}
	switch format {
	case model.RelationFormat_shorttext:
		params.Format = "c-shortText"
		params.Value = BasicTemplate(params, relationValue.GetStringValue())
	case model.RelationFormat_longtext:
		params.Format = "c-longText"
		params.Value = BasicTemplate(params, relationValue.GetStringValue())
	case model.RelationFormat_number:
		params.Format = "c-number"
		number := fmt.Sprintf("%f", relationValue.GetNumberValue())
		params.Value = BasicTemplate(params, number)
	case model.RelationFormat_status, model.RelationFormat_tag:
		params.Format = "c-select"
		var (
			elements       []templ.Component
			relationValues = relationValue.GetListValue().Values
		)
		if len(relationValues) == 0 {
			relationValues = []*types.Value{relationValue}
		}
		for _, value := range relationValues {
			if tag, ok := r.CachedPbFiles[value.GetStringValue()]; ok {
				name := tag.GetSnapshot().GetData().GetDetails().GetFields()[bundle.RelationKeyName.String()].GetStringValue()
				elements = append(elements, ListElement(name))
			}
		}
		params.Value = ListTemplate(elements)
	case model.RelationFormat_object:
		params.Format = "c-object"
		spaceId := r.Sp.GetSnapshot().GetData().GetDetails().GetFields()[bundle.RelationKeySpaceId.String()]
		var elements []templ.Component
		for _, value := range relationValue.GetListValue().Values {
			link := fmt.Sprintf(linkTemplate, value, spaceId)
			elements = append(elements, ListElement(link))
		}
		params.Value = ListTemplate(elements)
	case model.RelationFormat_file:
		params.Format = "c-file"
		var elements []templ.Component
		for _, value := range relationValue.GetListValue().Values {
			url, err := r.getFileUrl(value.GetStringValue())
			if err != nil {
				continue
			}
			fileBlock, err := r.getFileBlock(value.GetStringValue())
			if err != nil {
				continue
			}
			switch fileBlock.GetType() {
			case model.BlockContentFile_Audio:
				elements = append(elements, AudioIconTemplate(url))
			case model.BlockContentFile_Image:
				elements = append(elements, ImageIconTemplate(url, fileBlock.GetName()))
			case model.BlockContentFile_Video:
				elements = append(elements, VideoIconTemplate(url, fileBlock.GetName()))
			default:
				elements = append(elements, FileIconTemplate(url))
			}
		}
	case model.RelationFormat_phone:
		params.Format = "c-phone"
		params.Value = BasicTemplate(params, relationValue.GetStringValue())
	case model.RelationFormat_email:
		params.Format = "c-email"
		params.Value = BasicTemplate(params, relationValue.GetStringValue())
	case model.RelationFormat_url:
		params.Format = "c-url"
		params.Value = BasicTemplate(params, relationValue.GetStringValue())
	case model.RelationFormat_date:
		params.Format = "c-date"
		date := time.Unix(0, int64(relationValue.GetNumberValue()))
		params.Value = BasicTemplate(params, date.Format("02 Jan 2006"))
	case model.RelationFormat_checkbox:
		params.Format = "c-checkbox"
		if relationValue.GetBoolValue() {
			params.Value = ActiveCheckBoxTemplate(params)
		} else {
			params.Value = DisabledCheckBoxTemplate(params)
		}
	default:
		params.Format = "c-longText"
		params.Value = BasicTemplate(params, relationValue.GetStringValue())
	}
	return

}
func (r *Renderer) RenderRelations(b *model.Block) templ.Component {
	params := r.MakeRelationRenderParams(b)
	return RelationTemplate(params)
}
