package blockutils

import (
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"
)

const DetailsKeyFieldName = "_detailsKey"

type DetailsKeys struct {
	Text    string
	Checked string
}

func getKeysList(block *model.Block) []string {
	return pbtypes.GetStringList(block.Fields, DetailsKeyFieldName)
}

func getDetailsForKeyslist(keysList []string) (keys DetailsKeys) {
	if len(keysList) > 1 {
		keys.Text = keysList[0]
		keys.Checked = keysList[1]
		return
	}
	if len(keysList) > 0 {
		keys.Text = keysList[0]
	}

	return
}

func getDetailsKeys(block *model.Block) (keys DetailsKeys) {
	keysList := getKeysList(block)
	return getDetailsForKeyslist(keysList)
}
