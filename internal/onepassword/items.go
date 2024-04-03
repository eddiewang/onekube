package onepassword

import (
	"fmt"
	"strings"
)

func (c *Client) Item(itemIDOrName string) (*Item, error) {
	var out Item
	err := c.runOpAndUnmarshal("item", []string{"get", itemIDOrName}, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *Client) Items(tag string) ([]Item, error) {
	var out []Item
	err := c.runOpAndUnmarshal("item", []string{"list", fmt.Sprintf("--tags=%s", tag)}, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) Read(lookupIdentifier string) (string, error) {
	out, err := c.runOp("read", []string{lookupIdentifier})
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func (c *Client) ReadItemField(vaultIdOrName string, itemIdOrName string, fieldName string) (string, error) {
	lookupString := fmt.Sprintf("op://%s/%s/%s", vaultIdOrName, itemIdOrName, fieldName)
	return c.Read(lookupString)
}
