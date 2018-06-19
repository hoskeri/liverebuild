package main

import log "github.com/sirupsen/logrus"
import "flag"
import "net/http"
import "path"
import "github.com/omeid/go-livereload"
import "github.com/fsnotify/fsnotify"
import "github.com/rakyll/globalconf"
import "strings"
import "os"

type BuildAction interface {
	FileChanged(path string) error
}

type FileSet struct {
	name    string
	baseDir string
	match   string
	watcher *fsnotify.Watcher
}

func New(name, pat string) (f *FileSet) {
	f = &FileSet{
		name:    name,
		baseDir: path.Dir(pat),
		match:   path.Base(pat),
	}

	f.watcher, _ = fsnotify.NewWatcher()
	f.watcher.Add(f.baseDir)
	return
}

func (r *LiveRebuild) Run() (err error) {
	r.lrServer = livereload.New("liverebuild")

	r.lrMux = http.NewServeMux()
	r.lrMux.HandleFunc("/livereload.js", livereload.LivereloadScript)
	r.lrMux.Handle("/", r.lrServer)

	// FIXME: call fallback url on 404.
	//        set caching headers.
	//        call backend on matched requests.
	r.staticMux = http.NewServeMux()
	r.staticMux.Handle("/", http.FileServer(http.Dir(r.watchServeRoot)))

	go func() {
		var err = http.ListenAndServe(r.listenStatic, r.staticMux)
		if err != nil {
			log.Error(err)
		}
	}()

	go func() {
		var err = http.ListenAndServe(r.listenLR, r.lrMux)
		if err != nil {
			log.Error(err)
		}
	}()

	const interestedEvents = fsnotify.Create | fsnotify.Remove | fsnotify.Rename

	for {
		select {
		case event := <-r.fileSet[0].watcher.Events:
			if event.Op&interestedEvents > 0 {
				log.Debugf("reload file %s: %s", event.Name, event.Op)
				r.lrServer.Reload(event.Name, false)
			}
		case e := <-r.fileSet[0].watcher.Errors:
			log.Debugf("caught error %s", e)
		}
	}
}

type LiveRebuild struct {
	listenStatic string
	listenLR     string
	lrServer     *livereload.Server
	lrMux        *http.ServeMux
	staticMux    *http.ServeMux

	fileSet []*FileSet

	buildActionRoot string
	buildAction     string

	watchServeRoot     string
	watchServeFallback string
}

func main() {
	service := new(LiveRebuild)

	verbose := flag.Bool("verbose", false, "verbose logging")

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

	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	flag.VisitAll(func(f *flag.Flag) { log.Debugln(f.Name, "->", f.Value) })

	service.listenStatic = *listenStatic
	service.listenLR = *listenLR

	service.buildActionRoot = *buildActionRoot
	service.buildAction = *buildAction

	for _, e := range strings.Split(*buildFiles, " ") {
		service.fileSet = append(service.fileSet, New("build", e))
	}

	service.watchServeRoot = *watchServeRoot
	service.watchServeFallback = *watchServeFallback

	for _, e := range strings.Split(*watchFiles, " ") {
		service.fileSet = append(service.fileSet, New("watch", e))
	}

	log.Debugln("starting liverebuild")
	err = service.Run()
	if err != nil {
		log.Fatalf("error: %s", err)
		os.Exit(1)
	}
}
