package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/gkwa/miragewyrm/fetch"
)

var (
	count  int
	outDir string
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch files from S3",
	Long:  `Fetch files from an S3 bucket with various options`,
}

var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Randomly fetch files",
	Long:  `Randomly select and fetch files from S3 that don't exist locally`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := LoggerFrom(cmd.Context())
		fetcher := fetch.NewS3Fetcher(logger, bucket)
		fetcher.SetOutputDir(outDir)
		if err := fetcher.FetchRandom(context.Background(), count); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to fetch files: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
	fetchCmd.AddCommand(randomCmd)

	randomCmd.Flags().IntVarP(&count, "count", "n", 1, "Number of random files to fetch")
	randomCmd.Flags().StringVarP(&outDir, "outdir", "o", ".", "Output directory for downloaded files")
}
