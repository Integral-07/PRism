package webhook

type PullRequestEvent struct {
	Action      string       `json:"action"`
	PullRequest PullRequest  `json:"pull_request"`
	Repository  Repository   `json:"repository"`
	Installation Installation `json:"installation"`
}

type PullRequest struct {
	Number  int    `json:"number"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	HTMLURL string `json:"html_url"`
}

type Repository struct {
	FullName string `json:"full_name"`
}

type Installation struct {
	ID int64 `json:"id"`
}
