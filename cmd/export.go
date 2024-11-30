package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/eddymoulton/onekube/internal/items"
	"github.com/eddymoulton/onekube/internal/onepassword"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Downloads the requested configuration and outputs the path",
	Long:  `Call eval $(onekube set ...) to set the KUBECONFIG environment variable`,
	Run: func(cmd *cobra.Command, args []string) {
		client := onepassword.NewOpClient()

		allConfigItems, err := items.Load(client, false)
		if err != nil {
			log.Fatal(err)
		}

		// Create ~/.kube/onekube directory if it doesn't exist
		onekubePath := fmt.Sprintf("%s/.kube/onekube", os.Getenv("HOME"))
		if err := os.MkdirAll(onekubePath, 0755); err != nil {
			log.Fatal(err)
		}

		if len(args) == 0 {
			// Export all configs
			for _, item := range allConfigItems {
				if err := exportConfig(client, item, onekubePath); err != nil {
					log.Printf("Failed to export config for %s: %v", item.Title, err)
				}
			}
			return
		}

		if len(args) > 1 {
			log.Fatal("Please provide a single name only")
		}

		// Export single config
		itemName := args[0]
		item, err := items.Find(allConfigItems, itemName)
		if err != nil {
			log.Fatal(err)
		}

		if err := exportConfig(client, item, onekubePath); err != nil {
			log.Fatal(err)
		}
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		client := onepassword.NewOpClient()
		items, err := items.Load(client, false)

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

func exportConfig(client *onepassword.Client, item onepassword.Item, basePath string) error {
	rawKubeConfig, err := client.ReadItemField(item.Vault.ID, item.ID, "config")
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	updatedKubeConfig := fmt.Sprintf("# Managed by onekube\n%s", rawKubeConfig)
	configPath := fmt.Sprintf("%s/%s.yaml", basePath, item.Title)

	if err := os.WriteFile(configPath, []byte(updatedKubeConfig), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("Exported kubeconfig for '%s' to %s\n", item.Title, configPath)
	return nil
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
