package cmd

import (
	"github.com/spf13/cobra"

	"caraxes/internal/repo"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Shows commit history",
	Long:  "This command walks the commit history starting from the current branch",

	RunE: func(cmd *cobra.Command, args []string) error {
		return repo.LogRepo()
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
	