package renderer

import (
	"testing"

	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/stretchr/testify/assert"
)

func TestMakeRenderFileParams(t *testing.T) {
	r := getTestRenderer()
	id := "66c7055b7e4bcd7bc81f3f37"
	imageBlock := r.BlocksById[id]

	expected := &FileRenderParams{
		Id:      id,
		Classes: "align1",
		Src:     "/../snapshot_pb/files/img_5296.jpeg",
		Type:    model.BlockContentFile_Image,
	}

	actual, err := r.MakeRenderFileParams(imageBlock)
	if assert.NoError(t, err) {
		assert.Equal(t, expected.Id, actual.Id)
		assert.Equal(t, expected.Classes, actual.Classes)
		assert.Equal(t, expected.Src, actual.Src)
		assert.Equal(t, expected.Type, actual.Type)

	}
}
