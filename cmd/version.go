package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Print the declutter version number",
	Long:    `All software has versions. This is the declutter's'`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("v%s\n", cmd.Root().Version)
	},
}
