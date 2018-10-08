package updater

import (
	log "github.com/hoskeri/liverebuild/llog"
	livereload "github.com/omeid/go-livereload"
	"net/http"
	"time"
)

type LiveReload struct {
	lr     *livereload.Server
	server *http.Server
}

var _ Updater = (*LiveReload)(nil)

func (l *LiveReload) Name() string { return "livereload" }

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

	go func() { lr.server.ListenAndServe() }()

	return &lr, nil
}

func (l *LiveReload) Update(ts time.Duration, name ...string) {
	for _, e := range name {
		l.lr.Reload(e, false)
	}
}

func init() {
	livereload.Log = log.Blackhole
}
