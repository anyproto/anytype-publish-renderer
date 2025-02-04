package renderer

import (
	"fmt"
	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/core/domain"
	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/gogo/protobuf/types"
	"go.uber.org/zap"
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

	name, format, found := r.retrieveRelationInfo(relationBlock.GetKey())
	params.Name = name

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

	params.Format = r.getFormatClass(format)
	switch params.Format {
	case "c-object", "c-file", "c-select":
		params.Value = ListTemplate(params.Format, r.populateRelationListValue(format, relationValue))
	default:
		params.Value = WrapToCellTemplate(params.Format, r.populateRelationValue(format, relationValue))
	}
	return params
}

func (r *Renderer) retrieveRelationInfo(key string) (string, model.RelationFormat, bool) {
	relationKey := domain.RelationKey(key)
	relation, _ := bundle.GetRelation(relationKey)

	name, format, found := r.fetchRelationMetadata(relation, relationKey)
	if name == "" {
		name = defaultName
	}
	return name, format, found
}

func (r *Renderer) fetchRelationMetadata(relation *model.Relation, relationKey domain.RelationKey) (string, model.RelationFormat, bool) {
	if relation != nil {
		return relation.Name, relation.Format, true
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
			return name, format, true
		}
	}
	return "", model.RelationFormat_longtext, false
}

func (r *Renderer) populateRelationListValue(format model.RelationFormat, relationValue *types.Value) []templ.Component {
	switch format {
	case model.RelationFormat_status, model.RelationFormat_tag:
		return r.generateSelectOptions(format, relationValue)
	case model.RelationFormat_object:
		return r.generateObjectLinks(relationValue)
	case model.RelationFormat_file:
		return r.generateFileIcons(relationValue)
	}
	return nil
}

func (r *Renderer) populateRelationValue(format model.RelationFormat, relationValue *types.Value) templ.Component {
	switch format {
	case model.RelationFormat_shorttext, model.RelationFormat_longtext:
		return BasicTemplate(relationValue.GetStringValue())
	case model.RelationFormat_number:
		return BasicTemplate(fmt.Sprintf("%f", relationValue.GetNumberValue()))
	case model.RelationFormat_phone, model.RelationFormat_email, model.RelationFormat_url:
		return BasicTemplate(relationValue.GetStringValue())
	case model.RelationFormat_date:
		return BasicTemplate(r.formatDate(relationValue.GetNumberValue()))
	case model.RelationFormat_checkbox:
		return r.generateCheckbox(relationValue.GetBoolValue())
	}
	return nil
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

func (r *Renderer) generateCheckbox(checked bool) templ.Component {
	if checked {
		return ActiveCheckBoxTemplate()
	}
	return DisabledCheckBoxTemplate()
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
		elements = append(elements, OptionElement(name, color, relationType))
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
		spaceId := details.GetFields()[bundle.RelationKeySpaceId.String()].GetStringValue()
		name := details.GetFields()[bundle.RelationKeyName.String()].GetStringValue()
		if name == "" {
			name = defaultName
		}
		icon, class := r.getIconFromDetails(details, "c20")
		layoutClass := getLayoutClass(details)
		link := fmt.Sprintf(linkTemplate, objectId, spaceId)
		elements = append(elements, ObjectsListElement(layoutClass, icon, class, name, templ.SafeURL(link)))
	}
	return elements
}

func (r *Renderer) getObjectSnapshot(objectId string) *pb.SnapshotWithType {
	directories := []string{"objects", "relations", "types", "templates", "filesObjects"}
	var (
		snapshot *pb.SnapshotWithType
		err      error
	)
	for _, dir := range directories {
		path := filepath.Join(dir, objectId+".pb")
		snapshot, err = r.ReadJsonpbSnapshot(path)
		if err == nil {
			return snapshot
		}
	}
	log.Error("failed to get snapshot for object", zap.String("objectId", objectId), zap.Error(err))
	return nil
}

func (r *Renderer) generateFileIcons(relationValue *types.Value) []templ.Component {
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
	return elements
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
