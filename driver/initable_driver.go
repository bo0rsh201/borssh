package driver

import (
	"errors"
	"fmt"

	"github.com/bo0rsh201/borssh/common"
)

type InitableDriver struct {
	ph *common.PathHelper
	c  *common.Config
}

func (d *InitableDriver) Init(ph *common.PathHelper, c *common.Config) {
	d.ph = ph
	d.c = c
}

func (d *InitableDriver) baseGetSourceFileNames(files []string) []string {
	// source files are always local
	baseDir := d.ph.GetBaseDir(true)
	res := make([]string, 0, len(files))
	for _, file := range files {
		res = append(res, baseDir+"/"+file)
	}
	return res
}

func (d *InitableDriver) baseCleanInstalled(ex *common.Executor, compiledFile, parentFileName string) error {
	parentFile := d.ph.GetHomeDir(ex.IsLocal) + "/" + parentFileName
	parentFileTmp := parentFile + ".tmp"
	cmd := ex.Command(fmt.Sprintf(
		"if [ -f %s ]; then cat %s | grep -vF \"%s\" | cat > %s && mv %s %s; fi;",
		parentFile,
		parentFile,
		compiledFile,
		parentFileTmp,
		parentFileTmp,
		parentFile,
	), true)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprint("Clean error: ", err.Error(), " : ", string(out)))
	}
	return nil
}

func (d *InitableDriver) baseInstall(ex *common.Executor, compiledFile, parentFileName, includeCommand string) error {
	parentFile := d.ph.GetHomeDir(ex.IsLocal) + "/" + parentFileName
	cmd := ex.Command(fmt.Sprintf(
		"echo \"%s %s\" >> %s",
		includeCommand,
		compiledFile,
		parentFile,
	), true)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprint("Install error: ", err.Error(), " : ", string(out)))
	}
	return nil
}
