package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const (
	DEFAULT_CONNECT_TIMEOUT = 10
	SERVER_ALIVE_INTERVAL   = 3
	SERVER_ALIVE_COUNT_MAX  = 4
	// this looks a bit ugly,
	// because we could have situation
	// when client login shell exits with this code
	// and our wrapper will have a false-positive reaction to hash mismatch
	// but it's probability is too small
	// because login shell exits with 0 by default,
	// if you haven't executed "exit 192" explicitly
	EXIT_HASH_MISMATCH = 192
)

func NewExecutor(host string) (e executor, err error) {
	e.host = host
	sshBinary, err := exec.LookPath("ssh")
	if err != nil {
		return
	}
	e.sshBinary = sshBinary
	return
}

func NewLocalExecutor() (e executor) {
	e.isLocal = true
	return e
}

type executor struct {
	host      string
	isLocal   bool
	sshBinary string
}

func (e executor) getSshOptions() []string {
	return []string{
		"-o", fmt.Sprint("ConnectTimeout=", DEFAULT_CONNECT_TIMEOUT),
		"-o", "LogLevel=ERROR",
		"-o", fmt.Sprint("ServerAliveInterval=", SERVER_ALIVE_INTERVAL),
		"-o", fmt.Sprint("ServerAliveCountMax=", SERVER_ALIVE_COUNT_MAX),
		"-o", "BatchMode=yes",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
	}
}

func (e executor) tryToConnect(hash, hashPath string) (bool, int, error) {
	cmd := e.command(fmt.Sprintf(
		"if [ -f %s ] && [ $(cat %s) == %s ]; then exec bash -l; else exit %d; fi;",
		hashPath,
		hashPath,
		hash,
		EXIT_HASH_MISMATCH,
	), true)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	exitCode, err := parseExecError(cmd.Run())
	// ssh/generic error
	if err != nil {
		return false, exitCode, err
	}
	if exitCode == EXIT_HASH_MISMATCH {
		return false, exitCode, nil
	}
	// remote command success/error
	return true, exitCode, nil
}

func (e executor) command(cmd string, isTerminal bool) *exec.Cmd {
	if e.isLocal {
		return exec.Command("bash", "-c", cmd)
	}
	sshOptions := e.getSshOptions()
	if isTerminal {
		sshOptions = append(sshOptions, "-t")
	}
	return exec.Command(
		e.sshBinary,
		append(sshOptions, e.host, "bash", fmt.Sprintf("-c '%s'", cmd))...,
	)
}

func (e executor) rsync(src, dst string, moreOpts ...string) (*exec.Cmd, error) {
	if e.isLocal {
		return nil, errors.New("Cannot perform rsync on local executor")
	}
	rsyncBinary, err := exec.LookPath("rsync")
	if err != nil {
		return nil, err
	}
	args := append(
		[]string{
			"-rlptz",
			"-e", e.sshBinary + " " + strings.Join(e.getSshOptions(), " "),
		},
		moreOpts...,
	)
	args = append(
		args,
		src,
		e.host+":"+dst,
	)
	return exec.Command(
		rsyncBinary,
		args...,
	), nil
}

func parseExecError(execError error) (int, error) {
	if execError == nil {
		return 0, nil
	}
	exitError, ok := execError.(*exec.ExitError)
	if !ok {
		return 1, execError
	}
	status, ok := exitError.Sys().(syscall.WaitStatus)
	if !ok {
		return 1, errors.New("Seems you don't have POSIX-compatible OS :" + exitError.Error())
	}
	return status.ExitStatus(), nil
}
