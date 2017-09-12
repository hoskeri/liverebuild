package main

import log "github.com/sirupsen/logrus"
import "flag"
import "net/http"
import "path/filepath"
import "github.com/omeid/go-livereload"
import "gopkg.in/fsnotify.v1"
import "github.com/rakyll/globalconf"
import "strings"
import "os"

type BuildAction interface {
	FileChanged(path string) error
}

type FileSet struct {
	name    string
	baseDir string
	pattern []string
	matched []string
	watcher *fsnotify.Watcher
}

func rescanDir(curdir string, pattern []string) (m []string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	for _, p := range pattern {
		rp := filepath.Join(wd, curdir, p)
		g, err2 := filepath.Glob(rp)
		if err != nil {
			return nil, err2
		}

		for _, e := range g {
			var fi, err = os.Stat(e)
			switch {
			case err != nil:
				return nil, err
			case fi.Mode().IsRegular():
				m = append(m, e)
			}
		}
	}

	return
}

func (f *FileSet) Rescan() (err error) {
	log.Debugf("matching pattern %s", f.pattern)

	f.matched, err = rescanDir(f.baseDir, f.pattern)
	if err != nil {
		return err
	}

	log.Debugln(f.matched)

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

func (r *LiveRebuild) Run() (err error) {
	r.lrServer = livereload.New("liverebuild")

	r.lrMux = http.NewServeMux()
	r.lrMux.HandleFunc("/livereload.js", livereload.LivereloadScript)
	r.lrMux.Handle("/", r.lrServer)

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

	err = r.buildFileSet.Rescan()
	if err != nil {
		log.Error(err)
		return err
	}
	err = r.watchFileSet.Rescan()
	if err != nil {
		log.Error(err)
		return err
	}

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
	listenStatic string
	listenLR     string
	lrServer     *livereload.Server
	lrMux        *http.ServeMux
	staticMux    *http.ServeMux

	buildActionRoot string
	buildAction     string
	buildFileSet    FileSet

	watchServeRoot     string
	watchFileSet       FileSet
	watchServeFallback string
}

func main() {
	service := new(LiveRebuild)

	listenStatic := flag.String("listenstatic", ":4000", "shell command on build file change")
	listenLR := flag.String("listenlivereload", ":35729", "shell command on build file change")
	verbose := flag.Bool("verbose", false, "verbose logging")

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
	service.buildFileSet.pattern = strings.Split(*buildFiles, " ")
	service.buildFileSet.baseDir = *buildActionRoot

	service.watchServeRoot = *watchServeRoot
	service.watchServeFallback = *watchServeFallback
	service.watchFileSet.pattern = strings.Split(*watchFiles, " ")
	service.watchFileSet.baseDir = *watchServeRoot

	log.Debugln("starting liverebuild")
	err = service.Run()
	if err != nil {
		log.Fatalf("error: %s", err)
		os.Exit(1)
	}
}
