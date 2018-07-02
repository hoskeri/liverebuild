package main

import log "github.com/sirupsen/logrus"
import "github.com/omeid/go-livereload"
import "github.com/fsnotify/fsnotify"
import "github.com/rakyll/globalconf"
import "flag"
import "net/http"
import "path"
import "strings"
import "os"
import "os/exec"
import "time"

type BuildAction interface {
	FileChanged(path string) error
}

type FileSet struct {
	key     string
	baseDir string
	match   string
}

func (fs *FileSet) Match(name string) bool {
	m, _ := path.Match(path.Join(fs.baseDir, fs.match), name)
	return m
}

func NewFileSet(key, pat string) (f *FileSet) {
	f = &FileSet{
		key:     key,
		baseDir: path.Dir(pat),
		match:   path.Base(pat),
	}

	return
}

func (r *LiveRebuild) Watch() {
	for _, f := range r.fileSet {
		if err := r.watcher.Add(f.baseDir); err != nil {
			log.Debugf("failed to watch: %s:%s", f.baseDir, err)
		}
	}

	for {
		select {
		case e := <-r.watcher.Events:
			log.Debugf("event[%s] %s", e.Op, e.Name)
			for _, fs := range r.fileSet {
				if fs.Match(e.Name) {
					switch fs.key {
					case "reload":
						r.lrServer.Reload(e.Name, false)
					case "build":
						go func() {
							cmd := exec.Command(r.buildAction)
							op, excode := cmd.CombinedOutput()
							log.Debugf("build command result: %s", op, excode)
						}()
					default:
						log.Debugf("unknown key: %s", fs.key)
					}
				}

			}
		case e := <-r.watcher.Errors:
			log.Debugf("error %s", e)
		}
	}
}

type LiveRebuild struct {
	listenStatic string
	listenLR     string
	lrServer     *livereload.Server
	lrMux        *http.ServeMux
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

	r.lrServer = livereload.New("liverebuild")
	r.lrMux = http.NewServeMux()
	r.lrMux.HandleFunc("/livereload.js", livereload.LivereloadScript)
	r.lrMux.Handle("/", r.lrServer)

	go func() {
		var err = http.ListenAndServe(r.listenStatic, r.staticMux)
		if err != nil {
			log.Error(err)
		}
	}()

	r.staticMux = http.NewServeMux()
	r.staticMux.Handle("/", http.FileServer(http.Dir(r.watchServeRoot)))

	go func() {
		var err = http.ListenAndServe(r.listenLR, r.lrMux)
		if err != nil {
			log.Error(err)
		}
	}()

	r.Watch()

	return
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
		service.fileSet = append(service.fileSet, NewFileSet("build", e))
	}

	service.watchServeRoot = *watchServeRoot
	service.watchServeFallback = *watchServeFallback

	for _, e := range strings.Split(*watchFiles, " ") {
		service.fileSet = append(service.fileSet, NewFileSet("watch", e))
	}

	log.Debugln("starting liverebuild")
	err = service.Run()
	if err != nil {
		log.Fatalf("error: %s", err)
		os.Exit(1)
	}
}
