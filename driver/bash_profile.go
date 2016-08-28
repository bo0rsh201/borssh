package driver

import (
	"github.com/bo0rsh201/borssh/common"
)

const COMPILED_BASH_PROFILE_FILE = "bash_profile.compiled"
const BASH_PROFILE_FILE = ".bash_profile"

type BashProfileDriver struct {
	InitableDriver
}

func (d *BashProfileDriver) GetSourceFileNames() []string {
	return d.baseGetSourceFileNames(d.c.BashProfile)
}

func (d *BashProfileDriver) GetDstFileName(isLocal bool) string {
	return d.ph.GetBaseDir(isLocal) + "/" + COMPILED_BASH_PROFILE_FILE
}

func (d *BashProfileDriver) IsEmpty() bool {
	return len(d.c.BashProfile) == 0
}

func (d *BashProfileDriver) CleanInstalled(ex *common.Executor) error {
	return d.baseCleanInstalled(
		ex,
		d.GetDstFileName(ex.IsLocal),
		BASH_PROFILE_FILE,
	)
}

func (d *BashProfileDriver) Install(ex *common.Executor) error {
	return d.baseInstall(
		ex,
		d.GetDstFileName(ex.IsLocal),
		BASH_PROFILE_FILE,
		"source",
	)
}
