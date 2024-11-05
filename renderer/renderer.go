package renderer

import (
	"context"
	"fmt"
	"strings"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"

	"github.com/a-h/templ"
	"github.com/gogo/protobuf/proto"
)

type Renderer struct {
	Sp   *pb.SnapshotWithType
	Html *strings.Builder
}

func NewRenderer(snapshotData []byte) (r *Renderer, err error) {
	snapshot := pb.SnapshotWithType{}
	err = proto.Unmarshal(snapshotData, &snapshot)
	if err != nil {
		return
	}

	if snapshot.SbType != model.SmartBlockType_Page {
		err = fmt.Errorf("published snaphost is not Page, %d", snapshot.SbType)
		return
	}

	r = &Renderer{
		Sp:   &snapshot,
		Html: &strings.Builder{},
	}

	return

}

func (r *Renderer) HTML() string {
	return r.Html.String()
}

func (r *Renderer) templToString(component templ.Component) (err error) {
	err = component.Render(context.Background(), r.Html)
	return
}
