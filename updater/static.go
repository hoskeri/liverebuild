package updater

import (
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type StaticServer struct {
	server *http.Server
}

var _ Updater = (*StaticServer)(nil)

func (s *StaticServer) Name() string { return "static" }

func (s *StaticServer) Update(ts time.Duration, path ...string) {
	return
}

func NewStaticServer(address, dir, fallback string) (*StaticServer, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var lp = filepath.Join(dir, filepath.Clean(r.URL.EscapedPath()))
		if _, err := os.Stat(lp); err == nil {
			http.ServeFile(w, r, lp)
		} else {
			http.ServeFile(w, r, filepath.Join(dir, fallback))
		}
	})

	s := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	go func() { s.ListenAndServe() }()

	return &StaticServer{
		server: s,
	}, nil
}
