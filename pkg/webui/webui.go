package webui

// Implements a web interface that makes it easier to categorize and manage wallpapers

import (
	_ "embed"

	"context"
	"html/template"
	"net/http"
	"os"
	"path"
	"reflect"
	"strings"

	"yoink/pkg/config"
	"yoink/pkg/log"
)

////go:embed template.html
//var indexHtml string
// var imageDir = path.Join(config.DataPath(), "new")

// var keepDir = path.Join(config.DataPath(), "categorised", "keep")
// var discardDir = path.Join(config.DataPath(), "categorised", "discard")
// var animeDir = path.Join(config.DataPath(), "categorised", "anime-nsfw")
// var nsfwDir = path.Join(config.DataPath(), "categorised", "nsfw")

func Listen(ctx context.Context) {
	logger := log.Default()
	mux := http.NewServeMux()
	s := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	mux.Handle("GET /img/", http.StripPrefix("/img/", http.FileServer(http.Dir(config.NewDir))))

	mux.HandleFunc("GET /", handleIndex)
	mux.HandleFunc("POST /{imageId}/{verb}", handlePost)

	go func() {
		err := s.ListenAndServe()

		if err != nil {
			if strings.Contains(err.Error(), "Server closed") {
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

func handleIndex(w http.ResponseWriter, r *http.Request) {
	logger := log.Default()

	files, err := os.ReadDir(config.NewDir)
	if err != nil {
		logger.Error("Could not read data directory", "error", err)
		return
	}

	// get the first file that's an image
	var imageFile os.DirEntry
	for _, file := range files {
		e := path.Ext(file.Name())
		if e == ".jpg" || e == ".png" || e == ".jpeg" {
			imageFile = file
			break
		}
	}

	// for development, read the ./template.html file and use it as the template
	cwd, err := os.Getwd()
	if err != nil {
		logger.Error("Could not get working directory", "error", err)
		return
	}

	indexHtml, err := os.ReadFile(cwd + "/pkg/webui/template.html")
	if err != nil {
		logger.Error("Could not read template.html", "error", err)
		return
	}

	data := struct {
		Image string
		Count int
	}{
		Image: imageFile.Name(),
		Count: len(files),
	}

	tmpl, err := template.New("index").Parse(string(indexHtml))
	if err != nil {
		logger.Error("Could not parse template", "error", err)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		logger.Error("Could not execute template", "error", err)
		return
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	logger := log.Default()

	imageId := r.PathValue("imageId")
	verb := r.PathValue("verb")
	logger.Info("Received POST request", "imageId", imageId, "verb", verb)

	verbMapping := map[string]string{
		"keep":    config.KeepDir,
		"discard": config.DiscardDir,
		"anime":   config.AnimeDir,
		"nsfw":    config.NsfwDir,
	}

	if _, ok := verbMapping[verb]; !ok {
		logger.Warn("Unknown verb", "verb", verb)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := os.Rename(path.Join(config.NewDir, imageId), path.Join(verbMapping[verb], imageId))
	if err != nil {
		logger.Error("Could not move file", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
