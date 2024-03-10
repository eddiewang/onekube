/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/eddymoulton/onekube/internal/funcs"
	"github.com/eddymoulton/onekube/onepassword"
	"github.com/spf13/cobra"
)

var Force bool

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available configs from 1password",
	Long: `Lists available configs from 1password currently stored locally.
Force the update to re-check what's available in 1password`,
	Run: func(cmd *cobra.Command, args []string) {
		client := onepassword.NewOpClient()

		items, err := funcs.LoadItems(client, Force)

		if err != nil {
			log.Fatal(err)
		}

		for _, item := range items {
			fmt.Println(item.Title)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&Force, "force", "f", false, "Force reload from one password")
}
