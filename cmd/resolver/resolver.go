package resolver

import (
	"fmt"
	"os"
	"path/filepath"
)

type SimpleAssetResolver struct {
	CdnUrl      string
	SnapshotDir string
	RootPageId  string
}

func (r SimpleAssetResolver) GetRootPagePath() string {
	return fmt.Sprintf("objects/%s.pb", r.RootPageId)
}

func (r SimpleAssetResolver) GetJoinSpaceLink() string {
	return "https://invite.any.coop/bafybeib3eh3aowv24v5rv4japcrgqw4ly7fx3h4lca2vbh3kqdy7hasoxe#29SnvBDxo83r5MooE2FdSa6wPmwJKxLkuZEinqfvoCKt"
}

func (r SimpleAssetResolver) GetEmojiUrl(code rune) string {
	return fmt.Sprintf("%s/emojies/%x.png", r.CdnUrl, code)
}

func (r SimpleAssetResolver) GetAssetUrl(source string) string {
	return fmt.Sprintf("/%s", filepath.Join(r.SnapshotDir, source))
}

func (r SimpleAssetResolver) GetSnapshotPbFile(path string) (snapshotData []byte, err error) {
	snapshotData, err = os.ReadFile(filepath.Join(r.SnapshotDir, path))
	if err != nil {
		return
	}

	return
}
