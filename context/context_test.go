package context

import (
	"os"
	"reflect"
	"testing"

	"os/exec"

	"github.com/activatedio/wrangle/config"
)

func TestNewContext(t *testing.T) {

	wd, err := os.Getwd()

	check(err)

	cases := map[string]struct {
		config   *config.Config
		expected *Context
	}{
		"empty": {
			config: &config.Config{},
			expected: &Context{
				Delegate: &exec.Cmd{
					Path:   "/usr/local/bin/terraform",
					Args:   os.Args,
					Env:    os.Environ(),
					Dir:    wd,
					Stdin:  os.Stdin,
					Stdout: os.Stdout,
					Stderr: os.Stderr,
				},
			},
		},
		"delegate": {
			config: &config.Config{
				Delegate: "/bin/ls",
			},
			expected: &Context{
				Delegate: &exec.Cmd{
					Path:   "/bin/ls",
					Args:   os.Args,
					Env:    os.Environ(),
					Dir:    wd,
					Stdin:  os.Stdin,
					Stdout: os.Stdout,
					Stderr: os.Stderr,
				},
			},
		},
	}

	for k, v := range cases {

		t.Run(k, func(t *testing.T) {

			if _, err := os.Stat(v.expected.Delegate.Path); os.IsNotExist(err) {
				t.Skipf("%s executable does not exist", v.expected.Delegate)
			}

			got, err := NewContext(v.config)

			if err != nil {
				t.Fatalf("Unexpected error %s", err)
			}

			if !reflect.DeepEqual(v.expected, got) {
				t.Fatalf("Wanted %v, got %v", v.expected, got)
			}
		})
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
