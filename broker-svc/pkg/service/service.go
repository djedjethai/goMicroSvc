package service

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/djedjethai/broker/pkg/domain/jsdomain"

	//"github.com/djedjethai/broker/pkg/dto"
	//"github.com/djedjethai/broker/pkg/http/rest"
	"net/http"
)

type Service interface {
	GetResponse() []JsonResponse
	AddData(JsonResponse) JsonSavedResponse
	Authenticate(a AuthPayload) (error, int, *JsonResponse)
	LogItem(rp LogPayload) (error, int, *JsonResponse)
	SendMail(rp MailPayload) (error, int, *JsonResponse)
}

type service struct {
	r jsdomain.RepoInterf
}

func NewService(r jsdomain.RepoInterf) Service {
	return &service{r}
}

func (s *service) SendMail(rp MailPayload) (error, int, *JsonResponse) {
	var payload JsonResponse

	jsonData, _ := json.MarshalIndent(rp, "", "\t")

	mailServiceURL := "http://mailer-service/send"

	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err, http.StatusNotFound, &payload
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err, http.StatusNotFound, &payload
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return errors.New("Error connecting with the mailer"), response.StatusCode, &payload
	}

	payload.Error = false
	payload.Message = "message sent"

	return nil, http.StatusAccepted, &payload
}

func (s *service) LogItem(rp LogPayload) (error, int, *JsonResponse) {

	var payload JsonResponse

	jsonData, _ := json.MarshalIndent(rp, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err, http.StatusNotFound, &payload
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err, http.StatusNotFound, &payload
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return errors.New("Error connecting with the logger"), response.StatusCode, &payload
	}

	payload.Error = false
	payload.Message = "logged"

	return nil, http.StatusAccepted, &payload

}

func (s *service) Authenticate(a AuthPayload) (error, int, *JsonResponse) {

	var payload *JsonResponse

	// create some json we will send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		return err, http.StatusNotFound, nil
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err, http.StatusInternalServerError, nil
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		return errors.New("invalid credentials"), http.StatusUnauthorized, nil
	} else if response.StatusCode != http.StatusAccepted {
		return errors.New("error calling auth service"), http.StatusUnauthorized, nil
	}

	var jsonFromService JsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		return err, http.StatusInternalServerError, nil
	}

	if jsonFromService.Error {
		return err, http.StatusUnauthorized, nil
	}

	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	return nil, http.StatusAccepted, payload
}

// return the JsonResp in mem to the client
func (s *service) GetResponse() []JsonResponse {
	resp := s.r.RepoResponse()

	var FinalResp []JsonResponse
	for _, r := range resp {
		fr := JsonResponse{
			Error:   r.Error,
			Message: r.Message,
			Data:    r.Data,
		}
		FinalResp = append(FinalResp, fr)

	}

	return FinalResp
}

// add some datas to the memory
func (s *service) AddData(dt JsonResponse) JsonSavedResponse {
	js := jsdomain.JsonResponse{
		Error:   dt.Error,
		Message: dt.Message,
		Data:    dt.Data,
	}
	bl, _ := s.r.RepoAddData(js)

	return JsonSavedResponse{OK: bl}
}
