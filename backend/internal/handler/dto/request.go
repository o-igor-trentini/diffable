package dto

import (
	"fmt"

	"github.com/igor-trentini/diffable/backend/internal/domain"
)

var validLevels = map[string]bool{
	"technical":  true,
	"functional": true,
	"executive":  true,
}

func validateLevel(level string) error {
	if level != "" && !validLevels[level] {
		return fmt.Errorf("%w: level must be one of: technical, functional, executive", domain.ErrValidation)
	}
	return nil
}

type AnalyzeCommitRequest struct {
	Workspace  string `json:"workspace"`
	RepoSlug   string `json:"repo_slug"`
	CommitHash string `json:"commit_hash"`
	RawDiff    string `json:"raw_diff"`
	Level      string `json:"level"`
}

func (r *AnalyzeCommitRequest) Validate() error {
	hasHash := r.CommitHash != ""
	hasRawDiff := r.RawDiff != ""

	if !hasHash && !hasRawDiff {
		return fmt.Errorf("%w: commit_hash (with workspace and repo_slug) or raw_diff is required", domain.ErrValidation)
	}

	if hasHash && (r.Workspace == "" || r.RepoSlug == "") {
		return fmt.Errorf("%w: workspace and repo_slug are required when using commit_hash", domain.ErrValidation)
	}

	if err := validateLevel(r.Level); err != nil {
		return err
	}

	return nil
}

type AnalyzeRangeRequest struct {
	Workspace string `json:"workspace"`
	RepoSlug  string `json:"repo_slug"`
	FromHash  string `json:"from_hash"`
	ToHash    string `json:"to_hash"`
	Level     string `json:"level"`
}

func (r *AnalyzeRangeRequest) Validate() error {
	if r.Workspace == "" {
		return fmt.Errorf("%w: workspace is required", domain.ErrValidation)
	}
	if r.RepoSlug == "" {
		return fmt.Errorf("%w: repo_slug is required", domain.ErrValidation)
	}
	if r.FromHash == "" {
		return fmt.Errorf("%w: from_hash is required", domain.ErrValidation)
	}
	if r.ToHash == "" {
		return fmt.Errorf("%w: to_hash is required", domain.ErrValidation)
	}
	if err := validateLevel(r.Level); err != nil {
		return err
	}
	return nil
}

type AnalyzePRRequest struct {
	Workspace     string `json:"workspace"`
	RepoSlug      string `json:"repo_slug"`
	PRID          int    `json:"pr_id"`
	RawDiff       string `json:"raw_diff"`
	PRTitle       string `json:"pr_title"`
	PRDescription string `json:"pr_description"`
	Level         string `json:"level"`
}

func (r *AnalyzePRRequest) Validate() error {
	hasPRID := r.PRID > 0
	hasRawDiff := r.RawDiff != ""

	if !hasPRID && !hasRawDiff {
		return fmt.Errorf("%w: pr_id (with workspace and repo_slug) or raw_diff (with pr_title) is required", domain.ErrValidation)
	}

	if hasPRID && (r.Workspace == "" || r.RepoSlug == "") {
		return fmt.Errorf("%w: workspace and repo_slug are required when using pr_id", domain.ErrValidation)
	}

	if hasRawDiff && !hasPRID && r.PRTitle == "" {
		return fmt.Errorf("%w: pr_title is required when using raw_diff without pr_id", domain.ErrValidation)
	}

	if err := validateLevel(r.Level); err != nil {
		return err
	}

	return nil
}

type RefineRequest struct {
	Instruction string `json:"instruction"`
}

func (r *RefineRequest) Validate() error {
	if r.Instruction == "" {
		return fmt.Errorf("%w: instruction is required", domain.ErrValidation)
	}
	return nil
}
