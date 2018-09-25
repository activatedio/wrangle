package context

import (
	"os"
	"reflect"
	"testing"

	"errors"
	"os/exec"

	"github.com/activatedio/wrangle/config"
)

func TestNewContext(t *testing.T) {

	var args []string

	origGetArgs := getArgs
	defer func() {
		getArgs = origGetArgs
	}()

	getArgs = func() []string {
		return args
	}

	wd, err := os.Getwd()

	check(err)

	_config := &config.Config{
		Executables: map[string]*config.Executable{
			"dummy.sh": &config.Executable{
				Plugins: map[string]interface{}{
					"a": "a1",
				},
			},
			"terraform": &config.Executable{
				Plugins: map[string]interface{}{
					"b": "b1",
				},
			},
			"ls": &config.Executable{
				Plugins: map[string]interface{}{
					"c": "c1",
				},
			},
		},
	}

	cases := map[string]struct {
		args     []string
		config   *config.Config
		expected *Context
		err      error
	}{
		"current-directory-path": {
			args:   []string{"wrangle", "./dummy.sh", "~/"},
			config: _config,
			expected: &Context{
				Variables: map[string]string{},
				Delegate: &exec.Cmd{
					Path:   "./dummy.sh",
					Args:   []string{"./dummy.sh", "~/"},
					Env:    os.Environ(),
					Dir:    wd,
					Stdin:  os.Stdin,
					Stdout: os.Stdout,
					Stderr: os.Stderr,
				},
				Executable: _config.Executables["dummy.sh"],
			},
		},
		"relative": {
			args:   []string{"wrangle", "ls", "~/"},
			config: _config,
			expected: &Context{
				Variables: map[string]string{},
				Delegate: &exec.Cmd{
					Path:   "/bin/ls",
					Args:   []string{"ls", "~/"},
					Env:    os.Environ(),
					Dir:    wd,
					Stdin:  os.Stdin,
					Stdout: os.Stdout,
					Stderr: os.Stderr,
				},
				Executable: _config.Executables["ls"],
			},
		},
		"full-path": {
			args:   []string{"wrangle", "/bin/ls", "~/"},
			config: _config,
			expected: &Context{
				Variables: map[string]string{},
				Delegate: &exec.Cmd{
					Path:   "/bin/ls",
					Args:   []string{"/bin/ls", "~/"},
					Env:    os.Environ(),
					Dir:    wd,
					Stdin:  os.Stdin,
					Stdout: os.Stdout,
					Stderr: os.Stderr,
				},
				Executable: _config.Executables["ls"],
			},
		},
		"wr-vars": {
			args:   []string{"wrangle", "-wr_var=a=abc", "-wr_var=b=def", "./dummy.sh", "~/"},
			config: _config,
			expected: &Context{
				Variables: map[string]string{
					"a": "abc",
					"b": "def",
				},
				Delegate: &exec.Cmd{
					Path:   "./dummy.sh",
					Args:   []string{"./dummy.sh", "~/"},
					Env:    os.Environ(),
					Dir:    wd,
					Stdin:  os.Stdin,
					Stdout: os.Stdout,
					Stderr: os.Stderr,
				},
				Executable: _config.Executables["dummy.sh"],
			},
		},
		"wr-vars-quoted": {
			args:   []string{"wrangle", "-wr_var='a=abc'", "-wr_var='b=def'", "./dummy.sh", "~/"},
			config: _config,
			expected: &Context{
				Variables: map[string]string{
					"a": "abc",
					"b": "def",
				},
				Delegate: &exec.Cmd{
					Path:   "./dummy.sh",
					Args:   []string{"./dummy.sh", "~/"},
					Env:    os.Environ(),
					Dir:    wd,
					Stdin:  os.Stdin,
					Stdout: os.Stdout,
					Stderr: os.Stderr,
				},
				Executable: _config.Executables["dummy.sh"],
			},
		},
		"relative-config-not-found": {
			args:   []string{"wrangle", "cat", "version"},
			config: _config,
			err:    errors.New("Executable [cat] not found in config."),
		},
		"full-path-config-not-found": {
			args:   []string{"wrangle", "/bin/cat", "~/"},
			config: _config,
			err:    errors.New("Executable [cat] not found in config."),
		},
		"too-few-args": {
			args:   []string{"wrangle"},
			config: _config,
			err:    errors.New("Syntax: wrangle [wrangle options] [program name] [program args]"),
		},
		// "executable-not-found-relative
		// "executable-not-found-absolute
	}

	for k, v := range cases {

		t.Run(k, func(t *testing.T) {

			if v.expected != nil {
				if _, err := os.Stat(v.expected.Delegate.Path); os.IsNotExist(err) {
					t.Skipf("%s executable does not exist", v.expected.Delegate)
				}
			}

			args = v.args

			got, err := NewContext(v.config)

			if err != nil && v.err == nil {
				t.Fatalf("Unexpected error %s", err)
			} else if !reflect.DeepEqual(err, v.err) {
				t.Fatalf("Error [%+v], got [%+v]", v.err, err)
			}

			if !reflect.DeepEqual(v.expected, got) {
				t.Fatalf("Wanted %+v, got %+v", v.expected, got)
			}
		})
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
