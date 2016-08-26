package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

type compiler struct {
	ph pathHelper
}

func (c compiler) getLocalHash() (h string, err error) {
	hashPath := c.ph.getLocalHashPath()
	compiledHashBytes, err := ioutil.ReadFile(hashPath)
	if err != nil {
		return
	}
	h = strings.TrimRight(string(compiledHashBytes), "\n")
	return
}

func (c compiler) compile() error {
	configData, err := ioutil.ReadFile(c.ph.getConfigPath())
	if err != nil {
		return err
	}
	config := &Config{}
	_, err = toml.Decode(string(configData), config)
	if err != nil {
		return err
	}
	compiledFileName := c.ph.getLocalCompiledBashProfilePath()
	compiledFileNameTmp := compiledFileName + ".tmp"
	allContent := bytes.NewBuffer([]byte{})
	baseDir := c.ph.getLocalBaseDir()
	for _, file := range config.BashProfile {
		fullPath := baseDir + "/" + file
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
	err = ioutil.WriteFile(compiledFileNameTmp, allContent.Bytes(), 0755)
	if err != nil {
		return err
	}
	err = os.Rename(compiledFileNameTmp, compiledFileName)
	if err != nil {
		return err
	}
	hashFileName := c.ph.getLocalHashPath()
	hashFileNameTmp := hashFileName + ".tmp"
	sum := md5.Sum(allContent.Bytes())
	err = ioutil.WriteFile(
		hashFileNameTmp,
		[]byte(hex.EncodeToString(sum[0:len(sum)])),
		0755,
	)
	if err != nil {
		return err
	}
	return os.Rename(hashFileNameTmp, hashFileName)
}

func (c compiler) install(ex executor, profilePath string, compiledProfilePath string) error {
	profilePathTmp := profilePath + ".tmp"
	cmd := ex.command(fmt.Sprintf("touch %s", profilePath), false)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprint("Touch error: ", err.Error(), " : ", string(out)))
	}
	cmd = ex.command(fmt.Sprintf(
		"cat %s | grep -vF %s | cat > %s && echo \"source %s\" >> %s && mv %s %s",
		profilePath,
		BORSSH_DIR,
		profilePathTmp,
		compiledProfilePath,
		profilePathTmp,
		profilePathTmp,
		profilePath,
	), true)
	out, err = cmd.CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprint("Replace error: ", err.Error(), " : ", string(out)))
	}
	return nil
}

func NewCompiler() (c compiler, err error) {
	ph, err := getPathHelper()
	if err != nil {
		return
	}
	c.ph = ph
	return
}
