package resolver

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"

	"github.com/gogo/protobuf/proto"
)

type SimpleAssetResolver struct {
	CdnUrl      string
	SnapshotDir string
	RootPageId  string
}

func (r SimpleAssetResolver) GetRootPagePath() string {
	rootPbPath := filepath.Join(r.SnapshotDir, "objects", r.RootPageId+".pb")
	return rootPbPath
}

func (r SimpleAssetResolver) GetJoinSpaceLink() string {
	return "https://invite.any.coop/bafybeib3eh3aowv24v5rv4japcrgqw4ly7fx3h4lca2vbh3kqdy7hasoxe#29SnvBDxo83r5MooE2FdSa6wPmwJKxLkuZEinqfvoCKt"
}

func (r SimpleAssetResolver) ByEmojiCode(code rune) string {
	return fmt.Sprintf("%s/emojies/%x.png", r.CdnUrl, code)
}

func (r SimpleAssetResolver) ByTargetObjectId(id string) (path string, err error) {
	filePbPath := filepath.Join(r.SnapshotDir, "filesObjects", id+".pb")
	snapshotData, err := r.GetSnapshotPbFile(filePbPath)
	if err != nil {
		return
	}

	snapshot := pb.SnapshotWithType{}
	err = proto.Unmarshal(snapshotData, &snapshot)
	if err != nil {
		return
	}

	if snapshot.SbType != model.SmartBlockType_FileObject {
		err = fmt.Errorf("snaphot %s is not FileObjects, %d", filePbPath, snapshot.SbType)
		return
	}

	fields := snapshot.Snapshot.Data.GetDetails()
	source := pbtypes.GetString(fields, "source")
	path = "/" + filepath.Join(r.SnapshotDir, source)
	return
}

func (r SimpleAssetResolver) GetSnapshotPbFile(path string) (snapshotData []byte, err error) {
	snapshotData, err = os.ReadFile(path)
	if err != nil {
		return
	}

	return
}
