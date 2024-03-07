//go:build !windows && !wasm
// +build !windows,!wasm

// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs

package sys

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

func checkPid(pid int, processName string) bool {
	switch runtime.GOOS {
	case "linux":
		return checkPidLinux(pid, processName)
	case "darwin":
		return checkPidDarwin(pid, processName)
	default:
		return true
	}
}

func checkPidLinux(pid int, processName string) bool {
	cmdlinePath := fmt.Sprintf("/proc/%d/cmdline", pid)
	data, err := os.ReadFile(cmdlinePath)
	if err != nil {
		return false
	}
	cmdLine := strings.FieldsFunc(string(data), func(r rune) bool { return r == 0 })
	fullCmdLine := strings.Join(cmdLine, " ")
	return strings.Contains(fullCmdLine, processName)
}

func checkPidDarwin(pid int, processName string) bool {
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-ocommand=")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	trimmedOutput := strings.TrimSpace(string(output))
	return strings.Contains(trimmedOutput, processName)
}

func sendTerminated(pidStr string) {
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		fmt.Printf("invalid process ID: %s\n", pidStr)
		os.Exit(1)
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Printf("could not find process with PID %d: %v\n", pid, err)
		os.Exit(1)
	}

	if !checkPid(pid, "wfs") {
		fmt.Printf("not a wfs execution process %d\n", pid)
		os.Exit(1)
	}

	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		fmt.Printf("failed to send SIGTERM to process %d: %v\n", pid, err)
		os.Exit(1)
	}
	fmt.Printf("stop wfs for process %d.\n", pid)
}
