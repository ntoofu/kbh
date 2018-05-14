package lib

import (
	"fmt"
	"time"

	"github.com/xanzy/go-gitlab"
)

type GitlabClient struct {
	Url string
	AuthToken string
}

func (c GitlabClient) initClient() *gitlab.Client {
	client := gitlab.NewClient(nil, c.AuthToken)
	client.SetBaseURL(c.Url)
	return client
}

func (c GitlabClient) CreateIssue(board *Board, draft *Issue) (*Issue, error) {
	return nil, nil
}

func (c GitlabClient) UpdateIssue(board *Board, toBeUpdated *Issue) error {
	return nil
}

func (c GitlabClient) ReadIssue(board *Board, issueId string) (*Issue, error) {
	return nil, nil
}

func (c GitlabClient) QueryIssue(board *Board, condition StateCondDef) ([]*Issue, error) {
	client := c.initClient()
	var issueOpt gitlab.ListProjectIssuesOptions
	if condition.LabelName.Valid {
		issueOpt.Labels = []string{condition.LabelName.Value}
	}
	if condition.Asignee.Valid {
		// TODO: cache user ID
		var usersOpt gitlab.ListUsersOptions
		usersOpt.Username = &condition.Asignee.Value
		users, _, err := client.Users.ListUsers(&usersOpt)
		if err != nil {
			return nil, err
		}
		if len(users) != 1 {
			return nil, fmt.Errorf("Found %d users whose name is %s", len(users), condition.Asignee.Value)
		}
		userid := users[0].ID
		issueOpt.AssigneeID = &userid
	}
	if condition.IsClosed.Valid {
		if condition.IsClosed.Value {
			issueStatusClosed := "closed"
			issueOpt.State = &issueStatusClosed
		} else {
			issueStatusOpened := "opened"
			issueOpt.State = &issueStatusOpened
		}
	}
	issueScopeAll := "all"
	issueOpt.Scope = &issueScopeAll

	gitlabIssues, _, err := client.Issues.ListProjectIssues(board.Name, &issueOpt)
	if err != nil {
		return []*Issue{}, err
	}

	now := time.Now()
	issues := make([]*Issue, 0)
	for _, gitlabIssue := range gitlabIssues {
		if condition.MaxDaysWOUpdate.Valid {
			// Because GitLab API client cannot filter issues by updated date,
			// issues must be filtered here.
			// TODO: use GitLab API client's filtering option for query, if that feature is implemented
			maxHours := time.Duration(condition.MaxDaysWOUpdate.Value) * 24 * time.Hour
			if now.After(gitlabIssue.UpdatedAt.Add(maxHours)) {
				continue
			}
		}
		issue := NewIssue(
			fmt.Sprintf("%d", gitlabIssue.IID),
			board.Client,
			gitlabIssue.Title,
			gitlabIssue.Description,
			gitlabIssue.Assignee.Name,
			gitlabIssue.Labels,
			gitlabIssue.State == "closed",
			*gitlabIssue.UpdatedAt)
		issues = append(issues, issue)
	}
	return issues, nil
}
