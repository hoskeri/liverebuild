package main

import (
	"flag"
	"net/http"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	log "github.com/hoskeri/liverebuild/llog"
	"github.com/hoskeri/liverebuild/updater"
	"github.com/rakyll/globalconf"
)

type LiveRebuild struct {
	listenStatic string
	listenLR     string
	staticMux    *http.ServeMux

	watcher *fsnotify.Watcher
	fileSet []*FileSet

	buildActionRoot string
	buildAction     string

	watchServeRoot     string
	watchServeFallback string
}

func (r *LiveRebuild) Run() (err error) {
	r.watcher, err = fsnotify.NewWatcher()

	go func() {
		var err = http.ListenAndServe(r.listenStatic, r.staticMux)
		if err != nil {
			log.Error(err)
		}
	}()

	r.staticMux = http.NewServeMux()
	r.staticMux.Handle("/", http.FileServer(http.Dir(r.watchServeRoot)))

	r.Watch()

	return
}

func main() {
	service := new(LiveRebuild)

	listenStatic := flag.String("listenstatic", ":4000", "static file listen address")
	listenLR := flag.String("listenlivereload", ":35729", "livereload listener address")

	buildActionRoot := flag.String("buildcommandroot", "", "base directory for build")
	buildAction := flag.String("buildcommand", "", "command to build")
	buildFiles := flag.String("buildfiles", "", "set of files to rebuild")

	watchServeRoot := flag.String("watchserveroot", "", "static server document root")
	watchFiles := flag.String("watchservefiles", "", "set of files to livereload")
	watchServeFallback := flag.String("watchservefallback", "index.html", "path to render on fallback")

	var opts = globalconf.Options{
		Filename:  ".liverebuildrc",
		EnvPrefix: "LIRB_",
	}

	var conf, err = globalconf.NewWithOptions(&opts)

	if err != nil {
		log.Fatal(err)
	}

	conf.ParseAll()

	var nothing = new(updater.Nothing)

	flag.VisitAll(func(f *flag.Flag) { log.Debug(f.Name, " -> ", f.Value) })

	service.listenStatic = *listenStatic
	service.listenLR = *listenLR

	service.buildActionRoot = *buildActionRoot
	service.buildAction = *buildAction

	for _, e := range strings.Split(*buildFiles, " ") {
		service.Add(e, nothing)
	}

	service.watchServeRoot = *watchServeRoot
	service.watchServeFallback = *watchServeFallback

	for _, e := range strings.Split(*watchFiles, " ") {
		service.Add(e, nothing)
	}

	log.Debug("starting liverebuild")
	err = service.Run()
	if err != nil {
		log.Fatalf("error: %s", err)
		os.Exit(1)
	}
}
