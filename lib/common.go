package lib

import "time"

type Task struct {
	Board *Board
	IssueId string
	Title string
	Description string
	State string
	Uri string
}

type Board struct {
	Name string
	Alias []string
	Client KanbanApiClient
	StateMapping []StateCondition
}

type StateCondition struct {
	StateName string
	Condition StateCondDef
}

type Issue struct {
	id string
	client KanbanApiClient
	Title, Description string
	Asignee string
	Label []string
	IsClosed bool
	updateTime time.Time
	Uri string
}

func (i Issue) Id() string {
	return i.id
}

func (i Issue) UpdateTime() time.Time {
	return i.updateTime
}

func NewIssue(id string, client KanbanApiClient, title, description, asignee string, label []string, isClosed bool, updateTime time.Time, uri string) *Issue {
	return &Issue {id, client, title, description, asignee, label, isClosed, updateTime, uri}
}

func NewEmptyIssue() *Issue {
	return NewIssue("", nil, "", "", "", []string{}, false, time.Unix(0,0), "")
}

type KanbanApiClient interface {
	CreateIssue(board *Board, draft *Issue) (*Issue, error)
	UpdateIssue(board *Board, toBeUpdated *Issue) error
	ReadIssue(board *Board, issueId string) (*Issue, error)
	QueryIssue(board *Board, condition StateCondDef) ([]*Issue, error)
}
