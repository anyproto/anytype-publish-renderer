package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	inputFile  string
	outputFile string
)

var pbCmd = &cobra.Command{
	Use:   "anytype-publish-renderer",
	Short: "Convert Anytype exported protobuf to HTML",
	Run: func(cmd *cobra.Command, args []string) {
		snapshotData, err := os.ReadFile(inputFile)
		if err != nil {
			fmt.Printf("Error reading protobuf snapshot: %v\n", err)
			return
		}

		// page, err := models.NewPageFromPbSnapshot(snapshotData)
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

func init() {
	pbCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input snapshot file")
	pbCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output HTML file")
}
