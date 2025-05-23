// This source file is part of the Gel open source project.
//
// Copyright Gel Data Inc. and the Gel authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/geldata/gel-go/internal/testserver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	dsn         string
	projectRoot = getProjectRoot()
)

func getProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "../..")
}

var tests = []struct {
	description string
	directory   string
	args        []string
}{
	{
		description: "invoke edgeql-go without args",
		directory:   "testdata/no-args",
		args:        []string{},
	},
	{
		description: "invoke edgeql-go with -mixedcaps",
		directory:   "testdata/mixedcaps",
		args:        []string{"-mixedcaps"},
	},
	{
		description: "invoke edgeql-go with -pubfuncs",
		directory:   "testdata/pubfuncs",
		args:        []string{"-pubfuncs"},
	},
	{
		description: "invoke edgeql-go with -pubtypes",
		directory:   "testdata/pubtypes",
		args:        []string{"-pubtypes"},
	},
	{
		description: "invoke edgeql-go with -rawmessage",
		directory:   "testdata/rawmessage",
		args:        []string{"-rawmessage"},
	},
}

func TestMain(m *testing.M) {
	o := testserver.Options()
	dsn = testserver.AsDSN(o)
	os.Exit(m.Run())
}

func TestEdgeQLGo(t *testing.T) {
	for _, test := range tests {
		t.Run(test.description, runTest(test.directory, test.args))
	}
}

func buildEdgeqlGo(t *testing.T, dir string) string {
	t.Log("building edgeql-go")
	edgeqlGo := filepath.Join(dir, "edgeql-go")
	requireRun(t, ".", "go", "build", "-o", edgeqlGo)
	return edgeqlGo
}

func replaceGelGoModule(t *testing.T, dir string) {
	// Run tests against the current checkout of gel-go instead of
	// against whatever older version is in the test project's
	// go.mod file.
	replace := fmt.Sprintf(
		"-replace=github.com/geldata/gel-go=%s",
		projectRoot,
	)
	requireRun(t, dir, "go", "mod", "edit", replace)
	requireRun(t, dir, "go", "mod", "tidy")
}

func runTest(dir string, args []string) func(*testing.T) {
	return func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "edgeql-go-*")
		require.NoError(t, err)
		defer func() {
			assert.NoError(t, os.RemoveAll(tmpDir))
		}()

		edgeqlGo := buildEdgeqlGo(t, tmpDir)

		var wg sync.WaitGroup
		err = filepath.WalkDir(
			dir,
			func(src string, d fs.DirEntry, e error) error {
				require.NoError(t, e)
				if src == dir {
					return nil
				}

				dst := filepath.Join(tmpDir, strings.TrimPrefix(src, dir))
				if d.IsDir() {
					e = os.Mkdir(dst, os.ModePerm)
					require.NoError(t, e)
				} else {
					wg.Add(1)
					go func() {
						defer wg.Done()
						copyFile(t, dst, src)
					}()
				}
				return nil
			},
		)
		require.NoError(t, err)
		wg.Wait()

		entries, err := os.ReadDir(tmpDir)
		require.NoError(t, err)
		for _, entry := range entries {
			if entry.Name() == "edgeql-go" {
				continue
			}

			t.Run(entry.Name(), func(t *testing.T) {
				projectDir := filepath.Join(tmpDir, entry.Name())

				replaceGelGoModule(t, projectDir)

				requireRun(t, projectDir, edgeqlGo, args...)
				requireRun(t, projectDir, "go", "run", "./...")
				er := filepath.WalkDir(
					projectDir,
					func(f string, d fs.DirEntry, e error) error {
						require.NoError(t, e)
						if strings.HasSuffix(f, ".go.assert") {
							checkAssertFile(t, f)
						}
						if strings.HasSuffix(f, ".go") &&
							!strings.HasSuffix(f, "ignore.go") &&
							!strings.HasSuffix(f, "_test.go") {
							checkGoFile(t, f)
						}
						return nil
					},
				)
				require.NoError(t, er)
				requireRun(t, projectDir, "go", "test", "-count=1", "./...")
			})
		}
	}
}

func checkAssertFile(t *testing.T, file string) {
	t.Helper()
	goFile := strings.TrimSuffix(file, ".assert")
	if assert.FileExistsf(t, goFile, "missing .go file for %s", file) {
		assertEqualFiles(t, file, goFile)
	}
}

func checkGoFile(t *testing.T, file string) {
	t.Helper()
	assertFile := file + ".assert"
	if assert.FileExistsf(t, assertFile,
		"missing .go.assert file for %s", file,
	) {
		assertEqualFiles(t, assertFile, file)
	}
}

func assertEqualFiles(t *testing.T, left, right string) {
	t.Helper()
	leftData, err := os.ReadFile(left)
	require.NoErrorf(t, err, "reading %s", left)

	rightData, err := os.ReadFile(right)
	require.NoErrorf(t, err, "reading %s", right)

	assert.Equal(t, string(leftData), string(rightData),
		"files are not equal: %s != %s", left, right,
	)
}

func copyFile(t *testing.T, to, from string) {
	toFd, err := os.Create(to)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, toFd.Close())
	}()

	fromFd, err := os.Open(from)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, fromFd.Close())
	}()

	_, err = io.Copy(toFd, fromFd)
	require.NoError(t, err)
}

func run(t *testing.T, dir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), fmt.Sprintf("EDGEDB_DSN=%s", dsn))
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func requireRun(t *testing.T, dir, name string, args ...string) {
	output, err := run(t, dir, name, args...)
	require.NoError(t, err, output)
}

func TestEdgeQLGoQueryError(t *testing.T) {
	srcDir := "testdata/query-error/"
	tmpDir := t.TempDir()
	tmpDir, err := filepath.EvalSymlinks(tmpDir)
	require.NoError(t, err)

	matches, err := filepath.Glob(filepath.Join(srcDir, "*"))
	require.NoError(t, err)

	for _, src := range matches {
		dst := filepath.Join(tmpDir, filepath.Base(src))
		copyFile(t, dst, src)
	}
	replaceGelGoModule(t, tmpDir)
	edgeqlGo := buildEdgeqlGo(t, tmpDir)
	output, err := run(t, tmpDir, edgeqlGo)
	assert.EqualError(t, err, "exit status 1")

	//nolint:lll
	expected := fmt.Sprintf(
		`edgeql-go: failed to setup query: error introspecting query "%[1]s/broken_query.edgeql": gel.EdgeQLSyntaxError: Unexpected 'malformed'
%[1]s/broken_query.edgeql:1:1

malformed
^ error
`,
		tmpDir,
	)
	assert.Equal(t, expected, output)
}
