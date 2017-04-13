package main

import log "github.com/sirupsen/logrus"
import "flag"
import "net/http"
import "path/filepath"
import "github.com/omeid/go-livereload"
import "gopkg.in/fsnotify.v1"

type BuildAction interface {
	FileChanged(path string) error
}

type FileSet struct {
	pattern string
	matched []string
	watcher *fsnotify.Watcher
}

func (f *FileSet) Rescan() (err error) {
	f.matched, err = filepath.Glob(f.pattern)
	if err != nil {
		return err
	}

	f.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	for _, e := range f.matched {
		log.Debugf("watching file %s", e)
		err = f.watcher.Add(e)
		if err != nil {
			f.watcher.Close()
			return err
		}
	}
	return err
}

type wrappedWriter struct {
	w http.ResponseWriter
	h http.Header
}

func (r *LiveRebuild) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.staticMux.ServeHTTP(w, req)
}

func (r *LiveRebuild) Run() error {
	r.lrServer = livereload.New("liverebuild")

	r.lrMux = http.NewServeMux()
	r.lrMux.HandleFunc("/livereload.js", livereload.LivereloadScript)
	r.lrMux.Handle("/", r.lrServer)

	r.staticMux.Handle("/", http.FileServer(http.Dir(r.watchServeRoot)))

	r.buildFileSet.Rescan()
	r.watchFileSet.Rescan()

	for {
		select {
		case event := <-r.watchFileSet.watcher.Events:
			if event.Op&(fsnotify.Rename|fsnotify.Create|fsnotify.Write) > 0 {
				log.Debug("reload file %s", event.Name)
			}
		case event := <-r.buildFileSet.watcher.Events:
			if event.Op&(fsnotify.Rename|fsnotify.Create|fsnotify.Write) > 0 {
				log.Debug("running buildAction %s", r.buildAction)
			}

		case e := <-r.buildFileSet.watcher.Errors:
			log.Debugf("caught error %s", e)

		case e := <-r.watchFileSet.watcher.Errors:
			log.Debugf("caught error %s", e)
		}
	}
	return nil
}

type LiveRebuild struct {
	lrServer           *livereload.Server
	lrMux              *http.ServeMux
	staticMux          *http.ServeMux
	buildAction        string
	buildFileSet       FileSet
	watchFileSet       FileSet
	watchServeRoot     string
	watchServeFallback string
}

func main() {
	var verbose = false
	service := new(LiveRebuild)

	flag.StringVar(&service.buildAction, "onbuild", "", "shell command on build file change")
	flag.StringVar(&service.buildFileSet.pattern, "buildfiles", "", "set of files to rebuild")
	flag.StringVar(&service.watchFileSet.pattern, "watchfiles", "", "set of files to livereload")
	flag.StringVar(&service.watchServeRoot, "root", "", "static server document root")
	flag.StringVar(&service.watchServeFallback, "fallback", "index.html", "path to render on fallback")
	flag.BoolVar(&verbose, "verbose", false, "verbose logging")

	flag.Parse()

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugln("starting liverebuild")
	service.Run()
}
