package common

import (
	"errors"
	"os"
)

const BORSSH_DIR = ".borssh"
const CONFIG_FILE = "config"
const COMPILED_HASH_FILE = "hash.compiled"

type PathHelper struct {
	localHomeDir  string
	remoteHomeDir string
	localBaseDir  string
	remoteBaseDir string
}

func (ph *PathHelper) init(homeDir string) {
	ph.localHomeDir = homeDir
	ph.localBaseDir = homeDir + "/" + BORSSH_DIR
	ph.remoteHomeDir = "~"
	ph.remoteBaseDir = ph.remoteHomeDir + "/" + BORSSH_DIR
}

func (ph *PathHelper) GetBaseDir(isLocal bool) string {
	if isLocal {
		return ph.localBaseDir
	}
	return ph.remoteBaseDir
}

func (ph *PathHelper) GetHomeDir(isLocal bool) string {
	if isLocal {
		return ph.localHomeDir
	}
	return ph.remoteHomeDir
}

func (ph *PathHelper) GetConfigPath(isLocal bool) string {
	return ph.GetBaseDir(isLocal) + "/" + CONFIG_FILE
}

func (ph *PathHelper) GetHashPath(isLocal bool) string {
	return ph.GetBaseDir(isLocal) + "/" + COMPILED_HASH_FILE
}

func GetPathHelper() (ph *PathHelper, err error) {
	homeDir, ok := os.LookupEnv("HOME")
	if !ok {
		err = errors.New("Cannot detect home dir")
		return
	}
	ph = &PathHelper{}
	ph.init(homeDir)
	return
}
