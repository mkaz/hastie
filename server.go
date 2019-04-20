package main

import "fmt"
import "net/http"
import "time"

import "github.com/radovskyb/watcher"

func (h *Hastie) liveReload(log Logger, port int) {

	w := watcher.New()

	if err := w.AddRecursive(h.config.SourceDir); err != nil {
		log.Fatal(err)
	}
	if err := w.AddRecursive(h.config.LayoutDir); err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event := <-w.Event:
				log.Warn(fmt.Sprintf("File changed: %v", event))
				log.Warn("Regenerating content ...")

				// Note: Regenerating all content, keeping it simple
				h.generate()

				log.Warn("Done")
			case err := <-w.Error:
				log.Fatal(err)
			case <-w.Closed:
				return
			}
		}
	}()

	go func() {
		if err := w.Start(time.Millisecond * 100); err != nil {
			log.Fatal(err)
		}
	}()

	log.Warn("Watching Hastie content.")

	log.Warn(fmt.Sprintf("Listening on http://localhost:%d...\n", port))

	// Start the web server and listens until killed

	fs := http.FileServer(http.Dir(h.config.PublishDir))
	http.Handle("/", fs)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
