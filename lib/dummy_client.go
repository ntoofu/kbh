package lib

import (
	"fmt"
	"sort"
	"time"
)

var DummyClientIssues map[string][]*Issue
var DummyClientIssueIdCounter map[string]uint

type DummyClient struct {
}

func initIssueListIfNecessary(clientId string) {
	if _, good := DummyClientIssueIdCounter[clientId]; !good {
		DummyClientIssues[clientId] = []*Issue{}
	}
	if _, good := DummyClientIssueIdCounter[clientId]; !good {
		DummyClientIssueIdCounter[clientId] = 0
	}
}

func (c DummyClient) CreateIssue(board *Board, draft *Issue) (*Issue, error) {
	initIssueListIfNecessary(board.Name)
	id := fmt.Sprintf("id%d", DummyClientIssueIdCounter[board.Name])
	uri := fmt.Sprintf("https://dummy/%s/%s", board.Name, id)
	issue := NewIssue(
		id,
		c,
		draft.Title,
		draft.Description,
		draft.Asignee,
		draft.Label,
		draft.IsClosed,
		time.Now(),
		uri)
	DummyClientIssues[board.Name] = append(DummyClientIssues[board.Name], issue)
	DummyClientIssueIdCounter[board.Name]++
	return issue, nil
}

func (c DummyClient) UpdateIssue(board *Board, toBeUpdated *Issue) error {
	initIssueListIfNecessary(board.Name)
	for i, issue := range DummyClientIssues[board.Name] {
		if issue.Id() == toBeUpdated.Id() {
			DummyClientIssues[board.Name][i] = toBeUpdated
			return nil
		}
	}
	return fmt.Errorf("No corresponding issue has found")
}

func (c DummyClient) ReadIssue(board *Board, issueId string) (*Issue, error) {
	initIssueListIfNecessary(board.Name)
	for _, issue := range DummyClientIssues[board.Name] {
		if issue.Id() == issueId {
			return issue, nil
		}
	}
	return &Issue{}, fmt.Errorf("No corresponding issue has found")
}

func (c DummyClient) QueryIssue(board *Board, condition StateCondDef) ([]*Issue, error) {
	foundIssues := make([]*Issue, 0)
	for _, issue := range DummyClientIssues[board.Name] {
		if (!condition.Asignee.Valid || issue.Asignee == condition.Asignee.Value) &&
			(isIncludedByStrSlice(condition.Labels, issue.Label)) &&
			(!condition.IsClosed.Valid || issue.IsClosed == condition.IsClosed.Value) {
			foundIssues = append(foundIssues, issue)
		}
	}
	return foundIssues, nil
}

func isIncludedByStrSlice(inclusion []string, slice []string) bool {
	sorted := make(sort.StringSlice, len(slice))
	copy(sorted, slice)
	sorted.Sort()

	for _, x := range inclusion {
		idx := sorted.Search(x)
		if idx >= len(sorted) || sorted[idx] != x {
			return false
		}
	}
	return true
}
