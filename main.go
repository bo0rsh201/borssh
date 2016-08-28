package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"github.com/bo0rsh201/borssh/common"
)

const COMMAND_COMPILE = "compile"
const COMMAND_CONNECT = "connect"
const FLAG_QUITE = "q"

func printUsage() {
	usage := "Usage: borssh COMMAND [flags...] [args...]\n\n"
	usage += "A simple ssh wrapper that helps you to keep all your dot files up to date\n\n"
	usage += "Commands:\n"
	usage += "\t%s\t\t\tCompiles dot files for latest config version\n"
	usage += "\t%s <target>\tConnect with latest compiled version\n"
	usage += "Flags:\n"
	usage += "\t-%s\t\t\tQuite mode (suppress all output)\n\n"
	usage += "More info at http://github.com/bo0rsh201/borssh\n"
	fmt.Fprintf(os.Stderr, usage, COMMAND_COMPILE, COMMAND_CONNECT, FLAG_QUITE)
	os.Exit(1)
}

func getLocalHash(localHashPath string) (h string, err error) {
	compiledHashBytes, err := ioutil.ReadFile(localHashPath)
	if err != nil {
		return
	}
	h = strings.TrimRight(string(compiledHashBytes), "\n")
	return
}

func main() {
	fset := flag.NewFlagSet("basic", flag.ContinueOnError)
	fset.SetOutput(ioutil.Discard)
	quite := fset.Bool("q", false, "")
	fset.Parse(os.Args[1:])
	args := fset.Args()
	if len(args) < 1 {
		printUsage()
	}
	pp := NewProgressPrinter(*quite)

	ph, err := common.GetPathHelper()
	pp.failOnError(err)

	switch args[0] {
	case COMMAND_CONNECT:
		if len(args) < 2 {
			printUsage()
		}

		c, err := NewCompiler(ph)
		pp.failOnError(err)

		localHash, err := getLocalHash(ph.GetHashPath(true))
		if os.IsNotExist(err) {
			err = errors.New("Cannot find compiled version of dot files. You should run compile command first")
		}
		pp.failOnError(err)

		ex, err := common.NewExecutor(args[1])
		pp.failOnError(err)

		ok, exitCode, err := ex.TryToConnect(localHash, ph.GetHashPath(false))
		pp.failOnError(err)
		if ok {
			os.Exit(exitCode)
		}

		pp.yellow("Hash file mismatch...")
		pp.Start("Syncing base dir")
		cmd, err := ex.Rsync(
			ph.GetBaseDir(true),
			ph.GetHomeDir(false)+"/",
			"--delete",
			"--copy-unsafe-links",
			"--exclude",
			common.COMPILED_HASH_FILE,
		)
		pp.CheckError(err)

		out, err := cmd.CombinedOutput()
		if err != nil {
			err = fmt.Errorf("Failed to rsync base dir: '%s' output: '%s'", err.Error(), string(out))
		}
		pp.CheckAndComplete(err)

		pp.Start("Installing remote")
		err = c.install(ex)
		pp.CheckAndComplete(err)

		pp.Start("Syncing hash file")
		cmd, err = ex.Rsync(
			ph.GetHashPath(true),
			ph.GetHashPath(false),
		)
		pp.CheckError(err)

		out, err = cmd.CombinedOutput()
		if err != nil {
			err = fmt.Errorf("Failed to rsync compiled hash: '%s' output: '%s'", err.Error(), string(out))
		}
		pp.CheckAndComplete(err)

		ok, exitCode, err = ex.TryToConnect(localHash, ph.GetHashPath(false))
		pp.failOnError(err)
		if ok {
			os.Exit(exitCode)
		}
		pp.failOnError(fmt.Errorf(
			"Hash mismatch after sync.\n"+
				"This should never happen.\n"+
				"Maybe, it's an exit_code collision (reserved code is %d)",
			common.EXIT_HASH_MISMATCH,
		))
		break
	case COMMAND_COMPILE:
		pp.Start("Compiling")
		c, err := NewCompiler(ph)
		pp.CheckError(err)
		err = c.compile(ph.GetHashPath(true))
		pp.CheckAndComplete(err)

		pp.Start("Installing")
		err = c.install(common.NewLocalExecutor())
		pp.CheckAndComplete(err)
		break
	default:
		printUsage()
	}
}
