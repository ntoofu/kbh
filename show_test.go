package main

import (
	"bytes"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/ntoofu/kbh/lib"
)

func TestShowOutputWithoutOrder(t *testing.T) {
	expectedExample := `team1/proj1:id0,title1,todo
team1/proj1:id1,title2,todo
team1/proj1:id2,title3,doing
team1/proj1:id3,title4,waiting
team1/proj1:id4,title5,closed
team1/proj2:id0,title1,todo
team1/proj2:id1,title2,todo
team1/proj2:id2,title3,doing
team1/proj2:id3,title4,waiting
team1/proj2:id4,title5,closed
individual:id0,title1,todo
individual:id1,title2,todo
individual:id2,title3,todo
individual:id3,title4,todo
individual:id4,title5,fin
individual:id6,title7,now
`
	expectedLines := strings.Split(expectedExample, "\n")
	sort.Strings(expectedLines)
	c := lib.DummyClient{}
	now := time.Now()
	var (
		about20DaysAgo time.Time = now.Add(time.Hour * 24 * -20)
		about5DaysAgo  time.Time = now.Add(time.Hour * 24 * -5)
		about2DaysAgo  time.Time = now.Add(time.Hour * 24 * -2)
	)
	lib.DummyClientIssues = map[string][]lib.Issue{
		"team1/proj1": []lib.Issue{
			*lib.NewIssue("id0", c, "title1", "", "my.name", []string{"label1", "label2"}, false, about5DaysAgo),
			*lib.NewIssue("id1", c, "title2", "", "my.name", []string{}, false, about5DaysAgo),
			*lib.NewIssue("id2", c, "title3", "", "my.name", []string{"doing-label", "label3"}, false, about5DaysAgo),
			*lib.NewIssue("id3", c, "title4", "", "my.name", []string{"label4", "waiting-label"}, false, about20DaysAgo),
			*lib.NewIssue("id4", c, "title5", "", "my.name", []string{"doing-label", "label5"}, true, about5DaysAgo),
			*lib.NewIssue("id5", c, "title6", "", "my.name", []string{"doing-label", "label6"}, true, about20DaysAgo),
			*lib.NewIssue("id6", c, "title7", "", "others.name", []string{"doing-label"}, false, about2DaysAgo),
		},
		"team1/proj2": []lib.Issue{
			*lib.NewIssue("id0", c, "title1", "", "my.name", []string{"label1", "label2"}, false, about5DaysAgo),
			*lib.NewIssue("id1", c, "title2", "", "my.name", []string{}, false, about5DaysAgo),
			*lib.NewIssue("id2", c, "title3", "", "my.name", []string{"doing-label", "label3"}, false, about5DaysAgo),
			*lib.NewIssue("id3", c, "title4", "", "my.name", []string{"label4", "waiting-label"}, false, about20DaysAgo),
			*lib.NewIssue("id4", c, "title5", "", "my.name", []string{"doing-label", "label5"}, true, about5DaysAgo),
			*lib.NewIssue("id5", c, "title6", "", "my.name", []string{"doing-label", "label6"}, true, about20DaysAgo),
			*lib.NewIssue("id6", c, "title7", "", "others.name", []string{"doing-label"}, false, about2DaysAgo),
		},
		"individual": []lib.Issue{
			*lib.NewIssue("id0", c, "title1", "", "my.name", []string{"label1", "label2"}, false, about5DaysAgo),
			*lib.NewIssue("id1", c, "title2", "", "my.name", []string{}, false, about5DaysAgo),
			*lib.NewIssue("id2", c, "title3", "", "my.name", []string{"doing-label", "label3"}, false, about5DaysAgo),
			*lib.NewIssue("id3", c, "title4", "", "my.name", []string{"label4", "waiting-label"}, false, about20DaysAgo),
			*lib.NewIssue("id4", c, "title5", "", "my.name", []string{"doing-label", "label5"}, true, about5DaysAgo),
			*lib.NewIssue("id5", c, "title6", "", "my.name", []string{"doing-label", "label6"}, true, about20DaysAgo),
			*lib.NewIssue("id6", c, "title7", "", "others.name", []string{"doing-label"}, false, about2DaysAgo),
		},
	}
	lib.DummyClientIssueIdCounter = map[string]uint{
		"team1/proj1": 6,
		"team1/proj2": 6,
		"individual":  6,
	}
	stdoutBuf := new(bytes.Buffer)
	testConf, err := lib.ParseConfig("config_test.yml")
	if err != nil {
		t.Errorf("An error has occured during lib.ParseConfig: %v", err)
	}
	err = showTask([]string{}, testConf, stdoutBuf)
	if err != nil {
		t.Errorf("An error has occured during running 'show' sub-command: %v", err)
	}
	outputLines := strings.Split(stdoutBuf.String(), "\n")
	sort.Strings(outputLines)
	if !reflect.DeepEqual(outputLines, expectedLines) {
		t.Errorf("Sorted output does not match against expected output. --- expected ---\n%v\n--- actual ---\n%v\n", expectedLines, outputLines)
	}
}
