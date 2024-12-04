package cmd

import (
	"declutter/tools/helpers"
	"github.com/spf13/cobra"
)

var sortCmd = &cobra.Command{
	Use:     "sort",
	Aliases: []string{"s"},
	Short:   "Clean directory files",
	Long:    `Clean the current directory's files into folders based on their extension`,
	Args:    cobra.ExactArgs(1),
	Run:     run,
}

func init() {
	sortCmd.Flags().Bool("dry-run", false, "Runs the command without actually doing anything")
	sortCmd.Flags().StringArray("exclude", []string{}, "Exclude files from the cleaning process")
	rootCmd.AddCommand(sortCmd)
}

func run(cmd *cobra.Command, args []string) {
	isDryRun, err := cmd.Flags().GetBool("dry-run")
	if err != nil {
		panic(err)
	}
	array, err := cmd.Flags().GetStringArray("exclude")
	if err != nil {
		return
	}
	fileSorter := helpers.InitializeFileHandler(args[0], array).FilterDirectories().SortFiles()
	// Check if it is a dry run
	if isDryRun {
		fileSorter.DryRun()
	} else {
		fileSorter.MoveFiles()
	}
}
