package cmd

import (
	"github.com/eddymoulton/onekube/internal/funcs"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Deletes all data created by onekube",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		funcs.CleanData()
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
