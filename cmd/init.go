package cmd

import (
	"github.com/eddymoulton/onekube/internal/config"
	"github.com/eddymoulton/onekube/internal/items"
	"github.com/eddymoulton/onekube/internal/onepassword"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Cleans and reloads available items from 1password",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client := onepassword.NewOpClient()

		config.Clean()
		items.Load(client, true)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

}
