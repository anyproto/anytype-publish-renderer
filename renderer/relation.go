package renderer

import (
	"fmt"
	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/core/domain"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/gogo/protobuf/types"
	"path/filepath"
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
	IsSystem        bool
}

func (r *Renderer) MakeRelationRenderParams(b *model.Block) *RelationRenderParams {
	relationBlock := b.GetRelation()
	color := b.GetBackgroundColor()

	params := &RelationRenderParams{
		Id: b.Id,
	}

	if color != "" {
		params.BackgroundColor = fmt.Sprintf("bgColor bgColor-%s", color)
	}

	relationKey := domain.RelationKey(relationBlock.GetKey())
	relation, _ := bundle.GetRelation(relationKey)

	name, format, isSystem, found := r.fetchRelationMetadata(relation, relationKey)
	if name == "" {
		name = defaultName
	}

	params.Name = name
	params.IsSystem = isSystem

	if !found {
		params.IsDeleted = "isDeleted"
		params.IsEmpty = "isEmpty"
		return params
	}

	relationValue := r.Sp.GetSnapshot().GetData().GetDetails().GetFields()[relationBlock.GetKey()]
	if relationValue == nil {
		params.IsEmpty = "isEmpty"
		return params
	}

	r.populateRelationValue(params, format, relationValue)
	return params
}

func (r *Renderer) fetchRelationMetadata(relation *model.Relation, relationKey domain.RelationKey) (string, model.RelationFormat, bool, bool) {
	if relation != nil {
		return relation.Name, relation.Format, true, true
	}

	for _, snapshot := range r.UberSp.PbFiles {
		sn, err := readJsonpbSnapshot(snapshot)
		if err != nil || sn.SbType != model.SmartBlockType_STRelation {
			continue
		}

		fields := sn.GetSnapshot().GetData().GetDetails().GetFields()
		if uniqueKey := fields[bundle.RelationKeyUniqueKey.String()]; uniqueKey != nil && uniqueKey.GetStringValue() == relationKey.URL() {
			name := fields[bundle.RelationKeyName.String()].GetStringValue()
			format := model.RelationFormat(int32(fields[bundle.RelationKeyRelationFormat.String()].GetNumberValue()))
			return name, format, false, true
		}
	}
	return "", model.RelationFormat_longtext, false, false
}

func (r *Renderer) populateRelationValue(params *RelationRenderParams, format model.RelationFormat, relationValue *types.Value) {
	switch format {
	case model.RelationFormat_shorttext, model.RelationFormat_longtext:
		params.Format = r.getFormatClass(format)
		params.Value = BasicTemplate(params, relationValue.GetStringValue())

	case model.RelationFormat_number:
		params.Format = "c-number"
		params.Value = BasicTemplate(params, fmt.Sprintf("%f", relationValue.GetNumberValue()))

	case model.RelationFormat_status, model.RelationFormat_tag:
		params.Format = "c-select"
		params.Value = r.generateSelectOptions(params, format, relationValue)

	case model.RelationFormat_object:
		params.Format = "c-object"
		params.Value = r.generateObjectLinks(params, relationValue)

	case model.RelationFormat_file:
		params.Format = "c-file"
		params.Value = r.generateFileIcons(params, relationValue)

	case model.RelationFormat_phone, model.RelationFormat_email, model.RelationFormat_url:
		params.Format = r.getFormatClass(format)
		params.Value = BasicTemplate(params, relationValue.GetStringValue())

	case model.RelationFormat_date:
		params.Format = "c-date"
		params.Value = BasicTemplate(params, r.formatDate(relationValue.GetNumberValue()))

	case model.RelationFormat_checkbox:
		params.Format = "c-checkbox"
		params.Value = r.generateCheckbox(params, relationValue.GetBoolValue())

	default:
		params.Format = "c-longText"
		params.Value = BasicTemplate(params, relationValue.GetStringValue())
	}
}

func (r *Renderer) getFormatClass(format model.RelationFormat) string {
	switch format {
	case model.RelationFormat_shorttext:
		return "c-shortText"
	case model.RelationFormat_longtext:
		return "c-longText"
	case model.RelationFormat_phone:
		return "c-phone"
	case model.RelationFormat_email:
		return "c-email"
	case model.RelationFormat_url:
		return "c-url"
	default:
		return ""
	}
}

func (r *Renderer) formatDate(timestamp float64) string {
	if timestamp == 0 {
		return ""
	}
	return time.Unix(int64(timestamp), 0).Format("02 Jan 2006")
}

func (r *Renderer) generateCheckbox(params *RelationRenderParams, checked bool) templ.Component {
	if checked {
		return ActiveCheckBoxTemplate(params)
	}
	return DisabledCheckBoxTemplate(params)
}

func (r *Renderer) generateSelectOptions(params *RelationRenderParams, format model.RelationFormat, relationValue *types.Value) templ.Component {
	var elements []templ.Component
	relationType := "isSelect"
	if format == model.RelationFormat_tag {
		relationType = "isMultiSelect"
	}

	for _, value := range r.extractRelationValues(relationValue) {
		tag, err := r.ReadJsonpbSnapshot(filepath.Join("relationsOptions", value.GetStringValue()+".pb"))
		if err != nil {
			continue
		}

		fields := tag.GetSnapshot().GetData().GetDetails().GetFields()
		name := fields[bundle.RelationKeyName.String()].GetStringValue()
		color := fields[bundle.RelationKeyRelationOptionColor.String()].GetStringValue()
		elements = append(elements, OptionElement(name, color, relationType))
	}

	return ListTemplate(params, elements)
}

func (r *Renderer) extractRelationValues(relationValue *types.Value) []*types.Value {
	if relationValue.GetListValue() != nil {
		return relationValue.GetListValue().Values
	}
	return []*types.Value{relationValue}
}

func (r *Renderer) generateObjectLinks(params *RelationRenderParams, relationValue *types.Value) templ.Component {
	var elements []templ.Component
	spaceId := r.Sp.GetSnapshot().GetData().GetDetails().GetFields()[bundle.RelationKeySpaceId.String()].GetStringValue()

	for _, value := range r.extractRelationValues(relationValue) {
		link := fmt.Sprintf(linkTemplate, value.GetStringValue(), spaceId)
		elements = append(elements, BasicListElement(link))
	}
	return ListTemplate(params, elements)
}

func (r *Renderer) generateFileIcons(params *RelationRenderParams, relationValue *types.Value) templ.Component {
	var elements []templ.Component

	for _, value := range r.extractRelationValues(relationValue) {
		url, err := r.getFileUrl(value.GetStringValue())
		if err != nil {
			continue
		}
		fileBlock, err := r.getFileBlock(value.GetStringValue())
		if err != nil {
			continue
		}

		elements = append(elements, r.createFileIcon(url, fileBlock))
	}
	return ListTemplate(params, elements)
}

func (r *Renderer) createFileIcon(url string, fileBlock *model.BlockContentFile) templ.Component {
	filename := filepath.Base(url)

	switch fileBlock.GetType() {
	case model.BlockContentFile_Audio:
		return AudioIconTemplate(filename)
	case model.BlockContentFile_Image:
		return ImageIconTemplate(url, filename)
	case model.BlockContentFile_Video:
		return VideoIconTemplate(url, filename)
	default:
		return FileIconTemplate(filename)
	}
}

func (r *Renderer) RenderRelations(b *model.Block) templ.Component {
	params := r.MakeRelationRenderParams(b)
	return RelationTemplate(params)
}
