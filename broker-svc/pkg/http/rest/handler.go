package rest

import (
	"encoding/json"
	"errors"
	"time"

	// "github.com/djedjethai/broker/pkg/dto"
	"net/http"

	"github.com/djedjethai/broker/pkg/logs"
	"github.com/djedjethai/broker/pkg/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	// the route for the grpc
	mux.Post("/log-grpc", logViaGRPC())

	return mux
}

func logViaGRPC() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestPayload service.RequestPayload
		err := readJson(w, r, &requestPayload)
		if err != nil {
			errorJson(w, err)
			return
		}

		// second arg is the creadential of the server to connect to
		//  third arg are the options
		conn, err := grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			errorJson(w, err)
			return
		}
		defer conn.Close()

		// create the client
		c := logs.NewLogServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_, err = c.WriteLog(ctx, &logs.LogRequest{
			LogEntry: &logs.Log{
				Name: requestPayload.Log.Name,
				Data: requestPayload.Log.Data,
			},
		})
		if err != nil {
			errorJson(w, err)
			return
		}

		var jr service.JsonResponse
		jr.Error = false
		jr.Message = "Logged with grpc"

		_ = writeJson(w, http.StatusOK, jr)
	}
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
			// in this case use rpc comminication
			err, status, payload := js.LogItemViaRPC(requestPayload.Log)
			if err != nil {
				errorJson(w, err, status)
			}
			writeJson(w, status, payload)

			// new method to push log to RabbitMQ, which will push to listener
			// err, status, payload := js.LogEventViaRabbit(requestPayload.Log)
			// if err != nil {
			// 	errorJson(w, err, status)
			// }
			// writeJson(w, status, payload)

			// old method to push log directly to DB
			// err, status, payload := js.LogItem(requestPayload.Log)
			// if err != nil {
			// 	errorJson(w, err, status)
			// }

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
