package context

import (
	"os"
	"reflect"
	"testing"

	"github.com/activatedio/wrangle/config"
)

func TestNewContext(t *testing.T) {

	args := os.Args[1:]

	cases := map[string]struct {
		config   *config.Config
		expected *Context
	}{
		"empty": {
			config: &config.Config{},
			expected: &Context{
				Delegate: "/usr/local/bin/terraform",
				Args:     args,
			},
		},
		"delegate": {
			config: &config.Config{
				Delegate: "/bin/ls",
			},
			expected: &Context{
				Delegate: "/bin/ls",
				Args:     args,
			},
		},
	}

	for k, v := range cases {

		t.Run(k, func(t *testing.T) {

			if _, err := os.Stat(v.expected.Delegate); os.IsNotExist(err) {
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
