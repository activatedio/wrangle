package context

import (
	"os/exec"

	"os"

	"errors"

	"fmt"

	"path/filepath"

	"regexp"

	"github.com/activatedio/wrangle/config"
)

type Context struct {
	Delegate   *exec.Cmd
	Executable *config.Executable
	Variables  map[string]string
}

var getArgs = func() []string {
	return os.Args
}

func NewContext(c *config.Config) (*Context, error) {

	rawArgs := getArgs()

	args, vars := getArgAndVars(rawArgs)

	// Need to pre-process the args for any suffixed with "wr"

	if len(args) < 2 {
		return nil, errors.New("Syntax: wrangle [wrangle options] [program name] [program args]")
	}

	executable := args[1]
	executablePath, err := exec.LookPath(executable)
	executableBase := filepath.Base(executable)

	if err != nil {
		return nil, err
	}

	executableConfig, ok := c.Executables[executableBase]

	if !ok {
		return nil, errors.New(fmt.Sprintf("Executable [%s] not found in config.", executableBase))
	}

	wd, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	return &Context{
		Variables: vars,
		Delegate: &exec.Cmd{
			Path:   executablePath,
			Args:   args[1:],
			Env:    os.Environ(),
			Dir:    wd,
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		},
		Executable: executableConfig,
	}, nil
}

var WR_VAR_PATTERN = regexp.MustCompile("^-wr_var='?([\\w\\-_]+)=(.+?)'?$")

func getArgAndVars(rawArgs []string) ([]string, map[string]string) {

	args := []string{}
	vars := map[string]string{}

	// TODO - This is quite brittle and needs something more robust
	// Need to test all cases, including:
	// - different positions in the command line
	// - different quoting
	for _, arg := range rawArgs {

		matches := WR_VAR_PATTERN.FindStringSubmatch(arg)

		if len(matches) == 3 {
			key := matches[1]
			value := matches[2]
			vars[key] = value
		} else {
			args = append(args, arg)
		}
	}

	return args, vars
}
