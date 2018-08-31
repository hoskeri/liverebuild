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
	Root     string
	Paths    []string
	Fallback string
}

type lrAction struct {
	Paths   []string
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

func parseConfig(cf string) {
	if _, err := toml.DecodeFile(cf, &Config); err != nil {
		log.Fatalf("config error: %+v", err)
	} else {
		if Config.Verbose {
			log.Debug("parsed config: %+v", Config)
		}
	}
}

type updaters struct {
	static *updater.StaticServer
	lr     *updater.LiveReload
	daemon *updater.ChildProcess
	cmd    *updater.RunCommand
}

func main() {
	parseConfig(".liverebuildrc")

	service := new(LiveRebuild)
	up := new(updaters)
	var err error

	if len(Config.Static.Paths) > 0 {
		if up.static, err = updater.NewStaticServer(
			Config.Static.Address,
			Config.Static.Root,
			Config.Static.Fallback); err != nil {
			log.Fatalf("failed to initialize static server: %s", err)
		} else {
			for _, p := range Config.Static.Paths {
				service.Add(p, up.static)
			}
		}
	}

	if len(Config.Lr.Paths) > 0 {
		if up.lr, err = updater.NewLiveReload(
			Config.Lr.Address); err != nil {
			log.Fatalf("failed to initialize static server: %s", err)
		} else {
			for _, p := range Config.Lr.Paths {
				service.Add(p, up.lr)
			}
		}
	}

	if len(Config.Build.Paths) > 0 {
		if up.cmd, err = updater.NewRunCommand(Config.Build.Cmd); err != nil {
			log.Fatalf("failed to initialize static server: %s", err)
		} else {
			for _, p := range Config.Build.Paths {
				service.Add(p, up.cmd)
			}
		}
	}

	if len(Config.Daemon.Paths) > 0 {
		if up.daemon, err = updater.NewChildProcess(Config.Daemon.Cmd); err != nil {
			log.Fatalf("failed to initialize static server: %s", err)
		} else {
			for _, p := range Config.Daemon.Paths {
				service.Add(p, up.daemon)
			}
		}
	}

	if err := service.Run(); err != nil {
		log.Fatalf("error: %s", err)
		os.Exit(1)
	}
}
