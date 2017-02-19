package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/bo0rsh201/borssh/common"
)

type compiler struct {
	config  *common.Config
	drivers []Driver
}

func (c compiler) compile() error {
	allHashes := bytes.NewBuffer([]byte{})
	for _, driver := range c.drivers {
		compiledFileName := driver.GetDstFileName(true)
		// removing previously compiled file if it had to be
		if compiledFileName != "" {
			err := os.Remove(compiledFileName)
			if err != nil && !os.IsNotExist(err) {
				return err
			}
		}
		if driver.IsEmpty() {
			continue
		}
		allContent := bytes.NewBuffer([]byte{})
		compiledFileNameTmp := compiledFileName + ".tmp"
		// join all files into one ".compiled"
		for _, fullPath := range driver.GetSourceFileNames() {
			content, err := ioutil.ReadFile(fullPath)
			if err != nil {
				return err
			}
			err = allContent.WriteByte('\n')
			if err != nil {
				return err
			}
			_, err = allContent.Write(content)
			if err != nil {
				return err
			}
		}
		// empty target name means that we should only do rsync
		// no compile and no install (for custom files e.x.)
		if compiledFileName != "" {
			err := ioutil.WriteFile(compiledFileNameTmp, allContent.Bytes(), 0755)
			if err != nil {
				return err
			}
			err = os.Rename(compiledFileNameTmp, compiledFileName)
			if err != nil {
				return err
			}
		}
		sum := md5.Sum(allContent.Bytes())
		_, err := allHashes.Write(sum[:])
		if err != nil {
			return err
		}
	}

	hashFileName := paths.GetHashPath(true)
	hashFileNameTmp := hashFileName + ".tmp"
	sum := md5.Sum(allHashes.Bytes())
	err := ioutil.WriteFile(
		hashFileNameTmp,
		[]byte(hex.EncodeToString(sum[:])),
		0755,
	)
	if err != nil {
		return err
	}
	return os.Rename(hashFileNameTmp, hashFileName)
}

func (c compiler) install(ex *common.Executor) error {
	for _, driver := range c.drivers {
		err := driver.CleanInstalled(ex)
		if err != nil {
			return err
		}
		if driver.IsEmpty() {
			continue
		}
		err = driver.Install(ex)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewCompiler() (c compiler, err error) {
	var configData []byte
	// we need to include local config, to we pass true
	configData, err = ioutil.ReadFile(paths.GetConfigPath(true))
	if err != nil {
		return
	}

	config := &common.Config{}
	_, err = toml.Decode(string(configData), config)
	if err != nil {
		return
	}
	c.config = config
	c.drivers = getAllDrivers(config)
	return
}
