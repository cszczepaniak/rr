package main

import (
	workoutservice "datastar/rr/service/workouts"
	"datastar/rr/ui"
	"datastar/rr/ui/components"
	"datastar/rr/ui/workouts"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("GET /web/", http.StripPrefix("/web", http.FileServer(http.Dir("web"))))

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		err := ui.Index2().Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	mux.HandleFunc("GET /timer", func(w http.ResponseWriter, r *http.Request) {
		err := components.Timer(components.TimerProps{
			Duration: 5 * time.Second,
			CountIn:  2 * time.Second,
		}).Render(r.Context(), w)
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
