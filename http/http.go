package http

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"github.com/mkaz/hastie"
	"net/http"
	"sync"
)

type Handler struct {
	mx      sync.RWMutex
	config  hastie.Config
	monitor hastie.Monitor
	handler http.Handler
}

// Creates an http.Handler for hastie files.
func Handle(config hastie.Config, monitor hastie.Monitor) *Handler {
	handler := &Handler{config: config, monitor: monitor}
	handler.Reload()
	return handler
}

// Reload compiles the hastie files and updates the handler.
func (h *Handler) Reload() {
	h.mx.Lock()
	defer h.mx.Unlock()

	if err := h.config.Compile(h.monitor); err != nil {
		h.handler = errorHandler{err}
	} else {
		h.handler = http.FileServer(http.Dir(h.config.PublishDir))
	}
}

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
		for {
			select {
			case _ = <-w.Event:
				reload()
			case err := <-w.Error:
				errs <- err
				break
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
