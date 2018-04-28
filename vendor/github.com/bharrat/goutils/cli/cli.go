package cli

import (
	"os"
	"os/exec"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

//Debug "Enable debug mode, that will print every command being executed"
var Debug = false

//DryRun "Enable Dry run mode, no command will actually be run"
var DryRun = false

//RunCommand Run a command and print output to stdout
func RunCommand(command string, args ...string) (exitCode int, err error) {

	if DryRun {
		log.Infof("[DRY RUN] %s %s", command, strings.Join(args, " "))
		exitCode = 0
		return
	} else if Debug {
		log.Infof("Running command: %s %s", command, strings.Join(args, " "))
	}

	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	exitCode = 1
	if err = cmd.Start(); err != nil {
		return exitCode, err
	}

	if err = cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		}
	} else {
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()

	}

	return
}

//"Run a shell command and return stdout, stderr and exitcode"
func GetOutputOfCommand(name string, args ...string) (output string, exitCode int, err error) {

	cmd := exec.Command(name, args...)
	out_bytes, err := cmd.CombinedOutput()

	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}
	output = string(out_bytes)
	return
}

//CommandExists Check whether a cli command exists or not, in DryRun mode always return True
func CommandExists(commandName string) (exists bool, err error) {

	if DryRun {
		exists = true
	} else {
		_, exitCode, e := GetOutputOfCommand("type", commandName)
		err = e
		if exitCode == 0 {
			exists = true
		}
	}

	return
}
