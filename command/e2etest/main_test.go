package e2etest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"os/exec"

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

func TestRun(t *testing.T) {

	cases := map[string]struct {
		delegate           string
		options            []string
		expectedExitStatus int
		verify             func(t *testing.T, b *e2e.Binary, stdout string, stderr string)
	}{
		"multiple-executables": {
			delegate: "./delegate2.sh",
			verify: func(t *testing.T, b *e2e.Binary, stdout string, stderr string) {
				wantStdout := "2\n"
				if wantStdout != stdout {
					t.Fatalf("Stdout: wanted \n[%s]\n, got \n[%s]",
						wantStdout, stdout)
				}
			},
		},
		"template-only": {
			delegate: "./delegate.sh",
			verify: func(t *testing.T, b *e2e.Binary, stdout string, stderr string) {
				f := "main.tf"
				if !b.FileExists(f) {
					t.Fatalf("Expected file %s to exist", f)
				}
				wantStdout := `    a = "a1"
    b = "b1"
`
				if wantStdout != stdout {
					t.Fatalf("Stdout: wanted \n[%s]\n, got \n[%s]",
						wantStdout, stdout)
				}
			},
		},
		"template-only-with-vars": {
			options: []string{
				"-wr_var=cs=cs1",
				"-wr_var=ds=ds1",
			},
			delegate: "./delegate.sh",
			verify: func(t *testing.T, b *e2e.Binary, stdout string, stderr string) {
				f := "main.tf"
				if !b.FileExists(f) {
					t.Fatalf("Expected file %s to exist", f)
				}
				wantStdout := `    a = "a1"
    b = "b1"
    cs = "cs1"
    ds = "ds1"
`
				if wantStdout != stdout {
					t.Fatalf("Stdout: wanted \n[%s]\n, got \n[%s]",
						wantStdout, stdout)
				}
			},
		},
		"template-only-with-vars-quoted": {
			options: []string{
				"-wr_var='cs=cs1'",
				"-wr_var='ds=ds1'",
			},
			delegate: "./delegate.sh",
			verify: func(t *testing.T, b *e2e.Binary, stdout string, stderr string) {
				f := "main.tf"
				if !b.FileExists(f) {
					t.Fatalf("Expected file %s to exist", f)
				}
				wantStdout := `    a = "a1"
    b = "b1"
    cs = "cs1"
    ds = "ds1"
`
				if wantStdout != stdout {
					t.Fatalf("Stdout: wanted \n[%s]\n, got \n[%s]",
						wantStdout, stdout)
				}
			},
		},
		"aws-user-data-only": {
			delegate: "./delegate.sh",
			verify: func(t *testing.T, b *e2e.Binary, stdout string, stderr string) {
				f := ".user-data.sh"
				if !b.FileExists(f) {
					t.Fatalf("Expected file %s to exist", f)
				}
				// TODO - Check some golden files here
			},
		},
		"error-1": {
			delegate:           "./delegate.sh",
			expectedExitStatus: 1,
			verify: func(t *testing.T, b *e2e.Binary, stdout string, stderr string) {

				stderrWantContains := " function \"regions\" not defined"

				if !strings.Contains(stderr, stderrWantContains) {
					t.Fatalf("Wanted stderr [%s] to contain [%s]", stderr, stderrWantContains)
				}

			},
		},
	}

	for k, v := range cases {

		t.Run(k, func(t *testing.T) {

			os.Setenv("AWS_USER_DATA_UID", "1234")
			defer func() {
				os.Unsetenv("AWS_USER_DATA_UID")
			}()

			fixturePath := filepath.Join("test-fixtures", k)
			wr := e2e.NewBinary(wrangleBin, fixturePath)
			defer wr.Close()

			args := []string{}

			args = append(args, v.options...)
			args = append(args, v.delegate)

			stdout, stderr, err := wr.Run(args...)
			fmt.Println(stdout)
			if err != nil {

				if exiterr, ok := err.(*exec.ExitError); ok {
					if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
						if status.ExitStatus() != v.expectedExitStatus {
							t.Fatalf("unexpected error: status: %d %s\nstderr:\n%s", status, err, stderr)
						}
					}
				} else {
					t.Fatalf("unexpected error: %s\nstderr:\n%s", err, stderr)
				}
			} else if v.expectedExitStatus != 0 {
				t.Fatalf("Expected exit status %d but got 0", v.expectedExitStatus)
			}

			v.verify(t, wr, stdout, stderr)
		})
	}

}
