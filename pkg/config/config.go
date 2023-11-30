package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/yandzee/wait-action/pkg/github"
)

type Config struct {
	GithubToken string
	PollDelay   time.Duration
	RepoOwner   string
	Repo        string
	HeadRef     string
}

func ParseEnv() (*Config, error) {
	ghToken := os.Getenv("GITHUB_TOKEN")

	pollDelay, err := time.ParseDuration(os.Getenv("INPUT_POLL_DELAY"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse poll delay: %s", err.Error())
	}

	repo := os.Getenv("GITHUB_REPOSITORY")
	parts := strings.SplitN(repo, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("failed to parse repo '%s' into owner and name", repo)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse repository: %s", err.Error())
	}

	headRef := os.Getenv("GITHUB_HEAD_REF")
	if len(headRef) == 0 {
		return nil, fmt.Errorf("GITHUB_HEAD_REF is not set")
	}

	return &Config{
		GithubToken: ghToken,
		PollDelay:   pollDelay,
		RepoOwner:   parts[0],
		Repo:        parts[1],
		HeadRef:     headRef,
	}, nil
}

func (c *Config) CommitSpec() github.CommitSpec {
	// TODO: How to do it right?
	return github.CommitSpec{
		Sha:    "",
		Branch: c.HeadRef,
	}
}
