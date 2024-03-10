package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/eddymoulton/onekube/internal/funcs"
	"github.com/eddymoulton/onekube/onepassword"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Downloads the requested configuration and outputs the path",
	Long:  `Call eval $(onekube set ...) to set the KUBECONFIG environment variable`,
	Run: func(cmd *cobra.Command, args []string) {
		client := onepassword.NewOpClient()

		items, err := funcs.LoadItems(client, false)

		if err != nil {
			log.Fatal(err)
		}

		if len(args) == 0 {
			log.Fatal("Missing name of config")
		}

		if len(args) > 1 {
			log.Fatal("Please provide a single name only")
		}

		itemName := args[0]

		item, err := funcs.FindItem(items, itemName)

		if err != nil {
			log.Fatal(err)
		}

		kubeConfigFile := funcs.GetKubeConfigFilePath(itemName)

		if _, err := os.Stat(kubeConfigFile); errors.Is(err, os.ErrNotExist) {
			config, err := client.ReadItemField(item.Vault.ID, item.ID, "config")

			if err != nil {
				log.Fatal(err)
			}

			err = os.WriteFile(kubeConfigFile, []byte(config), 0644)

			if err != nil {
				log.Fatal(err)
			}
		}

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("export KUBECONFIG=%s\n", kubeConfigFile)
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		client := onepassword.NewOpClient()
		items, err := funcs.LoadItems(client, false)

		if err != nil {
			log.Fatal(err)
		}

		var names []string

		for _, item := range items {
			names = append(names, item.Title)
		}

		return names, cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
