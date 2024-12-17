package cmd

import (
	"fmt"
	"os"

	"github.com/anyproto/anytype-publish-renderer/renderer"
	"github.com/spf13/cobra"
)

var pbCmd = &cobra.Command{
	Use:   `anytype-publish-renderer <snapshot-path>`,
	Args:  cobra.MinimumNArgs(1),
	Short: "Convert Anytype web publish package to HTML",
	Run: func(cmd *cobra.Command, args []string) {
		snapshotPath := args[0]
		config := renderer.RenderConfig{
			StaticFilesPath:  "/static",
			PublishFilesPath: snapshotPath,
			PrismJsCdnUrl:    "https://cdn.jsdelivr.net/npm/prismjs@1.29.0",
			AnytypeCdnUrl:    "https://anytype-static.fra1.cdn.digitaloceanspaces.com",
			AnalyticsCode:    `<script>console.log("sending dummy analytics...")</script>`,
		}

		r, err := renderer.NewRenderer(config)
		if err != nil {
			fmt.Printf("Error creating renderer: %v\n", err)
			return
		}

		err = r.Render(os.Stdout)
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
