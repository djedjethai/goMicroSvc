package handlehttp

import (
	"fmt"
	"github.com/djedjethai/mail/pkg/dto"
	"github.com/djedjethai/mail/pkg/sender"
	"github.com/djedjethai/toolbox"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

var tools = new(toolbox.Tools)

func Handler(s sender.Sender) http.Handler {
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
	mux.Post("/send", sendMail(s))

	return mux
}

func sendMail(s sender.Sender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestPayload dto.MailMessage

		err := tools.ReadJSON(w, r, &requestPayload)
		if err != nil {
			tools.ErrorJSON(w, err)
			return
		}

		err = s.SendSMTPMessage(requestPayload)
		if err != nil {
			tools.ErrorJSON(w, err)
			return
		}

		fmt.Println("what the fuckkkkk")

		payload := toolbox.JSONResponse{
			Error:   false,
			Message: "Sent to " + requestPayload.To,
		}

		tools.WriteJSON(w, http.StatusAccepted, payload)
	}
}
