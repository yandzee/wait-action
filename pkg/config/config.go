package config

import (
	"fmt"
	"log/slog"
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
	Workflows   string
}

func ParseEnv() (*Config, error) {
	ghToken := os.Getenv("GITHUB_TOKEN")

	pollDelayStr := strings.TrimSpace(os.Getenv("INPUT_POLL_DELAY"))
	if len(pollDelayStr) == 0 {
		pollDelayStr = "10s"
	}

	pollDelay, err := time.ParseDuration(pollDelayStr)
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

	workflows := os.Getenv("INPUT_WORKFLOWS")

	return &Config{
		GithubToken: ghToken,
		PollDelay:   pollDelay,
		RepoOwner:   parts[0],
		Repo:        parts[1],
		HeadRef:     headRef,
		Workflows:   workflows,
	}, nil
}

func (c *Config) LogAttrs() []any {
	return []any{
		slog.Bool("token-is-set", len(c.GithubToken) > 0),
		slog.String("poll-delay", c.PollDelay.String()),
		slog.String("head-ref", c.HeadRef),
		slog.String("workflows", c.Workflows),
	}
}

func (c *Config) CommitSpec() github.CommitSpec {
	// TODO: How to do it right?
	return github.CommitSpec{
		Sha:    "",
		Branch: c.HeadRef,
	}
}
