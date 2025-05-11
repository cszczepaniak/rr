package main

import (
	workoutservice "datastar/rr/service/workouts"
	"datastar/rr/ui"
	"datastar/rr/ui/workouts"
	"embed"
	"log/slog"
	"net/http"
	"os"
)

//go:embed web/*
var webFS embed.FS

func main() {
	mux := http.NewServeMux()

	if os.Getenv("RAILWAY_SERVICE_ID") != "" {
		// If deployed, use the embedded filesystem; we don't do this for development because it's
		// easier to get live updates from the actual directory.
		//
		// We don't need to strip the prefix because the embed.FS is relative to the root
		// directory.
		mux.Handle("GET /web/", http.FileServer(http.FS(webFS)))
		slog.Info("serving embedded static assets")
	} else {
		mux.Handle("GET /web/", http.StripPrefix("/web", http.FileServer(http.Dir("web"))))
		slog.Info("serving static assets from local file system")
	}

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		err := ui.Index().Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	hs := workoutservice.New()
	h := workouts.NewHandler(hs)
	mux.Handle("POST /workouts", h.CreateWorkout())
	mux.Handle("POST /workouts/{id}/advance", h.AdvanceWorkout())

	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error("Error starting server:", slog.String("error", err.Error()))
	}
}
