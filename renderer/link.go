package renderer

import (
	"fmt"
	"path/filepath"

	"strings"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/gogo/protobuf/types"
	"go.uber.org/zap"
)

const linkTemplate = "anytype://object?objectId=%s&spaceId=%s"

type LinkRenderParams struct {
	Id             string
	Classes        string
	ContentClasses string
	SidesClasses   string
	CardClasses    string
	IsDeleted      bool
	IsArchived     string
	Name           string
	Description    string
	Type           string
	IconTemplate   templ.Component
	CoverTemplate  templ.Component
	Url            templ.SafeURL
}

func (r *Renderer) MakeLinkRenderParams(b *model.Block) *LinkRenderParams {
	targetObjectId := b.GetLink().GetTargetBlockId()
	targetDetails := r.findTargetDetails(targetObjectId)
	if targetDetails == nil || len(targetDetails.Fields) == 0 {
		return &LinkRenderParams{IsDeleted: true}
	}

	linkTypeClass := strings.ToLower(b.GetLink().GetCardStyle().String())
	description := getDescription(b, targetDetails)
	if isDeleted(targetDetails) {
		return &LinkRenderParams{IsDeleted: true}
	}

	bgColor := b.GetBackgroundColor()
	name := getFieldValue(targetDetails, bundle.RelationKeyName.String(), defaultName)
	layoutClass := getLayoutClass(targetDetails)
	archiveClass := getArchiveClass(targetDetails)
	objectTypeName, coverTemplate := r.getAdditionalParams(b, targetDetails)
	spaceId := targetDetails.GetFields()[bundle.RelationKeySpaceId.String()].GetStringValue()
	link := fmt.Sprintf(linkTemplate, targetObjectId, spaceId)
	classes := []string{linkTypeClass, archiveClass}
	contentClasses := []string{"content"}
	sidesClasses := []string{"sides"}
	cardClasses := []string{"linkCard", layoutClass}

	if bgColor != "" {
		sidesClasses = append(sidesClasses, "withBgColor")
		contentClasses = append(contentClasses, "bgColor", "bgColor-"+bgColor)
	}

	size, iconSize := getLinkIconSize(b)

	params := r.MakeRenderIconObjectParams(targetDetails, &IconObjectProps{
		Size:     int32(size),
		IconSize: int32(iconSize),
	})
	iconTemplate := IconObjectTemplate(r, params)

	if iconTemplate != nil {
		cardClasses = append(cardClasses, "withIcon", fmt.Sprintf("c%d", size))
	}

	if coverTemplate != nil {
		cardClasses = append(cardClasses, "withCover")
	}

	n := 1
	if description != "" {
		n++
	}
	if objectTypeName != "" {
		n++
	}

	cardClasses = append(cardClasses, fmt.Sprintf("c%d", n))

	return &LinkRenderParams{
		Id:             b.GetId(),
		Classes:        strings.Join(classes, " "),
		ContentClasses: strings.Join(contentClasses, " "),
		SidesClasses:   strings.Join(sidesClasses, " "),
		CardClasses:    strings.Join(cardClasses, " "),
		IsArchived:     archiveClass,
		Name:           name,
		Description:    description,
		Type:           objectTypeName,
		Url:            templ.URL(link),
		CoverTemplate:  coverTemplate,
		IconTemplate:   iconTemplate,
	}
}

func getDescription(b *model.Block, details *types.Struct) string {
	var key string

	switch b.GetLink().GetDescription() {
	case model.BlockContentLink_Content:
		key = bundle.RelationKeySnippet.String()
	case model.BlockContentLink_Added:
		key = bundle.RelationKeyDescription.String()
	default:
		return ""
	}

	descriptionValue := details.GetFields()[key]
	if descriptionValue != nil {
		return descriptionValue.GetStringValue()
	}
	return ""
}

func isDeleted(details *types.Struct) bool {
	deletedValue := details.GetFields()[bundle.RelationKeyIsDeleted.String()]
	return deletedValue != nil && deletedValue.GetBoolValue()
}

func getFieldValue(details *types.Struct, key, defaultValue string) string {
	value := details.GetFields()[key]
	if value == nil || value.GetStringValue() == "" {
		return defaultValue
	}
	return value.GetStringValue()
}

func getLayoutClass(details *types.Struct) string {
	layout := model.ObjectTypeLayout(details.GetFields()[bundle.RelationKeyLayout.String()].GetNumberValue())

	switch layout {
	case model.ObjectType_participant:
		return "isParticipant"
	case model.ObjectType_profile:
		return "isHuman"
	case model.ObjectType_todo:
		return "isTask"
	case model.ObjectType_collection:
		return "isCollection"
	case model.ObjectType_set:
		return "isSet"
	default:
		return "isPage"
	}
}

func getArchiveClass(details *types.Struct) string {
	archivedValue := details.GetFields()[bundle.RelationKeyIsArchived.String()]
	if archivedValue != nil && archivedValue.GetBoolValue() {
		return "isArchived"
	}
	return ""
}

func getLinkIconSize(b *model.Block) (int, int) {
	link := b.GetLink()
	cardStyle := link.GetCardStyle()
	iconSize := link.GetIconSize()

	newSize := 20
	newIconSize := 20

	if (cardStyle != model.BlockContentLink_Text) && (iconSize == model.BlockContentLink_SizeMedium) {
		newSize = 48
		newIconSize = 28
	}

	return newSize, newIconSize
}

func (r *Renderer) getAdditionalParams(b *model.Block, details *types.Struct) (objectTypeName string, coverTemplate templ.Component) {
	for _, relation := range b.GetLink().GetRelations() {
		if relation == bundle.RelationKeyType.String() {
			objectType := details.GetFields()[bundle.RelationKeyType.String()].GetStringValue()
			snapshot, err := r.ReadJsonpbSnapshot(filepath.Join("types", objectType+".pb"))
			if err != nil {
				log.Error("failed to read jsonpb snapshot", zap.Error(err))
				continue
			}
			objectTypeName = snapshot.GetSnapshot().GetData().GetDetails().GetFields()[bundle.RelationKeyName.String()].GetStringValue()
		}

		if relation == "cover" {
			var err error
			coverParams, err := r.getCoverParams(details, false, false)
			if err == nil {
				coverTemplate = coverParams.CoverTemplate
			}
		}
	}
	return
}

func (r *Renderer) RenderLink(b *model.Block) templ.Component {
	params := r.MakeLinkRenderParams(b)
	return LinkTempl(params)
}
