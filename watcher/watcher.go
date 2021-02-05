package watcher

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/viknesh-nm/go-listener/builder"
	"github.com/viknesh-nm/go-listener/runner"
)

// Watcher holds the fields for the application that needs to be listened
type Watcher struct {
	builder   *builder.Builder
	runner    *runner.Runner
	dir       string
	buildOnly bool
}

// New create and return new watcher object.
func New(path, name string, buildOnly bool, commands []string) *Watcher {
	return &Watcher{
		dir:       path,
		buildOnly: buildOnly,
		builder:   builder.New(name, path),
		runner:    runner.New(name, path, commands),
	}
}

// Watch monitors the given path for changes
func (w *Watcher) Watch() error {
	log.Infof("Watching: %s", w.dir)

	w.run()

	stopWatch := make(chan error)
	go func() {
		ticker := time.NewTicker(400 * time.Millisecond)
		for {
			err := filepath.Walk(w.dir, w.watchFunc)
			if err != nil && err != filepath.SkipDir {
				stopWatch <- err
				break
			}
			<-ticker.C
		}
	}()

	return <-stopWatch
}

// watchFunc listens all the go files and rerun the application
func (w *Watcher) watchFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	// Skip all directories and files start with dot
	if strings.HasPrefix(filepath.Base(path), ".") {
		if info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}

	// Track all go files
	if filepath.Ext(path) == ".go" {
		// If any file gets modified after the last build, the build and run again.
		if info.ModTime().After(w.builder.GetLatestBuild()) {
			p, err := filepath.Rel(w.dir, path)
			if err != nil {
				return err
			}
			log.Infof("Modified: %s", p)
			w.run()
		}
	}
	return nil
}

// Run builds and runs watching project every time files get modified.
func (w *Watcher) run() {
	// build
	ok := w.builder.Build()
	if ok && !w.buildOnly {
		// run
		if err := w.runner.Run(); err != nil {
			log.Error(err)
		}
	}
}
