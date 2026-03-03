package bitbucket

import "time"

type Author struct {
	DisplayName string `json:"display_name"`
	UUID        string `json:"uuid"`
}

type CommitAuthor struct {
	Raw  string `json:"raw"`
	User Author `json:"user"`
}

type ParentCommit struct {
	Hash string `json:"hash"`
}

type Commit struct {
	Hash    string       `json:"hash"`
	Message string       `json:"message"`
	Date    time.Time    `json:"date"`
	Author  CommitAuthor `json:"author"`
	Parents []ParentCommit `json:"parents"`
}

type Branch struct {
	Name string `json:"name"`
}

type PRRef struct {
	Branch     Branch     `json:"branch"`
	Repository Repository `json:"repository"`
}

type PullRequest struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	State       string    `json:"state"`
	Source      PRRef     `json:"source"`
	Destination PRRef     `json:"destination"`
	Author      Author    `json:"author"`
	CreatedOn   time.Time `json:"created_on"`
	UpdatedOn   time.Time `json:"updated_on"`
}

type DiffstatPath struct {
	Old string `json:"old"`
	New string `json:"new"`
}

type DiffstatEntry struct {
	Status       string       `json:"status"`
	LinesAdded   int          `json:"lines_added"`
	LinesRemoved int          `json:"lines_removed"`
	Old          *DiffstatFile `json:"old"`
	New          *DiffstatFile `json:"new"`
}

type DiffstatFile struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

type Repository struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
}

type PaginatedResponse[T any] struct {
	Values  []T    `json:"values"`
	Next    string `json:"next"`
	PageLen int    `json:"pagelen"`
	Size    int    `json:"size"`
	Page    int    `json:"page"`
}
