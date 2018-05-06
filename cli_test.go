package main

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/urfave/cli"
)

func TestHelpMessage(t *testing.T) {
	expectedHelpMsg := []byte(
		`NAME:
   Kanban Boards Handler - handle many kanban boards or issue trackers in a persistent way

USAGE:
   kbh.test [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     create, c  create a new task
     update, u  update state of a task
     show, s    display a list of tasks
     help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -c value  config file (default: "config.yml")
   --help, -h                show help
   --version, -v             print the version
`)
	argsList := [][]string{
		[]string{"./kbh", "help"},
		[]string{"./kbh", "-c", "config_test.yml", "h"},
		[]string{"./kbh", "--help", "-c", "config_test.yml"},
		[]string{"./kbh", "-h"},
	}

	for _, args := range argsList {
		stdoutBuf := new(bytes.Buffer)
		stderrBuf := new(bytes.Buffer)
		app := generateCliApp(stdoutBuf, stderrBuf, dummyExiter(t), normalSubcommandsFactory{stdoutBuf})
		app.Run(args)
		if bytes.Compare(stdoutBuf.Bytes(), expectedHelpMsg) != 0 {
			t.Errorf("Output differs from the expected help message\n--- expected ---\n%s\n--- actual ---\n%s\n", expectedHelpMsg, stdoutBuf.Bytes())
		}
	}
}

func TestInvalidArguments(t *testing.T) {
	argsList := [][]string{
		[]string{"./kbh", "invalid_cmd"},
		[]string{"./kbh", "-c", ".not_existing_file", "show"},
	}

	for _, args := range argsList {
		stdoutBuf := new(bytes.Buffer)
		stderrBuf := new(bytes.Buffer)
		app := generateCliApp(stdoutBuf, stderrBuf, func(_ int) { panic("CLI exit with errors") }, normalSubcommandsFactory{stdoutBuf})
		panicTest(t, func() {
			err := app.Run(args)
			if err != nil {
				panic("CLI exit with errors")
			}
		}, args)
	}
}

type callHistory struct {
	Called bool
	Args   []string
}

type mockSubcommandsFactory struct {
	CreateTaskHistory callHistory
	UpdateTaskHistory callHistory
	ShowTaskHistory   callHistory
}

func (f *mockSubcommandsFactory) GenerateCreateTask() func(*cli.Context) error {
	return func(c *cli.Context) error {
		f.CreateTaskHistory = callHistory{true, c.Args()}
		return nil
	}
}
func (f *mockSubcommandsFactory) GenerateUpdateTask() func(*cli.Context) error {
	return func(c *cli.Context) error {
		f.UpdateTaskHistory = callHistory{true, c.Args()}
		return nil
	}
}
func (f *mockSubcommandsFactory) GenerateShowTask() func(*cli.Context) error {
	return func(c *cli.Context) error {
		f.ShowTaskHistory = callHistory{true, c.Args()}
		return nil
	}
}

func TestSubcommandCall(t *testing.T) {
	type argsAndExpectation struct {
		Args   []string
		Expect mockSubcommandsFactory
	}

	argsAndExpectations := []argsAndExpectation{
		argsAndExpectation{
			[]string{"./kbh", "-c", "config_test.yml", "create", "board", "title", "state"},
			mockSubcommandsFactory{
				callHistory{true, []string{"board", "title", "state"}},
				callHistory{false, nil},
				callHistory{false, nil},
			},
		},
		argsAndExpectation{
			[]string{"./kbh", "-c", "config_test.yml", "update", "id", "state"},
			mockSubcommandsFactory{
				callHistory{false, nil},
				callHistory{true, []string{"id", "state"}},
				callHistory{false, nil},
			},
		},
		argsAndExpectation{
			[]string{"./kbh", "-c", "config_test.yml", "show"},
			mockSubcommandsFactory{
				callHistory{false, nil},
				callHistory{false, nil},
				callHistory{true, nil},
			},
		},
	}

	for _, ae := range argsAndExpectations {
		stdoutBuf := new(bytes.Buffer)
		stderrBuf := new(bytes.Buffer)
		mock := &mockSubcommandsFactory{}
		app := generateCliApp(stdoutBuf, stderrBuf, dummyExiter(t), mock)
		err := app.Run(ae.Args)
		if err != nil {
			t.Errorf("An error has occured during running CLI")
		}
		if !reflect.DeepEqual(*mock, ae.Expect) {
			t.Errorf("Calling history of subcommands is different from that expected (%v)", ae.Args)
		}
	}
}

func panicTest(t *testing.T, f func(), info interface{}) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Function runs without errors against expectations (%v)", info)
		}
	}()
	f()
}

func dummyExiter(t *testing.T) func(int) {
	return func(i int) {
		t.Logf("CLI exit with error code %d", i)
		t.FailNow()
	}
}
