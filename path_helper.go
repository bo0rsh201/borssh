package main

import (
	"os"
	"errors"
	"fmt"
)

const BORSSH_DIR = ".borssh"
const CONFIG_FILE = "config.toml"
const COMPILED_HASH_FILE = "hash.compiled"
const COMPILED_BASH_PROFILE_FILE = "bash_profile.compiled"
const BASH_PROFILE_FILE = ".bash_profile"

type pathHelper struct {
	homeDir string
}

func (ph pathHelper) getConfigPath() string {
	return fmt.Sprintf("%s/%s/%s", ph.homeDir, BORSSH_DIR, CONFIG_FILE)
}

func (ph pathHelper) getLocalHashPath() string {
	return fmt.Sprintf("%s/%s/%s", ph.homeDir, BORSSH_DIR, COMPILED_HASH_FILE)
}

func (ph pathHelper) getRemoteHashPath() string {
	return fmt.Sprintf("~/%s/%s", BORSSH_DIR, COMPILED_HASH_FILE)
}

func (ph pathHelper) getLocalCompiledBashProfilePath() string {
	return fmt.Sprintf("%s/%s/%s", ph.homeDir, BORSSH_DIR, COMPILED_BASH_PROFILE_FILE)
}

func (ph pathHelper) getRemoteCompiledBashProfilePath() string {
	return fmt.Sprintf("~/%s/%s", BORSSH_DIR, COMPILED_BASH_PROFILE_FILE)
}

func (ph pathHelper) getLocalBashProfilePath() string {
	return fmt.Sprintf("%s/%s", ph.homeDir, BASH_PROFILE_FILE)
}

func (ph pathHelper) getRemoteBashProfilePath() string {
	return fmt.Sprintf("~/%s", BASH_PROFILE_FILE)
}

func (ph pathHelper) getLocalBaseDir() string {
	return fmt.Sprintf("%s/%s", ph.homeDir, BORSSH_DIR)
}

func (ph pathHelper) getRemoteHomePath() string {
	return "~/"
}

func (ph pathHelper) getRemoteBaseDir() string {
	return fmt.Sprintf("~/%s/compiled", BORSSH_DIR)
}

func (ph pathHelper) prepareTmpDir() (tmpDir string, err error) {
	tmpDir = os.TempDir() + "/borssh_tmp"
	err = os.RemoveAll(tmpDir)
	if err != nil {
		return
	}
	err = os.Mkdir(tmpDir, 0777)
	return
}

func getPathHelper() (ph pathHelper, err error) {
	homeDir, ok := os.LookupEnv("HOME")
	if !ok {
		err = errors.New("Cannot detect home dir")
		return
	}
	ph.homeDir = homeDir
	return
}
