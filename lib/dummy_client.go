package lib

import (
	"fmt"
	"time"
)

var DummyClientIssues map[string][]*Issue
var DummyClientIssueIdCounter map[string]uint

type DummyClient struct {
}

func initIssueListIfNecessary (clientId string) {
	if _, good := DummyClientIssueIdCounter[clientId]; !good {
		DummyClientIssues[clientId] = []*Issue{}
	}
	if _, good := DummyClientIssueIdCounter[clientId]; !good {
		DummyClientIssueIdCounter[clientId] = 0
	}
}

func (c DummyClient) CreateIssue(boardId string, draft *Issue) (*Issue, error) {
	initIssueListIfNecessary(boardId)
	issue := NewIssue(
		fmt.Sprintf("id%d", DummyClientIssueIdCounter[boardId]),
		c,
		draft.Title,
		draft.Description,
		draft.Asignee,
		draft.Label,
		draft.IsClosed,
		time.Now())
	DummyClientIssues[boardId] = append(DummyClientIssues[boardId], issue)
	DummyClientIssueIdCounter[boardId]++
	return issue, nil
}

func (c DummyClient) UpdateIssue(boardId string, toBeUpdated *Issue) error {
	initIssueListIfNecessary(boardId)
	for i, issue := range DummyClientIssues[boardId] {
		if issue.Id() == toBeUpdated.Id() {
			DummyClientIssues[boardId][i] = toBeUpdated
			return nil
		}
	}
	return fmt.Errorf("No corresponding issue has found")
}

func (c DummyClient) ReadIssue(boardId string, issueId string) (*Issue, error) {
	initIssueListIfNecessary(boardId)
	for _, issue := range DummyClientIssues[boardId] {
		if issue.Id() == issueId {
			return issue, nil
		}
	}
	return &Issue{}, fmt.Errorf("No corresponding issue has found")
}

func (c DummyClient) QueryIssue(boardId string, condition StateCondDef) ([]*Issue, error) {
	foundIssues := make([]*Issue, 0)
	for _, issue := range DummyClientIssues[boardId] {
		if ( (!condition.Asignee.Valid || issue.Asignee == condition.Asignee.Value) &&
			 (!condition.LabelName.Valid || isInStrSlice(condition.LabelName.Value, issue.Label)) &&
			 (!condition.IsClosed.Valid || issue.IsClosed == condition.IsClosed.Value)) {
			foundIssues = append(foundIssues, issue)
		}
	}
	return foundIssues, nil
}

func isInStrSlice(elem string, slice []string) bool {
	for _, x := range slice {
		if x == elem {
			return true
		}
	}
	return false
}
