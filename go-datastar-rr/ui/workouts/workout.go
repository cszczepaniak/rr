package workouts

import (
	"datastar/rr/service/workouts"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	datastar "github.com/starfederation/datastar/sdk/go"
)

type Handler struct {
	service *workouts.Service
}

func NewHandler(service *workouts.Service) Handler {
	return Handler{
		service: service,
	}
}

func (h Handler) CreateWorkout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		workout, err := h.service.Start(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Location", fmt.Sprintf("/workouts/%s", workout))
		w.WriteHeader(http.StatusSeeOther)
	})
}

func (h Handler) GetWorkout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		stage, err := h.service.GetCurrentStage(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		viewData, err := stageToViewData(stage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		viewData.WorkoutID = id

		err = Workout(viewData).Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
}

func (w loggingResponseWriter) Write(bs []byte) (int, error) {
	// fmt.Print(string(bs))
	return w.ResponseWriter.Write(bs)
}

func (w loggingResponseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (h Handler) AdvanceWorkout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w = loggingResponseWriter{ResponseWriter: w}

		id := r.PathValue("id")

		signalData := struct {
			// This interval corresponds to a timer that might have been running on the previous
			// workout. We'll stop it if it's there.
			Interval int `json:"interval"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&signalData)
		if err != nil {
			slog.Error("advance", "jsonErr", err)
		}

		newStage, err := h.service.Advance(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		viewData, err := stageToViewData(newStage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		viewData.WorkoutID = id

		sse := datastar.NewSSE(w, r)

		// If there's an interval signal, it means there's a timer that was running on the previous
		// page. We'd like to stop that timer so it doesn't keep running when the workout loads!
		if signalData.Interval != 0 {
			script := fmt.Sprintf("window.clearInterval(%v)", signalData.Interval)
			sse.ExecuteScript(script)
		}

		err = sse.MergeFragmentTempl(body(viewData))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

type ViewData struct {
	WorkoutID       string
	CurrentMovement Movement
	IsResting       bool
	IsDone          bool
	RestDuration    time.Duration
}

type Movement struct {
	Category string
	Name     string
	Reps     int
	Dur      time.Duration
}

func stageToViewData(s workouts.Stage) (ViewData, error) {
	viewData := ViewData{
		CurrentMovement: Movement{
			Category: s.Category(),
			Name:     s.Name(),
		},
	}

	switch s := s.(type) {
	case workouts.End:
		return ViewData{
			IsDone: true,
		}, nil
	case workouts.Hold:
		viewData.CurrentMovement.Dur = s.Duration
	case workouts.Reps:
		viewData.CurrentMovement.Reps = s.Reps
	case workouts.Rest:
		viewData.IsResting = true
		viewData.RestDuration = s.Duration
	default:
		return ViewData{}, fmt.Errorf("unexpected workouts.Stage: %#v", s)
	}

	return viewData, nil
}
