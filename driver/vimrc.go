package driver

import (
	"github.com/bo0rsh201/borssh/common"
)

const COMPILED_VIMRC_FILE = "vimrc.compiled"
const VIMRC_FILE = ".vimrc"

type VimRcDriver struct {
	InitableDriver
}

func (d *VimRcDriver) GetSourceFileNames() []string {
	return d.baseGetSourceFileNames(d.c.VimRc)
}

func (d *VimRcDriver) GetDstFileName(isLocal bool) string {
	return d.ph.GetBaseDir(isLocal) + "/" + COMPILED_VIMRC_FILE
}

func (d *VimRcDriver) IsEmpty() bool {
	return len(d.c.VimRc) == 0
}

func (d *VimRcDriver) CleanInstalled(ex *common.Executor) error {
	return d.baseCleanInstalled(
		ex,
		d.GetDstFileName(ex.IsLocal),
		VIMRC_FILE,
	)
}

func (d *VimRcDriver) Install(ex *common.Executor) error {
	return d.baseInstall(
		ex,
		d.GetDstFileName(ex.IsLocal),
		VIMRC_FILE,
		"source",
	)
}
