package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeRenderCoverParams(t *testing.T) {
	r := getTestRenderer("snapshot_pb")
	_, err := r.MakeRenderPageCoverParams()
	assert.Error(t, err)

}
