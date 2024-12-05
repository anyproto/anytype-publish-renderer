package cmd

import (
	"fmt"
	"os"

	"github.com/anyproto/anytype-publish-renderer/cmd/resolver"
	"github.com/anyproto/anytype-publish-renderer/renderer"
	"github.com/spf13/cobra"
)

const CDN_URL = "https://anytype-static.fra1.cdn.digitaloceanspaces.com"

var pbCmd = &cobra.Command{
	Use: `anytype-publish-renderer <snapshot-dir> <root-page-id>

<root-page-id> is one of the pages .pb files, typically stored in snapshot-folder/objects/*.pb
`,
	Args:  cobra.MinimumNArgs(2),
	Short: "Convert Anytype exported protobuf folder to HTML",
	Run: func(cmd *cobra.Command, args []string) {
		snapshotDir := args[0]
		rootId := args[1]
		resolver := resolver.SimpleAssetResolver{
			CdnUrl:      CDN_URL,
			SnapshotDir: snapshotDir,
			RootPageId:  rootId,
		}
		r, err := renderer.NewRenderer(resolver, os.Stdout)
		if err != nil {
			fmt.Printf("Error rendering snapshot: %v\n", err)
			return
		}

		r.RootComp = r.RenderPage()
		err = r.Render()
		if err != nil {
			fmt.Printf("Error rendering page: %v\n", err)
			return
		}
	},
}

func Execute() {
	if err := pbCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
