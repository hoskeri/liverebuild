package updater

import (
	"os/exec"
	"time"

	livereload "github.com/omeid/go-livereload"
)

type UpdateFunc func(time.Duration, string)

type Updater interface {
	Update(time.Duration, string)
}

type Nothing struct {
	Updater
}

func (u *Nothing) Update(ts time.Duration, name string) {
	return
}

type RunCommand struct {
	cmd  string
	args []string
}

func (u *RunCommand) Update(ts time.Duration, name string) {
	cmd := exec.Command(u.cmd, u.args...)
	cmd.Run()
}

type LiveReload struct {
	Updater
	lr *livereload.Server
}

type ChildProcess struct {
	Updater
	oldProc *exec.Cmd
	proc    *exec.Cmd
}

func (u *ChildProcess) Update(ts time.Duration, name string) {
}
