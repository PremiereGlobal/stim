package shell

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"syscall"
)

type ShellCommand struct {
	Shell   []string
	Envs    []string
	Command []string
}

// Run runs a shell command and returns the output
// NEEDS WORk to handle output/errors better
func Run(shellCommand ShellCommand) (string, error) {

	// Set default shell
	if len(shellCommand.Shell) == 0 {
		shellCommand.Shell = []string{"/bin/sh", "-c"}
	}

	if len(shellCommand.Command) == 0 {
		return "", errors.New("No shell command specified")
	}

	fullCommand := append(shellCommand.Shell, shellCommand.Command...)
	cmd := exec.Command(fullCommand[0], fullCommand[1:]...)
	cmd.Env = shellCommand.Envs

	// Capture stdout messages
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error creating stdout pipe. %v", err))
	}

	// Capture stderr messages
	// stderr, err := cmd.StderrPipe()
	// if err != nil {
	//   return "", errors.New(fmt.Sprintf("Error creating stderr pipe. %v", err))
	// }

	// Run the command (async)
	if err := cmd.Start(); err != nil {
		return "", errors.New(fmt.Sprintf("Error starting command %v", err))
	}

	// Output stdout
	out, _ := ioutil.ReadAll(stdout)

	// Output stderr
	// scanner := bufio.NewScanner(stderr)
	// for scanner.Scan() {
	//   d.log.Warn(scanner.Text())
	// }

	// Wait for command to finish
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return "", errors.New(fmt.Sprintf("Shell command exit with code %d. %v", status.ExitStatus(), err))
			}
		}
	}

	return string(out), nil
}
