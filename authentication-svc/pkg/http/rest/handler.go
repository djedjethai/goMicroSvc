package rest

import (
	"errors"
	"fmt"
	"github.com/djedjethai/authentication/pkg/checking"
	"github.com/djedjethai/authentication/pkg/http/dto"
	"github.com/djedjethai/toolbox"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

var tools = new(toolbox.Tools)

func Handler(ch checking.Service) http.Handler {
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
	mux.Post("/authenticate", Authenticate(ch))

	return mux
}

func Authenticate(ch checking.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var rp dto.RequestPayload

		err := tools.ReadJSON(w, r, &rp)
		if err != nil {
			tools.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}

		// validate  the user
		user, err := ch.GetByEmail(&rp)
		if err != nil {
			tools.ErrorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
			return
		}

		payload := dto.JsonResponse{
			Error:   false,
			Message: fmt.Sprintf("Logged in user %s", user.Email),
			Data:    user,
		}

		tools.WriteJSON(w, http.StatusAccepted, payload)
	}
}
