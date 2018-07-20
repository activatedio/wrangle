package e2etest

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/activatedio/wrangle/e2e"
)

// Based on code from https://raw.githubusercontent.com/hashicorp/wrangle/master/command/e2etest/main_test.go

var wrangleBin string

func TestMain(m *testing.M) {
	teardown := setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() func() {
	if wrangleBin != "" {
		// this is pre-set when we're running in a binary produced from
		// the make-archive.sh script, since that builds a ready-to-go
		// binary into the archive. However, we do need to turn it into
		// an absolute path so that we can find it when we change the
		// working directory during tests.
		var err error
		wrangleBin, err = filepath.Abs(wrangleBin)
		if err != nil {
			panic(fmt.Sprintf("failed to find absolute path of wrangle executable: %s", err))
		}
		return func() {}
	}

	tmpFilename := e2e.GoBuild("github.com/activatedio/wrangle", "wrangle")

	// Make the executable available for use in tests
	wrangleBin = tmpFilename

	return func() {
		os.Remove(tmpFilename)
	}
}

/*
func helperProcess(s ...string) *exec.Cmd {
	cs := []string {"-test.run=TestHelperProcess", "--"}
	cs = append(cs, s...)
	env := [] string {
		"GO_WANT_HELPER_PROCESS=1",
	}

	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = append(env, os.Environ()...)
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)

	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}

		args = args[:1]
	}

	fmt.Println("Output here")
}
*/

func TestRun(t *testing.T) {

	cases := map[string]struct {
	} {
		"template-only": {},
	}

	for k, _ := range cases {

		t.Run(k, func(*testing.T) {

			fixturePath := filepath.Join("test-fixtures", k)
			wr := e2e.NewBinary(wrangleBin, fixturePath)
			defer wr.Close()

			_, stderr, err := wr.Run()
			if err != nil {
				t.Fatalf("unexpected init error: %s\nstderr:\n%s", err, stderr)
			}
		})
	}


	//cmd := helperProcess()

}


