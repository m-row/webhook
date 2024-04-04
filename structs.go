package main

type Payload struct {
	CallbackURL string     `json:"callback_url"`
	PushData    PushData   `json:"push_data"`
	Repository  Repository `json:"repository"`
}

type PushData struct {
	Images    []any  `json:"images"`
	MediaType string `json:"media_type"`
	PushedAt  int    `json:"pushed_at"`
	Pusher    string `json:"pusher"`
	Tag       string `json:"tag"`
}

type Repository struct {
	DateCreated     int    `json:"date_created"`
	Description     string `json:"description"`
	FullDescription any    `json:"full_description"`
	IsOfficial      bool   `json:"is_official"`
	IsPrivate       bool   `json:"is_private"`
	IsTrusted       bool   `json:"is_trusted"`
	Name            string `json:"name"`
	Namespace       string `json:"namespace"`
	Owner           string `json:"owner"`
	RepoName        string `json:"repo_name"`
	RepoURL         string `json:"repo_url"`
	StarCount       int    `json:"star_count"`
	Status          string `json:"status"`
}
