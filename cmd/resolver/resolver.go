package resolver

import (
	"fmt"
	"path/filepath"
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

func (r SimpleAssetResolver) ByEmojiCode(code rune) string {
	return fmt.Sprintf("%s/emojies/%x.png", r.CdnUrl, code)
}

func (r SimpleAssetResolver) ByTargetObjectId(id string) (path string, err error) {
	filePbPath := filepath.Join(r.SnapshotDir, "filesObjects", id+".pb")
	return filePbPath, nil
}
