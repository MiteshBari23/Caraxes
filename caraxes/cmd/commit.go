package cmd

import (
	"fmt"
	"caraxes/internal/repo"

	"github.com/spf13/cobra"
)

var commitMessage string

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Records staged changes as a new commit",
	Long:  "This command created a new commit from what's currently in the index",

	RunE: func(cmd *cobra.Command, args []string) error {
		if commitMessage == "" {
			return fmt.Errorf("commit message is required, use -m")
		}
		return repo.CommitRepo(commitMessage)
	},
}

func init() {
	commitCmd.Flags().StringVarP(&commitMessage, "message", "m", "", "commit message")
	rootCmd.AddCommand(commitCmd)
}
