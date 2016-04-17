package main

import (
	"os"
	"errors"
)

const BORSSH_DIR = ".borssh"
const CONFIG_FILE = "config.toml"
const COMPILED_HASH_FILE = "hash.compiled"
const COMPILED_BASH_PROFILE_FILE = "bash_profile.compiled"
const BASH_PROFILE_FILE = ".bash_profile"

type pathHelper struct {
	homeDir string
	remoteHomeDir string
	localBaseDir string
	remoteBaseDir string
}

func (ph *pathHelper) init(homeDir string) {
	ph.homeDir = homeDir
	ph.localBaseDir = homeDir + "/" + BORSSH_DIR
	ph.remoteHomeDir = "~"
	ph.remoteBaseDir = ph.remoteHomeDir + "/" + BORSSH_DIR
}

func (ph pathHelper) getConfigPath() string {
	return ph.localBaseDir + "/" + CONFIG_FILE
}

func (ph pathHelper) getLocalHashPath() string {
	return ph.localBaseDir + "/" + COMPILED_HASH_FILE
}
func (ph pathHelper) getRemoteHashPath() string {
	return ph.remoteBaseDir + "/" + COMPILED_HASH_FILE
}

func (ph pathHelper) getLocalCompiledBashProfilePath() string {
	return ph.localBaseDir + "/" + COMPILED_BASH_PROFILE_FILE
}
func (ph pathHelper) getRemoteCompiledBashProfilePath() string {
	return ph.remoteBaseDir + "/" + COMPILED_BASH_PROFILE_FILE
}

func (ph pathHelper) getLocalBashProfilePath() string {
	return ph.homeDir + "/" + BASH_PROFILE_FILE
}
func (ph pathHelper) getRemoteBashProfilePath() string {
	return ph.remoteHomeDir + "/" + BASH_PROFILE_FILE
}

func (ph pathHelper) getLocalBaseDir() string {
	return ph.localBaseDir
}
func (ph pathHelper) getRemoteHomePath() string {
	return ph.remoteHomeDir
}

func getPathHelper() (ph pathHelper, err error) {
	homeDir, ok := os.LookupEnv("HOME")
	if !ok {
		err = errors.New("Cannot detect home dir")
		return
	}
	ph.init(homeDir)
	return
}
