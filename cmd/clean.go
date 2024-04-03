package cmd

import (
	"github.com/eddymoulton/onekube/internal/config"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Deletes all data created by onekube",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		config.Clean()
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
