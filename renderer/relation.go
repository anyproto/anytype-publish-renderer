package renderer

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/core/domain"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/localstore/addr"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
	"github.com/gogo/protobuf/types"
)

const defaultName = "Untitled"

type RelationRenderSetting struct {
	Id           string
	Key          string
	Name         string
	Featured     bool
	LimitDisplay bool
	Classes      []string
}

func (r *Renderer) makeRelationTemplate(b *model.Block) templ.Component {
	relationBlock := b.GetRelation()
	key := relationBlock.GetKey()
	if key == "" {
		return nil
	}
	params := &RelationRenderSetting{Key: key}
	relationComponent := r.buildRelationComponents(params)
	if relationComponent == nil {
		return nil
	}
	return BlocksWrapper(&BlockWrapperParams{Classes: []string{"sides"}, Components: relationComponent})
}

func (r *Renderer) buildRelationComponents(params *RelationRenderSetting) []templ.Component {
	name, format, key, found := r.retrieveRelationInfo(params)
	if !found {
		return nil
	}
	var components []templ.Component
	if !params.Featured {
		components = append(components, BlocksWrapper(&BlockWrapperParams{
			Classes:    []string{"info"},
			Components: []templ.Component{BasicTemplate("name", name)},
		}))
	}
	relationValue := r.Sp.GetSnapshot().GetData().GetDetails().GetFields()[key]
	if relationValue == nil {
		params.Classes = append(params.Classes, "isEmpty")
		return append(components, CellTemplate(params, BasicTemplate("empty", "")))
	}
	formatClass := r.getFormatClass(format)
	params.Classes = append(params.Classes, formatClass)
	switch format {
	case model.RelationFormat_object, model.RelationFormat_tag, model.RelationFormat_status, model.RelationFormat_file:
		listTemplate := r.buildListComponent(params, format, relationValue)
		if listTemplate == nil {
			return components
		}
		components = append(components, CellTemplate(params, listTemplate))
	default:
		params.Name = name
		var component = r.populateRelationValue(params, format, relationValue)
		if component == nil {
			return components
		}
		components = append(components, CellTemplate(params, component))
	}
	return components
}

func (r *Renderer) buildListComponent(params *RelationRenderSetting, format model.RelationFormat, relationValue *types.Value) templ.Component {
	components := r.populateRelationListValue(format, relationValue)
	if len(components) == 0 {
		return nil
	}
	var listTemplate templ.Component
	if params.LimitDisplay && (format == model.RelationFormat_object || format == model.RelationFormat_file) && len(components) > 1 {
		more := fmt.Sprintf("+%s", strconv.FormatInt(int64(len(components)-1), 10))
		listTemplate = ListTemplate(more, components[0:1])
	} else {
		listTemplate = ListTemplate("", components)
	}
	return listTemplate
}

func (r *Renderer) retrieveRelationInfo(params *RelationRenderSetting) (string, model.RelationFormat, string, bool) {
	name, format, key, found := r.fetchRelationMetadata(params)
	if name == "" {
		name = defaultName
	}
	return name, format, key, found
}

func (r *Renderer) fetchRelationMetadata(params *RelationRenderSetting) (string, model.RelationFormat, string, bool) {
	if params.Id != "" {
		return r.getRelationDataById(params)
	}
	if params.Key != "" {
		return r.getRelationByKey(params.Key)
	}
	return "", 0, "", false
}

func (r *Renderer) getRelationByKey(key string) (string, model.RelationFormat, string, bool) {
	relationKey := domain.RelationKey(key)
	relation, _ := bundle.GetRelation(relationKey)
	if relation != nil {
		return relation.Name, relation.Format, key, true
	}
	for _, sn := range r.UberSp.PbFiles {
		sn, err := readJsonpbSnapshot(sn)
		if err != nil || sn.SbType != model.SmartBlockType_STRelation {
			continue
		}
		fields := sn.GetSnapshot().GetData().GetDetails().GetFields()
		if uniqueKey := fields[bundle.RelationKeyUniqueKey.String()]; uniqueKey != nil && uniqueKey.GetStringValue() == relationKey.URL() {
			name := fields[bundle.RelationKeyName.String()].GetStringValue()
			format := model.RelationFormat(int32(fields[bundle.RelationKeyRelationFormat.String()].GetNumberValue()))
			return name, format, key, true
		}
	}
	return "", 0, "", false
}

func (r *Renderer) getRelationDataById(params *RelationRenderSetting) (string, model.RelationFormat, string, bool) {
	snapshot := r.getObjectSnapshot(params.Id)
	if snapshot != nil {
		name := getRelationField(snapshot.GetSnapshot().GetData().GetDetails(), bundle.RelationKeyName, relationToString)
		format := getRelationField(snapshot.GetSnapshot().GetData().GetDetails(), bundle.RelationKeyRelationFormat, relationToRelationFormat)
		uk := getRelationField(snapshot.GetSnapshot().GetData().GetDetails(), bundle.RelationKeyUniqueKey, relationToString)
		return name, format, strings.TrimPrefix(uk, addr.RelationKeyToIdPrefix), true
	}
	return "", 0, "", false
}

