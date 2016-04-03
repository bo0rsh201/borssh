package main
import (
	"os"
	"fmt"
	"syscall"
)
const COMMAND_COMPILE = "compile"
const COMMAND_CONNECT = "connect"

func failOnErr(err error) {
	if err != nil {
		fatal(err.Error())
	}
}

func fatal(msg string) {
	fmt.Fprintf(os.Stderr, msg + "\n")
	os.Exit(1)
}

func printUsage() {
	usage := "Usage: borssh COMMAND [arg...]\n\n"
	usage += "A simple ssh wrapper that helps you to keep all your dot files up to date\n\n"
	usage += "Commands:\n"
	usage += "\t%s\t\t\tCompiles dot files for latest config version\n"
	usage += "\t%s <target>\tConnect with latest compiled version\n\n"
	usage += "More info at http://github.com/bo0rsh201/borssh\n"
	fmt.Fprintf(os.Stderr, usage, COMMAND_COMPILE, COMMAND_CONNECT)
	os.Exit(1)
}

func main()  {
	if len(os.Args) < 2 {
		printUsage()
	}
	switch os.Args[1] {
	case COMMAND_CONNECT:
		if len(os.Args) < 3 {
			printUsage()
		}
		c, err := NewCompiler()
		failOnErr(err)
		localHash, err := c.getLocalHash()
		if os.IsNotExist(err) {
			fatal("Cannot find compiled version of dot files\nYou should run compile command first")
		}
		ex, err := NewExecutor(os.Args[2])
		failOnErr(err)
		remoteHash, err := c.getRemoteHash(ex)
		failOnErr(err)
		if remoteHash != localHash {
			cmd, err := ex.rsync(
				c.ph.getLocalBaseDir(),
				c.ph.getRemoteHomePath(),
				"--delete",
				"--copy-unsafe-links",
				"--exclude",
				COMPILED_HASH_FILE,
			)
			failOnErr(err)
			out, err := cmd.CombinedOutput()
			if err != nil {
				fatal(fmt.Sprint("Failed to rsync base dir: ", err.Error(), " output: ", string(out)))
			}
			err = c.install(
				ex,
				c.ph.getRemoteBashProfilePath(),
				c.ph.getRemoteCompiledBashProfilePath(),
			)
			failOnErr(err)

			cmd, err = ex.rsync(
				c.ph.getLocalHashPath(),
				c.ph.getRemoteHashPath(),
			)
			failOnErr(err)
			out, err = cmd.CombinedOutput()
			if err != nil {
				fatal(fmt.Sprint("Failed to rsync compiled hash: ", err.Error(), " output: ", string(out)))
			}
		}
		err = syscall.Exec(ex.sshBinary, []string{"-t", ex.host}, os.Environ())
		failOnErr(err)
		break
	case COMMAND_COMPILE:
		fmt.Print("Compiling...\n")
		c, err := NewCompiler()
		failOnErr(err)
		err = c.compile()
		failOnErr(err)
		fmt.Print("Done\nInstalling...\n")
		err = c.install(
			NewLocalExecutor(),
			c.ph.getLocalBashProfilePath(),
			c.ph.getLocalCompiledBashProfilePath(),
		)
		failOnErr(err)
		fmt.Print("Done\n")
		break
	default:
		printUsage()
	}
}