package main

import (
	"os"

	"github.com/BurntSushi/toml"
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

type staticAction struct {
	Address  string
	Paths    []string
	Fallback string
}

type lrAction struct {
	Address string
}

type buildAction struct {
	Paths []string
	Cmd   string
}

type daemonAction struct {
	Paths []string
	Cmd   string
}

type config struct {
	Verbose bool
	Static  staticAction
	Lr      lrAction
	Build   buildAction
	Daemon  daemonAction
}

var Config = config{}

func parseConfig(cf string) error {
	_, err := toml.DecodeFile(cf, &Config)
	return err
}

func main() {
	if err := parseConfig(".liverebuildrc"); err != nil {
		log.Fatalf("failed to parse .liverebuildrc: %+v", err)
	} else {
		log.Debugf("parsed config: %+v", Config)
	}

	service := new(LiveRebuild)

	var _ = updater.Nothing{}

	if err := service.Run(); err != nil {
		log.Fatalf("error: %s", err)
		os.Exit(1)
	}
}
