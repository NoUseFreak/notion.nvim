package cli

type IssueDBSpec struct {
	ID        string
	Title     string
	Assignees string
	Status    string

	IDPrefix string
}

type Issue struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Assignees []string `json:"assignees"`
	URL       string   `json:"url"`
	Content   []string `json:"content,omitempty"`
}
