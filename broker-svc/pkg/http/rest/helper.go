package rest

import (
	"encoding/json"
	"errors"
	// "github.com/djedjethai/broker/pkg/dto"
	"github.com/djedjethai/broker/pkg/service"
	"io"
	"net/http"
)

func readJson(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("Body must have only a single json value")
	}

	return nil
}

func writeJson(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return nil
	}

	return nil
}

func errorJson(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload service.JsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return writeJson(w, statusCode, payload)
}
