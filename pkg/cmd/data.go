package cmd

type Teams struct {
	Teams []Team `json:"teams"`
}

type Team struct {
	Name       string     `json:"name"`
	Categories []Category `json:"categories"`
}

type Category struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Sorting  string   `json:"sorting"`
	Channels []string `json:"channels"`
}
