package main

import (
	log "github.com/hoskeri/liverebuild/llog"
	"github.com/hoskeri/liverebuild/updater"
	"path"
	"time"
)

type FileSet struct {
	uf      updater.Updater
	last    time.Time
	baseDir string
	match   string
}

func (fs *FileSet) Match(name string) bool {
	m, _ := path.Match(path.Join(fs.baseDir, fs.match), name)
	return m
}

func (r *LiveRebuild) Add(pat string, fn updater.Updater) {
	f := &FileSet{
		uf:      fn,
		baseDir: path.Dir(pat),
		match:   path.Base(pat),
	}

	r.fileSet = append(r.fileSet, f)

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
					fs.uf.Update(time.Since(fs.last), e.Name)
					fs.last = time.Now()
				}
			}
		case e := <-r.watcher.Errors:
			log.Debugf("error %s", e)
		}
	}
}
