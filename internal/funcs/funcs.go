package funcs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/eddymoulton/onekube/onepassword"
	"github.com/iancoleman/strcase"
)

func GetKubeConfigFilePath(name string) string {
	configDirectory := getConfigDirectory()
	return filepath.Join(configDirectory, strcase.ToKebab(name))
}

func LoadItems(client *onepassword.Client, forceReload bool) ([]onepassword.Item, error) {
	items, err := readConfig()

	if err != nil || forceReload {
		items, err = loadItemsFromSource(client)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}

	return items, err
}

func FindItem(items []onepassword.Item, title string) (onepassword.Item, error) {

	for _, item := range items {
		if item.Title == title {
			return item, nil
		}
	}

	return onepassword.Item{}, fmt.Errorf("configuration for '%s' not found", title)
}

func CleanData() {
	os.RemoveAll(getConfigDirectory())
}

func getConfigDirectory() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(home, ".config", "onekube")
}

func getConfigFilePath() string {
	configDirectory := getConfigDirectory()
	return filepath.Join(configDirectory, "configs")
}

func readConfig() ([]onepassword.Item, error) {
	itemsJson, err := os.ReadFile(getConfigFilePath())

	if err != nil {
		return nil, err
	}

	var items []onepassword.Item

	err = json.Unmarshal(itemsJson, &items)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return items, nil
}

func writeConfig(items []onepassword.Item) error {
	itemsJson, _ := json.Marshal(items)

	err := os.WriteFile(getConfigFilePath(), []byte(itemsJson), 0644)

	if err != nil {
		log.Fatal(err)
	}

	return err
}

func ensureConfigDirectoryExists() error {
	configDirectory := getConfigDirectory()

	err := os.MkdirAll(configDirectory, os.ModePerm)

	if err != nil {
		log.Fatal(err)
	}

	return err
}

func loadItemsFromSource(client *onepassword.Client) ([]onepassword.Item, error) {
	items, _ := client.Items("kubeconfig")

	err := ensureConfigDirectoryExists()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = writeConfig(items)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return items, nil
}