func (r *Renderer) populateRelationListValue(format model.RelationFormat, relationValue *types.Value) []templ.Component {
	switch format {
	case model.RelationFormat_status, model.RelationFormat_tag:
		return r.generateSelectOptions(format, relationValue)
	case model.RelationFormat_object:
		return r.generateObjectLinks(relationValue)
	case model.RelationFormat_file:
		return r.generateFileComponent(relationValue)
	}
	return nil
}

func (r *Renderer) populateRelationValue(params *RelationRenderSetting, format model.RelationFormat, relationValue *types.Value) templ.Component {
	if format != model.RelationFormat_checkbox && pbtypes.IsEmptyValue(relationValue) {
		return nil
	}
	switch format {
	case model.RelationFormat_shorttext, model.RelationFormat_longtext:
		return BasicTemplate("name", relationValue.GetStringValue())
	case model.RelationFormat_number:
		return BasicTemplate("name", fmt.Sprintf("%g", relationValue.GetNumberValue()))
	case model.RelationFormat_phone, model.RelationFormat_email, model.RelationFormat_url:
		url := getUrlScheme(format, relationValue.GetStringValue()) + relationValue.GetStringValue()
		return ObjectElement(relationValue.GetStringValue(), templ.URL(url))
	case model.RelationFormat_date:
		return BasicTemplate("name", r.formatDate(relationValue.GetNumberValue()))
	case model.RelationFormat_checkbox:
		return r.generateCheckbox(params, relationValue.GetBoolValue())
	}
	return nil
}

func getUrlScheme(format model.RelationFormat, value string) string {
	if value == "" {
		return ""
	}
	if format == model.RelationFormat_url {
		parsedUrl, err := url.Parse(value)
		if err != nil {
			return ""
		}
		if parsedUrl.Scheme == "" {
			return "http://"
		}
	}
	if format == model.RelationFormat_email {
		return "mailto:"
	}
	if format == model.RelationFormat_phone {
		return "tel:"
	}
	return ""
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
	case model.RelationFormat_object:
		return "c-object"
	case model.RelationFormat_file:
		return "c-file"
	case model.RelationFormat_checkbox:
		return "c-checkbox"
	case model.RelationFormat_date:
		return "c-date"
	case model.RelationFormat_tag, model.RelationFormat_status:
		return "c-select"
	case model.RelationFormat_number:
		return "c-number"
	default:
		return "c-longText"
	}
}

func (r *Renderer) formatDate(timestamp float64) string {
	if timestamp == 0 {
		return ""
	}
	return time.Unix(int64(timestamp), 0).Format("02 Jan 2006")
}

func (r *Renderer) generateCheckbox(params *RelationRenderSetting, checked bool) templ.Component {
	if checked {
		return ActiveCheckBoxTemplate(params.Name, params.Featured)
	}
	return DisabledCheckBoxTemplate(params.Name, params.Featured)
}

func (r *Renderer) generateSelectOptions(format model.RelationFormat, relationValue *types.Value) []templ.Component {
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
		elements = append(elements, ListElement(OptionElement(name, color, relationType), nil))
	}

	return elements
}

func (r *Renderer) extractRelationValues(relationValue *types.Value) []*types.Value {
	if relationValue.GetListValue() != nil {
		return relationValue.GetListValue().Values
	}
	return []*types.Value{relationValue}
}

func (r *Renderer) generateObjectLinks(relationValue *types.Value) []templ.Component {
	var elements []templ.Component
	for _, value := range r.extractRelationValues(relationValue) {
		objectId := value.GetStringValue()
		snapshot := r.getObjectSnapshot(objectId)
		details := snapshot.GetSnapshot().GetData().GetDetails()
		if details == nil || len(details.GetFields()) == 0 {
			continue
		}

		name := getRelationField(details, bundle.RelationKeyName, relationToString)
		if name == "" {
			name = defaultName
		}
		icon := r.getIconFromDetails(details)
		link := r.makeAnytypeLink(details, objectId)
		elements = append(elements, ListElement(ObjectElement(name, templ.SafeURL(link)), icon))
	}
	return elements
}

func (r *Renderer) generateFileComponent(relationValue *types.Value) []templ.Component {
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
		icon := r.createFileIcon(fileBlock)
		elements = append(elements, ListElement(ObjectElement(filepath.Base(url), templ.URL(url)), icon))
	}
	return elements
}

func (r *Renderer) createFileIcon(fileBlock *model.Block) templ.Component {
	params, err := r.MakeRenderFileParams(fileBlock)
	if err != nil {
		return NoneTemplate(err.Error())
	}

	iconComp := r.FileIconBlock(fileBlock, params)
	return iconComp
}

func (r *Renderer) getIconFromDetails(details *types.Struct) templ.Component {
	props := &IconObjectProps{Size: 20}
	iconParams := r.MakeRenderIconObjectParams(details, props)
	return IconObjectTemplate(r, iconParams)
}

func (r *Renderer) RenderRelations(b *model.Block) templ.Component {
	component := r.makeRelationTemplate(b)
	if component == nil {
		return NoneTemplate("")
	}
	blockParams := makeDefaultBlockParams(b)
	color := b.GetBackgroundColor()
	if color != "" {
		blockParams.Classes = append(blockParams.Classes, fmt.Sprintf("bgColor bgColor-%s", color))
	}
	blockParams.Content = component

	return BlockTemplate(r, blockParams)
}
