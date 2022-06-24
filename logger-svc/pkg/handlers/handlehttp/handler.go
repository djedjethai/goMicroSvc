package handlehttp

import (
	"github.com/djedjethai/logger/pkg/service"
	"github.com/djedjethai/toolbox"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

var tools = new(toolbox.Tools)

func Handler(s service.Service) http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"PUT", "POST", "GET", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// set routes
	mux.Use(middleware.Heartbeat("/ping"))
	mux.Post("/log", writeLog(s))

	return mux
}

func writeLog(s service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestPayload service.JSONPayload
		_ = tools.ReadJSON(w, r, &requestPayload)

		resp, err := s.InsertPayload(requestPayload, tools)
		if err != nil {
			tools.ErrorJSON(w, err)
			return
		}

		tools.WriteJSON(w, http.StatusAccepted, resp)
	}
}
