package tasks

import "github.com/yandzee/wait-action/internal/utils"

type WaitTask struct {
	Workflows []string
}

func Parse(task string) ([]WaitTask, error) {
	workflows := utils.SplitStrings(task, ",")

	return []WaitTask{
		{
			Workflows: workflows,
		},
	}, nil
}
