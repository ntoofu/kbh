package main

import (
	"fmt"
	"io"
	"sort"
	"strings"

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

	errorChan := make(chan error)
	taskChan := make(chan *lib.Task)
	for _, bd := range boardList {
		go func(b *lib.Board, tch chan *lib.Task, ech chan error) {
			issues, err := getTaskCandidates(b)
			if err != nil {
				ech <- err
			}
			for _, issue := range issues {
				task, good := issueToTask(b, issue)
				if good {
					taskChan <- task
				}
			}
			ech <- nil
		}(bd, taskChan, errorChan)
	}

	taskList := make([]*lib.Task, 0)
	for i := len(boardList); i>0; {
		select {
		case err := <-errorChan:
			i--
			if err != nil {
				return err
			}
		case task := <-taskChan:
			taskList = append(taskList, task)
		}
	}

	showOpts := conf.Command.Show
	for _, task := range taskList {
		fieldStr := map[string]string{
			"state": task.State,
			"title": task.Title,
			"uri": task.Uri,
		}
		for i, field := range showOpts.Field {
			if i != 0 {
				fmt.Fprintf(stdout, "%s", showOpts.Delimiter)
			}
			replacedStr := strings.Replace(fieldStr[field], showOpts.Delimiter, showOpts.Replacer, -1)
			fmt.Fprintf(stdout, "%s", replacedStr)
		}
		fmt.Fprintf(stdout, "\n")
	}

	return nil
}

func getTaskCandidates(bd *lib.Board) ([]*lib.Issue, error) {
	errorChan := make(chan error)
	issuesChan := make(chan []*lib.Issue)
	defer func() {
		close(errorChan)
		close(issuesChan)
	}()

	issues := make([]*lib.Issue, 0)
	for _, cond := range bd.StateMapping {
		go func(b *lib.Board, c lib.StateCondDef, ich chan []*lib.Issue, ech chan error) {
			partialIssues, err := b.Client.QueryIssue(b, c)
			if err != nil {
				ech <- err
			}
			ech <- nil
			ich <- partialIssues
		}(bd, cond.Condition, issuesChan, errorChan)
	}

	dupCounter := make(map[string]struct{})
	for i:=0; i<len(bd.StateMapping); i++ {
		err := <-errorChan
		if err != nil {
			return []*lib.Issue{}, err
		}
		for _, is := range <-issuesChan {
			_, existence := dupCounter[is.Id()]
			if !existence {
				issues = append(issues, is)
				dupCounter[is.Id()] = struct{}{}
			}
		}
	}
	return issues, nil
}

func getBoardList(endpointDef []lib.EndpointDef, boardDef []lib.BoardDef) ([]*lib.Board, error) {
	apiClientList := map[string]lib.KanbanApiClient{}
	for _, endpoint := range endpointDef {
		var apiClient lib.KanbanApiClient
		switch endpoint.Type {
		case "gitlab":
			apiClient = lib.GitlabClient{endpoint.Url, endpoint.AuthToken}
		case "dummy":
			apiClient = lib.DummyClient{}
		}
		apiClientList[endpoint.Name] = apiClient
	}

	boardList := make([]*lib.Board, 0)
	for _, bd := range boardDef {
		stateMapping := reorderStateConditions(bd.Mapping.State)
		elem := lib.Board{bd.Name, bd.Alias, apiClientList[bd.Endpoint], stateMapping}
		boardList = append(boardList, &elem)
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

func issueToTask(bd *lib.Board, issue *lib.Issue) (*lib.Task, bool) {
	state := ""
	for _, cond := range bd.StateMapping {
		if cond.Condition.IsMatched(issue) {
			state = cond.StateName
			break
		}
	}
	if state == "" {
		return nil, false
	}
	return &lib.Task{bd, issue.Id(), issue.Title, issue.Description, state, issue.Uri}, true
}
