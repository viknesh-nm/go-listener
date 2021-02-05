package runner

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// Runner holds the field that needs to listen while running the applications
type Runner struct {
	cmd      *exec.Cmd
	path     string
	commands []string
}

// New returns new runner object
func New(name, path string, commands []string) *Runner {
	return &Runner{
		path:     filepath.Join(path, name),
		commands: commands,
	}
}

// Run our application
func (r *Runner) Run() error {

	// Terminate running process
	err := r.terminate()
	if err != nil {
		return err
	}
	log.Infof("Running: %s %s", filepath.Base(r.path), strings.Join(r.commands, " "))

	// Run the application
	r.cmd = exec.Command(r.path, r.commands...)
	r.cmd.Stderr = os.Stderr
	r.cmd.Stdout = os.Stdout
	return r.cmd.Start()
}

// terminate stop running process if anyone has a pending status
func (r *Runner) terminate() error {
	if r.cmd == nil {
		return nil
	}

	timer := time.NewTicker(time.Second)
	done := make(chan error, 1)

	// Try to stop the running processes
	go func() {
		done <- r.cmd.Wait()
	}()

	// Wait for atleast one second for the process to stop
	select {
	case <-timer.C:
		r.cmd.Process.Kill()
		<-done
	case err := <-done:
		if err != nil {
			if _, ok := err.(*exec.ExitError); !ok {
				return err
			}
		}
	}
	r.cmd = nil
	return nil
}
