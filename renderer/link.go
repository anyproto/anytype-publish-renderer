package renderer

import (
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/anyproto/anytype-heart/pkg/lib/bundle"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/gogo/protobuf/types"
)

type LinkRenderParams struct {
	Id           string
	SidesClasses string
	CardClasses  string
	Url          templ.SafeURL
	Components   []templ.Component
}

func (r *Renderer) makeLinkBlockParams(b *model.Block) *BlockParams {
	blockParams, sidesClasses := r.fillContentClasses(b)
	targetObjectId := b.GetLink().GetTargetBlockId()
	targetDetails := r.findTargetDetails(targetObjectId)
	if targetDetails == nil || len(targetDetails.Fields) == 0 || getRelationField(targetDetails, bundle.RelationKeyIsDeleted, relationToBool) {
		var icon templ.Component
		icon, blockParams.Classes = r.getIconTemplate(targetDetails, 20, 20, blockParams.Classes)
		blockParams.Content = DeletedLinkTemplate(icon)
		return blockParams
	}
	archiveClass := r.fillArchiveClass(targetDetails, b.GetLink().GetCardStyle(), blockParams)
	cardClasses := r.fillLayoutClass(targetDetails)

	name := getNameValue(targetDetails, bundle.RelationKeyName.String(), defaultName)
	description := getDescription(b, targetDetails)

	iconTemplate, cardClasses := r.getIconComponent(b, targetDetails, cardClasses)

	objectTypeName, coverTemplate := r.getAdditionalParams(b, targetDetails)
	linkComponents, cardClasses := r.getLinkComponent(coverTemplate, iconTemplate, cardClasses, name, description, objectTypeName, archiveClass)

	lp := &LinkRenderParams{
		Id:           b.GetId(),
		SidesClasses: strings.Join(sidesClasses, " "),
		CardClasses:  strings.Join(cardClasses, " "),
		Url:          templ.SafeURL(makeAnytypeLink(targetDetails, targetObjectId)),
		Components:   linkComponents,
	}
	blockParams.Content = LinkTemplate(lp)
	return blockParams
}

func (r *Renderer) fillContentClasses(b *model.Block) (*BlockParams, []string) {
	blockParams := makeDefaultBlockParams(b)
	bgColor := b.GetBackgroundColor()
	sidesClasses := []string{"sides"}
	if bgColor != "" {
		sidesClasses = append(sidesClasses, "withBgColor")
		blockParams.ContentClasses = append(blockParams.ContentClasses, "bgColor", "bgColor-"+bgColor)
	}
	return blockParams, sidesClasses
}

func (r *Renderer) fillArchiveClass(targetDetails *types.Struct, style model.BlockContentLinkCardStyle, blockParams *BlockParams) string {
	blockParams.Classes = append(blockParams.Classes, strings.ToLower(style.String()))
	archiveClass := getArchiveClass(targetDetails)
	if archiveClass != "" {
		blockParams.Classes = append(blockParams.Classes, archiveClass)
	}
	return archiveClass
}

func (r *Renderer) fillLayoutClass(targetDetails *types.Struct) []string {
	layout := getRelationField(targetDetails, bundle.RelationKeyLayout, relationToObjectTypeLayout)
	layoutClass := getLayoutClass(layout)
	cardClasses := []string{"linkCard", layoutClass}
	return cardClasses
}

func (r *Renderer) getIconComponent(b *model.Block, targetDetails *types.Struct, cardClasses []string) (templ.Component, []string) {
	size, iconSize := getLinkIconSize(b)

	iconTemplate := NoneTemplate("")

	if size != 0 {
		return r.getIconTemplate(targetDetails, size, iconSize, cardClasses)
	}
	return iconTemplate, cardClasses
}

func (r *Renderer) getIconTemplate(targetDetails *types.Struct, size, iconSize int, cardClasses []string) (templ.Component, []string) {
	params := r.MakeRenderIconObjectParams(targetDetails, &IconObjectProps{
		Size:     int32(size),
		IconSize: int32(iconSize),
	})
	iconTemplate := IconObjectTemplate(r, params)

	if iconTemplate != nil {
		cardClasses = append(cardClasses, "withIcon", fmt.Sprintf("c%d", size))
	}
	return iconTemplate, cardClasses
}

func (r *Renderer) getLinkComponent(coverTemplate, iconTemplate templ.Component, cardClasses []string, name, description, objectTypeName, archiveClass string) ([]templ.Component, []string) {
	var cardComponents, sideLeftComponents, sideRightComponents []templ.Component
	if coverTemplate != nil {
		cardClasses = append(cardClasses, "withCover")
		sideRightComponents = append(sideRightComponents, BlocksWrapper(&BlockWrapperParams{
			Classes:    []string{"side right"},
			Components: []templ.Component{coverTemplate},
		}))
	}
	cardComponents = append(cardComponents, iconTemplate, BasicTemplate("name", name))
	if archiveClass != "" {
		cardComponents = append(cardComponents, ArchivedLinkTemplate())
	}
	sideLeftComponents = append(sideLeftComponents, BlocksWrapper(&BlockWrapperParams{
		Classes:    []string{"cardName"},
		Components: cardComponents,
	}))
	n := 1
	if description != "" {
		n++
		sideLeftComponents = append(sideLeftComponents, LinkItemTemplate("cardDescription", "description", description))
	}
	if objectTypeName != "" {
		n++
		sideLeftComponents = append(sideLeftComponents, LinkItemTemplate("cardType", "item", objectTypeName))
	}
	cardClasses = append(cardClasses, fmt.Sprintf("c%d", n))
	linkComponents := []templ.Component{
		BlocksWrapper(&BlockWrapperParams{
			Classes:    []string{"side left"},
			Components: sideLeftComponents,
		}),
	}
	linkComponents = append(linkComponents, sideRightComponents...)
	return linkComponents, cardClasses
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

func getNameValue(details *types.Struct, key, defaultValue string) string {
	value := details.GetFields()[key]
	if value == nil || value.GetStringValue() == "" {
		return defaultValue
	}
	return value.GetStringValue()
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

	if iconSize == model.BlockContentLink_SizeNone {
		return 0, 0
	}

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
			objectType := getRelationField(details, bundle.RelationKeyType, relationToString)
			snapshot := r.getObjectSnapshot(objectType)
			if snapshot == nil {
				continue
			}
			objectTypeName = getRelationField(snapshot.GetSnapshot().GetData().GetDetails(), bundle.RelationKeyName, relationToString)
		}

		if relation == "cover" {
			var err error
			coverParams, err := r.getCoverParams(details, false, false, true)
			if err == nil {
				coverTemplate = coverParams.CoverTemplate
			}
		}
	}
	return
}

func (r *Renderer) RenderLink(b *model.Block) templ.Component {
	params := r.makeLinkBlockParams(b)
	return BlockTemplate(r, params)
}
