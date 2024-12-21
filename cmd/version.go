package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gkwa/miragewyrm/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of miragewyrm",
	Long:  `All software has versions. This is miragewyrm's`,
	Run: func(cmd *cobra.Command, args []string) {
		buildInfo := version.GetBuildInfo()
		fmt.Println(buildInfo)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
