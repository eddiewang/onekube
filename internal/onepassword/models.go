package onepassword

import "time"

type Client struct{}

type Vault struct {
	ContentVersion int    `json:"content_version"`
	ID             string `json:"id"`
	Name           string `json:"name"`
}

type Item struct {
	AdditionalInformation string    `json:"additional_information,omitempty"`
	Category              string    `json:"category"`
	CreatedAt             time.Time `json:"created_at"`
	Favorite              bool      `json:"favorite,omitempty"`
	ID                    string    `json:"id"`
	LastEditedBy          string    `json:"last_edited_by"`
	Tags                  []string  `json:"tags,omitempty"`
	Title                 string    `json:"title"`
	UpdatedAt             time.Time `json:"updated_at"`
	Urls                  []struct {
		Href    string `json:"href"`
		Label   string `json:"label,omitempty"`
		Primary bool   `json:"primary,omitempty"`
	} `json:"urls,omitempty"`
	Vault struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"vault"`
	Version int `json:"version"`
}
