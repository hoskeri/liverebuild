package main

import "log"
import "flag"
import "path/filepath"
import "github.com/omeid/go-livereload"

type FileSet struct {
	pattern string
	matched []string
}

type LiveRebuild struct {
	server             *livereload.Server
	buildAction        string
	buildFileSet       FileSet
	watchFileSet       FileSet
	watchServeRoot     string
	watchServeFallback string
}

// XXX: pattern is not recursive.
// XXX: initial path is not validated.
func (f *FileSet) Rescan() (err error) {
	f.matched, err = filepath.Glob(f.pattern)
	return err
}

func (r *LiveRebuild) Setup() error {
	r.server = livereload.New("liverebuild")
	r.watchFileSet.Rescan()
	r.buildFileSet.Rescan()

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

	service.Setup()

	log.Println(service)
}
