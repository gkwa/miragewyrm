package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/gkwa/miragewyrm/list"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List S3 objects recursively",
	Long:  `List objects in an S3 bucket recursively, filtering by file extension`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := LoggerFrom(cmd.Context())
		lister := list.NewS3Lister(logger, bucket)
		if err := lister.List(context.Background()); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to list objects: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
