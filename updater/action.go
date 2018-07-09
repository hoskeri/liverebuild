package updater

import (
	livereload "github.com/omeid/go-livereload"
	"net/http"
	"os/exec"
	"time"
)

type UpdateFunc func(time.Duration, string)

type Updater interface {
	Update(time.Duration, string)
	Start() error
	Stop() error
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
	lr     *livereload.Server
	server *http.Server
}

func NewLiveReload(address string) (*LiveReload, error) {
	mux := http.NewServeMux()

	lr := LiveReload{
		lr: livereload.New("liverebuild"),
		server: &http.Server{
			Addr:    address,
			Handler: mux,
		},
	}

	mux.HandleFunc("/livereload.js", livereload.LivereloadScript)
	mux.Handle("/", lr.lr)

	return &lr, lr.server.ListenAndServe()
}

type ChildProcess struct {
	Updater
	oldProc *exec.Cmd
	proc    *exec.Cmd
}

func (u *ChildProcess) Update(ts time.Duration, name string) {
}
