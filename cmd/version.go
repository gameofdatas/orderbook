package cmd

import (
	"fmt"

	"pricer/version"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the application version",
	Long:  "Prints the application version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Goal achiever Version %s\n", version.Version)
		fmt.Printf("Built the %s\n", version.BuildDate)
		fmt.Printf("Last commit %s\n", version.GitCommit)
		fmt.Printf("OS %s\n", version.OsArch)
		fmt.Printf("GO version %s\n", version.GoVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
