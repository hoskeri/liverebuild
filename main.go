package main

import "log"
import "flag"
import "path/filepath"
import "github.com/omeid/go-livereload"
import "gopkg.in/fsnotify.v1"

type FileSet struct {
	pattern string
	matched []string
	watcher *fsnotify.Watcher
}

type LiveRebuild struct {
	server             *livereload.Server
	buildAction        string
	buildFileSet       FileSet
	watchFileSet       FileSet
	watchServeRoot     string
	watchServeFallback string
}

func (f *FileSet) Rescan() (err error) {
	f.matched, err = filepath.Glob(f.pattern)
	return
}

func (f *FileSet) Watch() (err error) {
	return
}

func (r *LiveRebuild) Run() error {
	r.server = livereload.New("liverebuild")
	r.watchFileSet.Watch()
	r.buildFileSet.Watch()

	return nil
}

func main() {
	service := new(LiveRebuild)

	flag.StringVar(&service.buildAction, "onbuild", "", "shell command on build file change")
	flag.StringVar(&service.buildFileSet.pattern, "buildfiles", "", "set of files to rebuild")
	flag.StringVar(&service.watchFileSet.pattern, "watchfiles", "", "set of files to livereload")
	flag.StringVar(&service.watchServeRoot, "root", "", "set of files to livereload")
	flag.StringVar(&service.watchServeFallback, "fallback", "index.html", "set of files to livereload")

	flag.Parse()
	log.Println("starting liverebuild")

	service.Run()

	log.Println(service)
}
