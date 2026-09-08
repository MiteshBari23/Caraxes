package cmd

import (
	"caraxes/internal/repo"


	"github.com/spf13/cobra"
)


var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Shows the working tree status",
	Long:  "This command shows untracked and modified file compared to the staging area",

	RunE: func(cmd *cobra.Command, args []string) error {
		return repo.StatusRepo()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
