package blockutils

import (
	"fmt"

	"github.com/anyproto/anytype-heart/pkg/lib/logging"
	"github.com/anyproto/anytype-heart/pkg/lib/pb/model"
	"github.com/anyproto/anytype-heart/util/pbtypes"

	types "github.com/gogo/protobuf/types"
	"go.uber.org/zap"
)

var log = logging.Logger("renderer.blockutils").Desugar()

func HydrateBlock(block *model.Block, details *types.Struct) (err error) {
	detailsKeys := getDetailsKeys(block)
	if textBlock, ok := block.GetContent().(*model.BlockContentOfText); ok {
		if detailsKeys.Text != "" {
			text := pbtypes.GetString(details, detailsKeys.Text)
			log.Debug("details:", zap.String("id", block.Id), zap.String("text", text))
			textBlock.Text.Text = text
		}
		if detailsKeys.Checked != "" {
			checked := pbtypes.GetBool(details, detailsKeys.Checked)
			textBlock.Text.Checked = checked
		}
	} else {
		err = fmt.Errorf("hydrateBlock: expected block to be BlockContentOfText")
	}

	return
}
