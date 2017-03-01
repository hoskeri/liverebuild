package main

import "log"
import "flag"
import "github.com/omeid/go-livereload"

type FileSet struct {
	pattern string
	matched string
}

type LiveRebuild struct {
	server       *livereload.Server
	buildAction  string
	buildFileSet FileSet
	watchFileSet FileSet
}

func main() {
	service := new(LiveRebuild)

	flag.StringVar(&service.buildAction, "onbuild", "", "shell command on build file change")
	flag.StringVar(&service.buildFileSet.pattern, "buildfiles", "", "set of files to rebuild")
	flag.StringVar(&service.watchFileSet.pattern, "watchfiles", "", "set of files to livereload")
	log.Println("starting liverebuild")
	log.Println(service)
}
