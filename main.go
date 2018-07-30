package main

import (
	"os"

	"github.com/fsnotify/fsnotify"
	log "github.com/hoskeri/liverebuild/llog"
	"github.com/hoskeri/liverebuild/updater"
)

type LiveRebuild struct {
	watcher *fsnotify.Watcher
	fileSet []*FileSet
}

func (r *LiveRebuild) Run() (err error) {
	r.watcher, err = fsnotify.NewWatcher()

	r.Watch()
	return
}

func main() {
	service := new(LiveRebuild)

	var _ = updater.Nothing{}

	if err := service.Run(); err != nil {
		log.Fatalf("error: %s", err)
		os.Exit(1)
	}
}
