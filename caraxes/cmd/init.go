package cmd

import (
	"github.com/spf13/cobra"

	"caraxes/internal/banner"
	"caraxes/internal/repo"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes an empty caraxes repo",
	Long:  "This command initializes an empty caraxes repo",

	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}
		if err := repo.InitRepo(path); err != nil {
			return err
		}
		banner.Print()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
