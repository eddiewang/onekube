package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/eddymoulton/onekube/internal/items"
	"github.com/eddymoulton/onekube/internal/onepassword"
	"github.com/spf13/cobra"
)

var (
	pushAllFiles bool
	vaultName    string
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push kubeconfig files to 1Password",
	Long:  `Push kubeconfig files to 1Password with the 'kubeconfig' tag`,
	Run: func(cmd *cobra.Command, args []string) {
		client := onepassword.NewOpClient()

		// Get the directory where kubeconfig files are stored
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		// Location for onekube kubeconfig files
		onekubeDir := filepath.Join(home, ".kube", "onekube")

		// Create the directory if it doesn't exist
		if err := os.MkdirAll(onekubeDir, 0755); err != nil {
			log.Fatal(err)
		}

		if len(args) == 0 && !pushAllFiles {
			log.Fatal("Please provide a file path or use --all flag")
		}

		if pushAllFiles {
			// Push all kubeconfig files from .kube/onekube directory
			fmt.Println("Pushing all kubeconfig files from .kube/onekube directory...")

			// Read all files in the .kube/onekube directory
			files, err := os.ReadDir(onekubeDir)
			if err != nil {
				log.Fatal(err)
			}

			// Process each YAML file
			for _, file := range files {
				// Skip directories and non-YAML files
				if file.IsDir() || (!strings.HasSuffix(file.Name(), ".yaml") && !strings.HasSuffix(file.Name(), ".yml") && file.Name() != "config") {
					continue
				}

				filePath := filepath.Join(onekubeDir, file.Name())
				pushKubeconfig(client, filePath)
			}

			return
		}

		// Push specific files
		for _, filePath := range args {
			pushKubeconfig(client, filePath)
		}

		// Reload items from 1Password
		fmt.Println("Reloading configurations from 1Password...")
		items.Load(client, true)
	},
}

func pushKubeconfig(client *onepassword.Client, filePath string) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("File not found: %s", filePath)
		return
	}

	// Read the kubeconfig file
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Failed to read file %s: %v", filePath, err)
		return
	}

	// Remove the "# Managed by onekube" line if present
	configContent := string(content)
	configContent = strings.ReplaceAll(configContent, "# Managed by onekube\n", "")

	// Create a title from the filename (without extension)
	fileName := filepath.Base(filePath)
	title := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	// If the file is named "config" (default kubeconfig), use a more descriptive name
	if title == "config" {
		title = "default-kubeconfig"
	}

	// Create the 1Password item
	fmt.Printf("Saving kubeconfig '%s' to 1Password...\n", title)

	// Check if the item already exists
	listArgs := []string{
		"item", "list",
		"--tags", "kubeconfig",
		"--vault", vaultName,
		"--format", "json",
	}

	listCmd := exec.Command("op", listArgs...)
	var listOutput bytes.Buffer
	listCmd.Stdout = &listOutput
	var listStderr bytes.Buffer
	listCmd.Stderr = &listStderr

	err = listCmd.Run()
	if err != nil {
		errMsg := listStderr.String()
		if errMsg != "" {
			log.Printf("Failed to check if item exists: %v - %s", err, errMsg)
		} else {
			log.Printf("Failed to check if item exists: %v", err)
		}
		// Continue anyway, trying to create the item
	}

	// Check if the item already exists
	var existingItems []map[string]interface{}
	var itemID string

	if listOutput.Len() > 0 {
		if err := json.Unmarshal(listOutput.Bytes(), &existingItems); err == nil {
			// Look for an item with the matching title
			for _, item := range existingItems {
				if itemTitle, ok := item["title"].(string); ok && itemTitle == title {
					// Found the item
					itemID = item["id"].(string)
					break
				}
			}
		}
	}

	if itemID != "" {
		// Item exists, update it
		fmt.Printf("Updating existing kubeconfig '%s' in 1Password...\n", title)

		args := []string{
			"item", "edit", itemID,
			"--vault", vaultName,
			fmt.Sprintf("config=%s", configContent),
			"--format", "json",
		}

		cmd := exec.Command("op", args...)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		err = cmd.Run()
		if err != nil {
			errMsg := stderr.String()
			if errMsg != "" {
				log.Printf("Failed to update kubeconfig '%s' in 1Password: %v - %s", title, err, errMsg)
			} else {
				log.Printf("Failed to update kubeconfig '%s' in 1Password: %v", title, err)
			}
			return
		}

		fmt.Printf("Successfully updated kubeconfig '%s' in 1Password\n", title)
		return
	}

	// Item doesn't exist, create it
	fmt.Printf("Creating new kubeconfig '%s' in 1Password...\n", title)
	args := []string{
		"item", "create",
		"--category", "Secure Note",
		"--title", title,
		"--tags", "kubeconfig",
		"--vault", vaultName,
		fmt.Sprintf("config=%s", configContent),
		"--format", "json",
	}

	cmd := exec.Command("op", args...)

	// Capture stderr to get detailed error messages
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		errMsg := stderr.String()
		if errMsg != "" {
			log.Printf("Failed to save kubeconfig '%s' to 1Password: %v - %s", title, err, errMsg)
		} else {
			log.Printf("Failed to save kubeconfig '%s' to 1Password: %v", title, err)
		}
		return
	}

	fmt.Printf("Successfully saved kubeconfig '%s' to 1Password\n", title)
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().BoolVarP(&pushAllFiles, "all", "a", false, "Push all kubeconfig files from .kube/onekube directory")
	pushCmd.Flags().StringVarP(&vaultName, "vault", "v", "Homelab", "1Password vault to store kubeconfigs")
}
