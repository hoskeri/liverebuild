package updater

import (
	log "github.com/hoskeri/liverebuild/llog"
	"os/exec"
	"time"
)

type RunCommand struct {
	cmd  string
	args []string
}

var _ Updater = (*RunCommand)(nil)

func (r *RunCommand) Name() string { return "runc" }

func NewRunCommand(cmd string, args ...string) (*RunCommand, error) {
	return &RunCommand{
		cmd:  cmd,
		args: args,
	}, nil
}

func (u *RunCommand) Update(ts time.Duration, name ...string) {
	cmd := exec.Command(u.cmd, u.args...)
	if op, err := cmd.CombinedOutput(); err != nil {
		log.Debug("%s", string(op))
	}
}
