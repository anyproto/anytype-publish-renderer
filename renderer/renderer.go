package renderer

import (
	"context"
	"fmt"
	"io"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"

	"github.com/a-h/templ"
	"github.com/gogo/protobuf/proto"
)

type Renderer struct {
	Sp       *pb.SnapshotWithType
	Out      io.Writer
	RootComp templ.Component
}

func NewRenderer(snapshotData []byte, writer io.Writer) (r *Renderer, err error) {
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
		Sp:  &snapshot,
		Out: writer,
	}

	return
}

func (r *Renderer) Render() (err error) {
	err = r.RootComp.Render(context.Background(), r.Out)
	if err != nil {
		return
	}
	return
}
