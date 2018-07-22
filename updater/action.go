package updater

import (
	livereload "github.com/omeid/go-livereload"
	"log"
	"net/http"
	"os/exec"
	"time"
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

func NewRunCommand(cmd string, args ...string) (*RunCommand, error) {
	return &RunCommand{
		cmd:  cmd,
		args: args,
	}, nil
}

func (u *RunCommand) Update(ts time.Duration, name string) {
	cmd := exec.Command(u.cmd, u.args...)
	if op, err := cmd.CombinedOutput(); err != nil {
		log.Printf("%s", string(op))
	}
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
	cmd  string
	args []string
	proc *exec.Cmd
}

func NewChildProcess(cmd string, args ...string) (*ChildProcess, error) {
	return &ChildProcess{
		cmd:  cmd,
		args: args,
	}, nil
}

func (u *ChildProcess) Update(ts time.Duration, name string) {
	if u.proc != nil && u.proc.Process != nil {
		u.proc.Process.Kill()
		log.Printf("child process exited: %+v", u.proc.Wait())
	}

	u.proc = exec.Command(u.cmd, u.args...)
	u.proc.Start()
}

type StaticServer struct {
	Updater
	server *http.Server
}

func NewStaticServer(address, dir, fallback string) (*StaticServer, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	})

	s := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	return &StaticServer{
		server: s,
	}, s.ListenAndServe()
}
