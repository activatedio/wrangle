package context

import (
	"os/exec"

	"os"

	"github.com/activatedio/wrangle/config"
)

// TODO - eventually this will come from a config
const DEFAULT_DELEGATE_NAME = "terraform"

type Context struct {
	Delegate *exec.Cmd
}

func NewContext(config *config.Config) (*Context, error) {

	p := config.Delegate

	if p == "" {
		p = DEFAULT_DELEGATE_NAME
	}

	p, err := exec.LookPath(p)

	if err != nil {
		return nil, err
	}

	wd, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	return &Context{
		Delegate: &exec.Cmd{
			Path:   p,
			Args:   os.Args,
			Env:    os.Environ(),
			Dir:    wd,
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		},
	}, nil
}
