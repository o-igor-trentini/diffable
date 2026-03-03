package domain

import "time"

type AnalysisType string

const (
	AnalysisTypeSingleCommit AnalysisType = "single_commit"
	AnalysisTypeCommitRange  AnalysisType = "commit_range"
	AnalysisTypePullRequest  AnalysisType = "pull_request"
)

type AnalysisStatus string

const (
	AnalysisStatusPending    AnalysisStatus = "pending"
	AnalysisStatusProcessing AnalysisStatus = "processing"
	AnalysisStatusCompleted  AnalysisStatus = "completed"
	AnalysisStatusFailed     AnalysisStatus = "failed"
)

type Analysis struct {
	ID            string         `json:"id"`
	AnalysisType  AnalysisType   `json:"analysis_type"`
	Status        AnalysisStatus `json:"status"`
	Workspace     string         `json:"workspace,omitempty"`
	RepoSlug      string         `json:"repo_slug,omitempty"`
	CommitHash    string         `json:"commit_hash,omitempty"`
	FromHash      string         `json:"from_hash,omitempty"`
	ToHash        string         `json:"to_hash,omitempty"`
	PrID          *int           `json:"pr_id,omitempty"`
	RawDiff       string         `json:"-"`
	DiffHash      string         `json:"diff_hash,omitempty"`
	GeneratedDesc string         `json:"generated_description,omitempty"`
	ModelUsed     string         `json:"model_used,omitempty"`
	TokensUsed    *int           `json:"tokens_used,omitempty"`
	ErrorMessage  string         `json:"error_message,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type Refinement struct {
	ID           string    `json:"id"`
	AnalysisID   string    `json:"analysis_id"`
	Instruction  string    `json:"instruction"`
	OriginalDesc string    `json:"original_description"`
	RefinedDesc  string    `json:"refined_description,omitempty"`
	ModelUsed    string    `json:"model_used,omitempty"`
	TokensUsed   *int      `json:"tokens_used,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
