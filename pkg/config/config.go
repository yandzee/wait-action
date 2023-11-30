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
	Head        github.CommitSpec
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

	headSha := os.Getenv("INPUT_HEAD_SHA")
	if len(headSha) == 0 {
		return nil, fmt.Errorf("INPUT_HEAD_SHA is not set")
	}

	head := github.CommitSpec{
		Sha: headSha,
	}

	workflows := os.Getenv("INPUT_WORKFLOWS")

	return &Config{
		GithubToken: ghToken,
		PollDelay:   pollDelay,
		RepoOwner:   parts[0],
		Repo:        parts[1],
		Head:        head,
		Workflows:   workflows,
	}, nil
}

func (c *Config) LogAttrs() []any {
	return []any{
		slog.Bool("token-is-set", len(c.GithubToken) > 0),
		slog.String("poll-delay", c.PollDelay.String()),
		slog.String("head-sha", c.Head.Sha),
		slog.String("workflows", c.Workflows),
	}
}
