package main

import (
	"os/exec"
	"strings"
	"fmt"
	"errors"
)
const (
	DEFAULT_CONNECT_TIMEOUT = 10
	SERVER_ALIVE_INTERVAL   = 3
	SERVER_ALIVE_COUNT_MAX  = 4
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
	host string
	isLocal bool
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

func (e executor) command(cmd string) *exec.Cmd {
	if e.isLocal {
		return exec.Command("bash", "-c", cmd)
	}
	return exec.Command(
		e.sshBinary,
		append(e.getSshOptions(), e.host, "bash", fmt.Sprintf("-c '%s'", cmd))...,
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
		moreOpts...
	)
	args = append(
		args,
		src,
		e.host + ":" + dst,
	)
	return exec.Command(
		rsyncBinary,
		args...
	), nil
}