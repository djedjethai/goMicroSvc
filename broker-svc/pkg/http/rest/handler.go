package rest

import (
	"encoding/json"
	"errors"
	// "github.com/djedjethai/broker/pkg/dto"
	"net/http"

	"github.com/djedjethai/broker/pkg/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// pass svc as arg in the handler()
func Handler(js service.Service) http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"PUT", "POST", "GET", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// will allow us to ping the service
	mux.Use(middleware.Heartbeat("/ping"))

	mux.Get("/", getResponse(js))
	mux.Post("/", addResponse(js))

	// single entry point will handle all req
	mux.Post("/handle", handleSubmission(js))

	return mux
}

func handleSubmission(js service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestPayload service.RequestPayload

		err := readJson(w, r, &requestPayload)
		if err != nil {
			errorJson(w, err)
			return
		}

		// convert type from requestPayload.Auth to service.AuthPayload

		switch requestPayload.Action {
		case "auth":
			err, status, payload := js.Authenticate(requestPayload.Auth)
			if err != nil {
				errorJson(w, err, status)
			}

			writeJson(w, status, payload)
		case "log":
			err, status, payload := js.LogItem(requestPayload.Log)
			if err != nil {
				errorJson(w, err, status)
			}

			writeJson(w, status, payload)
		case "mail":
			err, status, payload := js.SendMail(requestPayload.Mail)
			if err != nil {
				errorJson(w, err, status)
			}

			writeJson(w, status, payload)
		default:
			errorJson(w, errors.New("Unknow action"))
		}
	}
}

// get []JsonResponse from memory
func getResponse(js service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := js.GetResponse()
		_ = writeJson(w, http.StatusOK, payload)
	}
}

// add JsonResponse
func addResponse(js service.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var jr service.JsonResponse
		_ = json.NewDecoder(r.Body).Decode(&jr)
		resp := js.AddData(jr)

		_ = writeJson(w, http.StatusAccepted, resp)
	}
}
