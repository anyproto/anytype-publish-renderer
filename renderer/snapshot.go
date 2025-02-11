package renderer

import (
	"path/filepath"
	"strings"

	"go.uber.org/zap"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/localstore/addr"
)

func (r *Renderer) getObjectSnapshot(objectId string) *pb.SnapshotWithType {
	if strings.HasPrefix(objectId, addr.DatePrefix) {
		return r.getDateSnapshot(objectId)
	}
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
