package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	GithubToken string
	PollDelay   time.Duration
	RepoOwner   string
	Repo        string
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

	return &Config{
		GithubToken: ghToken,
		PollDelay:   pollDelay,
		RepoOwner:   parts[0],
		Repo:        parts[1],
	}, nil
}
