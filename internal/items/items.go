package items

import (
	"fmt"
	"log"

	"github.com/eddymoulton/onekube/internal/config"
	"github.com/eddymoulton/onekube/internal/onepassword"
)

func Load(client *onepassword.Client, forceReload bool) ([]onepassword.Item, error) {
	items, err := config.Read()

	if err != nil || forceReload {
		items, err = loadItemsFromSource(client)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}

	return items, err
}

func Find(items []onepassword.Item, title string) (onepassword.Item, error) {

	for _, item := range items {
		if item.Title == title {
			return item, nil
		}
	}

	return onepassword.Item{}, fmt.Errorf("configuration for '%s' not found", title)
}

func loadItemsFromSource(client *onepassword.Client) ([]onepassword.Item, error) {
	items, _ := client.Items("kubeconfig")

	err := config.EnsureDirectoryExists()

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = config.Write(items)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return items, nil
}
