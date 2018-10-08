package main

import (
	log "github.com/hoskeri/liverebuild/llog"
	"github.com/hoskeri/liverebuild/updater"
	"path"
	"time"
)

type FileSet struct {
	uf      updater.Updater
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

const BatchDuration = 750 * time.Millisecond

func (r *LiveRebuild) Watch() {
	var ticker = time.NewTicker(BatchDuration)

	/* Group updates by updater, and dispatch them all at once */
	type p_u []string
	type p_u_map map[updater.Updater]p_u
	var pending = make(p_u_map)

	for _, f := range r.fileSet {
		if err := r.watcher.Add(f.baseDir); err != nil {
			log.Debugf("failed to watch: %s:%s", f.baseDir, err)
		}
	}

	for {
		select {
		case e := <-r.watcher.Events:
			for _, fs := range r.fileSet {
				if fs.Match(e.Name) {
					log.Debugf("enqueue[%s] %s", fs.uf.Name(), e.Name)
					if p_u, ok := pending[fs.uf]; ok {
						pending[fs.uf] = append(p_u, e.Name)
					} else {
						pending[fs.uf] = []string{e.Name}
					}
				}
			}

		case e := <-r.watcher.Errors:
			log.Debugf("error %s", e)

		case _ = <-ticker.C:
			for u_f, p_u := range pending {
				log.Debugf("update[%s]: %q", u_f.Name(), p_u)
				// TODO: deduplicate the p_u list
				u_f.Update(BatchDuration, p_u...)
			}
			pending = make(p_u_map)
		}
	}
}
