package cmd

import (
	"fmt"
	"os"

	"github.com/anyproto/anytype-publish-renderer/renderer"
	"github.com/spf13/cobra"
)

var pbCmd = &cobra.Command{
	Use:   "anytype-publish-renderer <objects/snapshot.pb>",
	Args:  cobra.MinimumNArgs(1),
	Short: "Convert Anytype exported protobuf to HTML",
	Run: func(cmd *cobra.Command, args []string) {
		snapshotPath := args[0]
		snapshotData, err := os.ReadFile(snapshotPath)
		if err != nil {
			fmt.Printf("Error reading protobuf snapshot: %v\n", err)
			return
		}
		r, err := renderer.NewRenderer(snapshotData)
		if err != nil {
			fmt.Printf("Error rendering snapshot: %v\n", err)
			return
		}

		err = r.RenderPage()
		if err != nil {
			fmt.Printf("Error rendering page: %v\n", err)
			return
		}
		fmt.Println(r.HTML())
		// if err != nil {
		// 	fmt.Printf("Error reading creating page from snapshot: %v\n", err)
		// 	return
		// }

		// fmt.Printf("%#v\n", page)

	},
}

func Execute() {
	if err := pbCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
