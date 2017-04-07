package main

import log "github.com/sirupsen/logrus"
import "flag"
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

func (r *LiveRebuild) Run() error {
	r.server = livereload.New("liverebuild")
	log.Errorln(r.watchFileSet.Rescan())
	log.Errorln(r.buildFileSet.Rescan())

	for {
		select {
		case event := <-r.watchFileSet.watcher.Events:
			if event.Op&(fsnotify.Rename|fsnotify.Create|fsnotify.Write) > 0 {
				log.Debug("running watchAction %s", r.buildAction)
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
	server             *livereload.Server
	buildAction        string
	buildFileSet       FileSet
	watchFileSet       FileSet
	watchServeRoot     string
	watchServeFallback string
}

func main() {
	service := new(LiveRebuild)

	flag.StringVar(&service.buildAction, "onbuild", "", "shell command on build file change")
	flag.StringVar(&service.buildFileSet.pattern, "buildfiles", "", "set of files to rebuild")
	flag.StringVar(&service.watchFileSet.pattern, "watchfiles", "", "set of files to livereload")
	flag.StringVar(&service.watchServeRoot, "root", "", "static server document root")
	flag.StringVar(&service.watchServeFallback, "fallback", "index.html", "path to render on fallback")

	log.SetLevel(log.DebugLevel)

	flag.Parse()

	log.Println("starting liverebuild")
	service.Run()
	log.Println("liverebuild is running")
}
