package main

import (
	"fmt"
	"io"
	"sort"

	"github.com/ntoofu/kbh/lib"
	"github.com/urfave/cli"
)

func (f normalSubcommandsFactory) GenerateShowTask() func(*cli.Context) error {
	return func(c *cli.Context) error {
		conf, err := lib.ParseConfig(c.GlobalString("config"))
		if err != nil {
			return err
		}
		return showTask(c.Args(), conf, f.stdout)
	}
}

func showTask(args []string, conf lib.GlobalConfig, stdout io.Writer) error {
	if len(args) > 0 {
		return fmt.Errorf("too many arguments are given: %v", args[0:])
	}

	boardList, err := getBoardList(conf.Endpoint, conf.Board)
	if err != nil {
		return err
	}

	taskList := make([]lib.Task, 0)
	for _, bd := range boardList {
		dupCounter := make(map[string]struct{})
		issues := make([]lib.Issue, 0)
		for _, cond := range bd.StateMapping {
			partialIssues, err := bd.Client.QueryIssue(bd.Name, cond.Condition)
			if err != nil {
				return err
			}
			for _, is := range partialIssues {
				_, existence := dupCounter[is.Id()]
				if !existence {
					issues = append(issues, is)
					dupCounter[is.Id()] = struct{}{}
				}
			}
		}

		for _, issue := range issues {
			task, good := issueToTask(bd, issue)
			if good {
				taskList = append(taskList, task)
			}
		}
	}

	for _, task := range taskList {
		fmt.Fprintf(stdout, "%s:%s,%s,%s\n", task.Board.Name, task.IssueId, task.Title, task.State)
	}

	return nil
}

func getBoardList(endpointDef []lib.EndpointDef, boardDef []lib.BoardDef) ([]lib.Board, error) {
	apiClientList := map[string]lib.KanbanApiClient{}
	for _, endpoint := range endpointDef {
		var apiClient lib.KanbanApiClient
		switch endpoint.Type {
		case "gitlab":
			apiClient = lib.GitlabClient{endpoint.Url, endpoint.ApiVersion}
		case "dummy":
			apiClient = lib.DummyClient{}
		}
		apiClientList[endpoint.Name] = apiClient
	}

	boardList := make([]lib.Board, 0)
	for _, bd := range boardDef {
		stateMapping := reorderStateConditions(bd.Mapping.State)
		elem := lib.Board{bd.Name, bd.Alias, apiClientList[bd.Endpoint], stateMapping}
		boardList = append(boardList, elem)
	}

	return boardList, nil
}

func reorderStateConditions(orgConditions map[string]lib.StateCondDef) []lib.StateCondition {
	conds := make([]lib.StateCondition, 0)
	for state, cond := range orgConditions {
		conds = append(conds, lib.StateCondition{state, cond})
	}
	sort.Slice(conds, func(i, j int) bool { return conds[i].Condition.Order < conds[j].Condition.Order })
	return conds
}

func issueToTask(bd lib.Board, issue lib.Issue) (lib.Task, bool) {
	state := ""
	for _, cond := range bd.StateMapping {
		if cond.Condition.IsMatched(issue) {
			state = cond.StateName
			break
		}
	}
	if state == "" {
		return lib.Task{}, false
	}
	return lib.Task{bd, issue.Id(), issue.Title, issue.Description, state}, true
}
