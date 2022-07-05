package service

type RPCPayload struct {
	Name string
	Data string
}

type JsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// return after a JsonResponse is saved
type JsonSavedResponse struct {
	OK bool `json:"ok"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

//
//
// // import (
// // 	"encoding/json"
// 	"errors"
// 	// "github.com/djedjethai/broker/pkg/dto"
// 	"io"
// 	"net/http"
// )

// func writeJson(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
// 	out, err := json.Marshal(data)
// 	if err != nil {
// 		return err
// 	}
//
// 	if len(headers) > 0 {
// 		for k, v := range headers[0] {
// 			w.Header()[k] = v
// 		}
// 	}
//
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(status)
// 	_, err = w.Write(out)
// 	if err != nil {
// 		return nil
// 	}
//
// 	return nil
// }
//
// func errorJson(w http.ResponseWriter, err error, status ...int) error {
// 	statusCode := http.StatusBadRequest
//
// 	if len(status) > 0 {
// 		statusCode = status[0]
// 	}
//
// 	var payload dto.JsonResponse
// 	payload.Error = true
// 	payload.Message = err.Error()
//
// 	return writeJson(w, statusCode, payload)
// }
