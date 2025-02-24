// This source file is part of the EdgeDB open source project.
//
// Copyright EdgeDB Inc. and the EdgeDB authors.
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

package testserver

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/geldata/gel-go/gelcfg"
	"github.com/geldata/gel-go/geltypes"
	gelint "github.com/geldata/gel-go/internal/client"
)

var (
	testServerInfo = filepath.Join(
		os.TempDir(),
		"edgedb-go-test-server-info",
	)
	opts gelcfg.Options
	once sync.Once
)

type info struct {
	TLSCertFile string `json:"tls_cert_file"`
	Port        int    `json:"port"`
	PID         int    `json:"pid"`
}

func (i *info) options() gelcfg.Options {
	return gelcfg.Options{
		Host:     "127.0.0.1",
		Port:     i.Port,
		User:     "test",
		Password: geltypes.NewOptionalStr("shhh"),
		TLSOptions: gelcfg.TLSOptions{
			CAFile:       i.TLSCertFile,
			SecurityMode: gelcfg.TLSModeNoHostVerification,
		},
	}
}

// Options starts the test server if it isn't already running and returns the
// connection options.
func Options() gelcfg.Options {
	once.Do(initServerInfo)
	return opts
}

// Fatal prints a stack trace, logs the error and exits the process.
func Fatal(err error) {
	debug.PrintStack()
	log.Fatal(err)
}

func readCachedInfoFile() (info *info, err error) {
	defer func() {
		if err != nil {
			log.Println("error reading test server info file:", err)
		}
	}()

	var data []byte
	data, err = os.ReadFile(testServerInfo)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}

	o := info.options()
	o.WaitUntilAvailable = 500 * time.Millisecond
	pool, err := gelint.NewPool("", o)
	if err != nil {
		return nil, err
	}

	err = pool.EnsureConnected(context.Background())
	if err != nil {
		return nil, err
	}

	return info, nil
}

// convert a windows path to a unix path for systems with WSL.
func getWSLPath(path string) string {
	path = strings.ReplaceAll(path, "C:", "/mnt/c")
	path = strings.ReplaceAll(path, `\`, "/")
	path = strings.ToLower(path)

	return path
}

func readServerStatusFile(fileName string) (*info, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close() // nolint:errcheck

	var line string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()
		if strings.HasPrefix(line, "READY=") {
			break
		}
	}

	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	if line == "" {
		return nil, errors.New("no data found in " + fileName)
	}

	var inf info
	line = strings.TrimPrefix(line, "READY=")
	err = json.Unmarshal([]byte(line), &inf)
	if err != nil {
		return nil, err
	}

	return &inf, nil
}

func startServerProcess() *info {
	log.Print("starting test server")
	defer log.Print("test server started")

	serverBin := os.Getenv("EDGEDB_SERVER_BIN")
	if serverBin == "" {
		serverBin = "edgedb-server"
	}

	dir, err := os.MkdirTemp("", "")
	if err != nil {
		Fatal(err)
	}

	statusFile := path.Join(dir, "status-file")
	log.Println("status file:", dir)

	statusFileUnix := getWSLPath(statusFile)

	args := []string{serverBin}
	if runtime.GOOS == "windows" {
		args = append([]string{"wsl", "-u", "edgedb"}, args...)
	}

	autoShutdownAfter := os.Getenv("EDGEDB_SERVER_AUTO_SHUTDOWN_AFTER_SECONDS")
	if autoShutdownAfter == "" {
		autoShutdownAfter = "10"
	}

	args = append(
		args,
		"--temp-dir",
		"--testmode",
		"--port=auto",
		"--emit-server-status="+statusFileUnix,
		"--tls-cert-mode=generate_self_signed",
		"--auto-shutdown-after="+autoShutdownAfter,
		`--bootstrap-command=`+
			`CREATE SUPERUSER ROLE test { SET password := "shhh" }`,
	)

	log.Println("starting server with:", strings.Join(args, " "))

	cmd := exec.Command(args[0], args[1:]...)

	if os.Getenv("EDGEDB_SILENT_SERVER") == "" {
		fmt.Print(`
-------------------------------------------------------------------------------
Forwarding server's stderr. Set EDGEDB_SILENT_SERVER=1 to suppress.
-------------------------------------------------------------------------------

`)
		cmd.Stderr = os.Stderr
	} else {
		fmt.Print(`
-------------------------------------------------------------------------------
EDGEDB_SILENT_SERVER is set. Hiding server's stderr.
-------------------------------------------------------------------------------

`)
	}

	if os.Getenv("EDGEDB_DEBUG_SERVER") != "" {
		fmt.Print(`
-------------------------------------------------------------------------------
EDGEDB_DEBUG_SERVER is set. Forwarding server's stdout.
-------------------------------------------------------------------------------

`)
		cmd.Stdout = os.Stdout
	} else {
		fmt.Print(`
-------------------------------------------------------------------------------
Set EDGEDB_DEBUG_SERVER=1 to see server debug logs.
-------------------------------------------------------------------------------

`)
	}

	if os.Getenv("CI") == "" && os.Getenv("EDGEDB_SERVER_BIN") == "" {
		cmd.Env = append(os.Environ(),
			"__EDGEDB_DEVMODE=1",
		)
	}

	err = cmd.Start()
	if err != nil {
		Fatal(err)
	}

	log.Println("waiting for test server connection info")
	var localInfo *info
	for i := 0; i < 250; i++ {
		localInfo, err = readServerStatusFile(statusFile)
		if err == nil && localInfo != nil {
			break
		}
		time.Sleep(time.Second)
	}

	if err != nil {
		_ = cmd.Process.Kill()
		Fatal(err)
	}

	if len(localInfo.TLSCertFile) != 0 && runtime.GOOS == "windows" {
		tmpFile := path.Join(dir, "edbtlscert.pem")
		_, err = exec.Command(
			"wsl",
			"-u",
			"edgedb",
			"cp",
			localInfo.TLSCertFile,
			getWSLPath(tmpFile),
		).Output()
		if err != nil {
			Fatal(err)
		}
		localInfo.TLSCertFile = tmpFile
	}

	localInfo.PID = cmd.Process.Pid
	data, err := json.Marshal(localInfo)
	if err != nil {
		Fatal(err)
	}

	err = os.WriteFile(testServerInfo, data, 0777)
	if err != nil {
		Fatal(err)
	}

	return localInfo
}

func initServerInfo() {
	i, err := readCachedInfoFile()
	if err != nil {
		i = startServerProcess()
	}

	opts = i.options()
}
