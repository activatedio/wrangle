package context

import (
	"os/exec"

	"os"

	"errors"

	"fmt"

	"path/filepath"

	"github.com/activatedio/wrangle/config"
)

type Context struct {
	Delegate   *exec.Cmd
	Executable *config.Executable
}

var getArgs = func() []string {
	return os.Args
}

func NewContext(c *config.Config) (*Context, error) {

	args := getArgs()

	if len(args) < 2 {
		return nil, errors.New("Syntax: wrangle [program name] [args]")
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
