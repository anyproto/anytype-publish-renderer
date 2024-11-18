package renderer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeRenderCoverParams(t *testing.T) {
	t.Run("no cover returns error", func(t *testing.T) {
		r := getTestRenderer("snapshot_pb")
		_, err := r.MakeRenderPageCoverParams()
		assert.Error(t, err)
	})

	t.Run("cover params", func(t *testing.T) {
		r := getTestRenderer("snapshot_pb2")
		expected := &CoverRenderParams{
			Id:      "bafyreie4vr5erfeecbhh5ocs6wzqmab2c3rndmp3e7aovhmcngj4geylf4",
			Classes: "type1",
			Src:     "/../test_snapshots/snapshot_pb2/files/1729523760-21-10-24_17-16-00.png",
		}

		actual, err := r.MakeRenderPageCoverParams()
		if assert.NoError(t, err) {
			assert.Equal(t, expected.Id, actual.Id)
			assert.Equal(t, expected.Classes, actual.Classes)
			assert.Equal(t, expected.Src, actual.Src)
		}
	})
}
