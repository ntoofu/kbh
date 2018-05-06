package lib

type GitlabClient struct {
	Url, ApiVersion string
}

func (c GitlabClient) CreateIssue(boardId string, draft Issue) (Issue, error) {
	// TODO
	return Issue{}, nil
}

func (c GitlabClient) UpdateIssue(boardId string, toBeUpdated Issue) error {
	// TODO
	return nil
}

func (c GitlabClient) ReadIssue(boardId string, issueId string) (Issue, error) {
	// TODO
	return Issue{}, nil
}

func (c GitlabClient) QueryIssue(boardId string, condition StateCondDef) ([]Issue, error) {
	// TODO
	return []Issue{}, nil
}
