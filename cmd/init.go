package cmd

import (
	"github.com/eddymoulton/onekube/internal/funcs"
	"github.com/eddymoulton/onekube/onepassword"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Cleans and reloads available items from 1password",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client := onepassword.NewOpClient()

		funcs.CleanData()
		funcs.LoadItems(client, true)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

}
