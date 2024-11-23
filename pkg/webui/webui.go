package webui

// Implements a web interface that makes it easier to categorize and manage wallpapers

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"yoink/pkg/log"
)

func Listen(ctx context.Context) {
	logger := log.Default()
	mux := http.NewServeMux()
	s := &http.Server{
		Addr: ":8081",
		Handler: mux,
	}

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
 		fmt.Fprintf(w, "Hello, Twitch! %s", r.URL.Path)
	})


	go func() {
		err := s.ListenAndServe()

		if err != nil {
			if strings.Contains(err.Error(), "Server closed"){
				logger.Warn("closing server")
			} else {
				logger.Error("error in webui", "error", err, "type", reflect.TypeOf(err))
			}
		}
	}()
	
	<-ctx.Done()
	s.Shutdown(ctx)
	logger.Info("Received ctx.Done, stopping WebUI")
}

