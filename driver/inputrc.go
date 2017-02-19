package driver

import (
	"github.com/bo0rsh201/borssh/common"
)

const COMPILED_INPUTRC_FILE = "inputrc.compiled"
const INPUTRC_FILE = ".inputrc"

type InputRcDriver struct {
	InitableDriver
}

func (d *InputRcDriver) GetSourceFileNames() []string {
	return d.baseGetSourceFileNames(d.c.InputRc)
}

func (d *InputRcDriver) GetDstFileName(isLocal bool) string {
	return d.ph.GetBaseDir(isLocal) + "/" + COMPILED_INPUTRC_FILE
}

func (d *InputRcDriver) IsEmpty() bool {
	return len(d.c.InputRc) == 0
}

func (d *InputRcDriver) CleanInstalled(ex *common.Executor) error {
	return d.baseCleanInstalled(
		ex,
		d.GetDstFileName(ex.IsLocal),
		INPUTRC_FILE,
	)
}

func (d *InputRcDriver) Install(ex *common.Executor) error {
	return d.baseInstall(
		ex,
		d.GetDstFileName(ex.IsLocal),
		INPUTRC_FILE,
		"\\$include",
	)
}
