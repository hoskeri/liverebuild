package updater

import (
	log "github.com/hoskeri/liverebuild/llog"
	"os/exec"
	"time"
)

type ChildProcess struct {
	Updater
	cmd  string
	args []string
	proc *exec.Cmd
}

var _ Updater = (*ChildProcess)(nil)

func (c *ChildProcess) Name() string { return "childprocess" }

func NewChildProcess(cmd string, args ...string) (*ChildProcess, error) {
	return &ChildProcess{
		cmd:  cmd,
		args: args,
	}, nil
}

func (u *ChildProcess) Update(ts time.Duration, name ...string) {
	if u.proc != nil && u.proc.Process != nil {
		u.proc.Process.Kill()
		log.Debug("child process exited: %+v", u.proc.Wait())
	}

	u.proc = exec.Command(u.cmd, u.args...)
	u.proc.Start()
}
