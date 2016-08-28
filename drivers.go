package main

import (
	"github.com/bo0rsh201/borssh/common"
	"github.com/bo0rsh201/borssh/driver"
)

type Driver interface {
	Init(ph *common.PathHelper, c *common.Config)
	IsEmpty() bool
	GetSourceFileNames() []string
	GetDstFileName(isLocal bool) string
	CleanInstalled(ex *common.Executor) error
	Install(ex *common.Executor) error
}

func getAllDrivers(ph *common.PathHelper, c *common.Config) []Driver {
	res := []Driver{
		&driver.BashProfileDriver{},
		&driver.VimRcDriver{},
		&driver.InputRcDriver{},
		&driver.CustomFilesDriver{},
	}
	for _, d := range res {
		d.Init(ph, c)
	}
	return res
}