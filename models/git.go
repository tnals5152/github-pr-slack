package models

type RepoInfo struct {
	Key     string `json:"-" yaml:"-"`
	Owner   string `json:"owner" yaml:"owner"`
	Repo    string `json:"repo" yaml:"repo"`
	Message struct {
		LastMessageDate string `json:"date" yaml:"date"`
		LastMessageID   string `json:"id" yaml:"id"`
	} `json:"message" yaml:"message"`
}

type PRInfo struct {
	Number             int32   `json:"number"`
	Url                string  `json:"url"`
	State              string  `json:"state"`
	User               User    `json:"user"`
	RequestedReviewers []User  `json:"requested_reviewers"`
	UpdatedAt          NewTime `json:"updated_at,omitempty"`
	CreatedAt          NewTime `json:"created_at,omitempty"`
	ClosedAt           NewTime `json:"closed_at,omitempty"`
	Links              struct {
		Html struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"_links"`
}

type User struct {
	Login string `json:"login"`
}

type PRReview struct {
	User  User   `json:"user"`
	State string `json:"state"`
}
