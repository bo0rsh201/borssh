package driver

import (
	"github.com/bo0rsh201/borssh/common"
)

type CustomFilesDriver struct {
	InitableDriver
}

func (d *CustomFilesDriver) GetSourceFileNames() []string {
	// source files are always local
	baseDir := d.ph.GetBaseDir(true)
	res := make([]string, 0, len(d.c.CustomFiles))
	for _, file := range d.c.CustomFiles {
		res = append(res, baseDir+"/"+file)
	}
	return res
}

func (d *CustomFilesDriver) GetDstFileName(isLocal bool) string {
	return ""
}

func (d *CustomFilesDriver) CleanInstalled(ex *common.Executor) error {
	return nil
}

func (d *CustomFilesDriver) IsEmpty() bool {
	return len(d.c.CustomFiles) == 0
}

func (d *CustomFilesDriver) Install(ex *common.Executor) error {
	return nil
}
