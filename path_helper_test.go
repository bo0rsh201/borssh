package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaths(t *testing.T) {
	ph := pathHelper{}
	ph.init("/home/testuser")

	assert.Equal(t, "/home/testuser/.borssh/config.toml", ph.getConfigPath())

	assert.Equal(t, "/home/testuser/.borssh/hash.compiled", ph.getLocalHashPath())
	assert.Equal(t, "~/.borssh/hash.compiled", ph.getRemoteHashPath())

	assert.Equal(t, "/home/testuser/.borssh/bash_profile.compiled", ph.getLocalCompiledBashProfilePath())
	assert.Equal(t, "~/.borssh/bash_profile.compiled", ph.getRemoteCompiledBashProfilePath())

	assert.Equal(t, "/home/testuser/.bash_profile", ph.getLocalBashProfilePath())
	assert.Equal(t, "~/.bash_profile", ph.getRemoteBashProfilePath())

	assert.Equal(t, "/home/testuser/.borssh", ph.getLocalBaseDir())

	assert.Equal(t, "~", ph.getRemoteHomePath())
}
