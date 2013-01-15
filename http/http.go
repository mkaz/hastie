package http

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"github.com/mkaz/hastie"
	"net/http"
	"regexp"
	"sync"
	"time"
)

var (
	WatchFilesRE = regexp.MustCompile("\\.(html)|(md)$") // Watch all .html and .md files
)

// Handler is a struct that can serve HTTP requests to hastie templates and reload on demand.
type Handler struct {
	mx      sync.RWMutex
	config  hastie.Config
	monitor hastie.Monitor
	handler http.Handler
}

// Handle creates a Handler struct which implements http.Handler to serve hastie generated files
// and can reload on demand.
func Handle(config hastie.Config, monitor hastie.Monitor) *Handler {
	handler := &Handler{config: config, monitor: monitor}
	handler.Reload()
	return handler
}

// Reload compiles the hastie templates and updates the handler to serve them.
func (h *Handler) Reload() {
	h.mx.Lock()
	defer h.mx.Unlock()

	if err := h.config.Compile(h.monitor); err != nil {
		h.handler = errorHandler{err}
	} else {
		h.handler = http.FileServer(http.Dir(h.config.PublishDir))
	}
}

// ServeHTTP will serve a HTTP request to the hastie templates.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mx.RLock()
	defer h.mx.RUnlock()

	h.handler.ServeHTTP(w, r)
}

// Watch checks to see if there are any changes to the SourceDir or LayoutDir
// and if so triggers a Reload.
// This is blocking so you'll probably want to call from a go routine.
func (h *Handler) Watch() error {
	f := watch([]string{h.config.SourceDir, h.config.LayoutDir}, func() {
		h.Reload()
	})

	return f
}

func watch(dirs []string, reload func()) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer w.Close()

	errs := make(chan error)
	defer close(errs)

	go func() {
		tick := time.Tick(time.Second)
		doReload := false
		for {
			select {
			case e := <-w.Event:
				if WatchFilesRE.MatchString(e.Name) {
					doReload = true // trigger reload
				}
			case err := <-w.Error:
				errs <- err
				break
			case _ = <-tick:
				if doReload {
					reload()
					doReload = false // Reset
				}
			}
		}
	}()

	for _, d := range dirs {
		if err := w.Watch(d); err != nil {
			return err
		}
	}

	err = <-errs

	return err
}

type errorHandler struct {
	error
}

func (h errorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Server Error", 500)
	fmt.Fprintf(w, h.Error())
}
