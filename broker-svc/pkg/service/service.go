package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/djedjethai/broker/pkg/consumer"
	"github.com/djedjethai/broker/pkg/domain/jsdomain"
	"github.com/djedjethai/broker/pkg/emitter"
	// "github.com/djedjethai/broker/pkg/event"
	"net/http"
	"net/rpc"
)

type Service interface {
	GetResponse() []JsonResponse
	AddData(JsonResponse) JsonSavedResponse
	Authenticate(a AuthPayload) (error, int, *JsonResponse)
	LogItemViaRPC(l LogPayload) (error, int, *JsonResponse)
	// LogEventViaRabbit(l LogPayload) (error, int, *JsonResponse)
	// LogItem(rp LogPayload) (error, int, *JsonResponse)
	SendMail(rp MailPayload) (error, int, *JsonResponse)
}

type service struct {
	r    jsdomain.RepoInterf
	cons consumer.Consumer
	emit emitter.Emitter
}

func NewService(r jsdomain.RepoInterf, cons consumer.Consumer, emit emitter.Emitter) Service {
	return &service{r, cons, emit}
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

// new new method which log event to MongoDB using RPC
func (s *service) LogItemViaRPC(l LogPayload) (error, int, *JsonResponse) {

	var payload = JsonResponse{}

	fmt.Println("in LogItemViaRPC")

	// logger-service is the docker-compose svc name
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		return errors.New("Error connecting to RPC server"), http.StatusInternalServerError, &payload
	}

	// the type we gonna use here HAVE TO BE EXACTLY SAME
	// the one the server is expecting to get
	rpcPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}

	var result string
	// "RPCServer is the type one the server side"
	// the last param &result is the response from the server
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		return errors.New("Error sending data to RPC server"), http.StatusInternalServerError, &payload
	}

	payload.Error = false
	payload.Message = result

	return nil, http.StatusAccepted, &payload
}

// new method which push logs to RabbitMQ, which then push to the listener
// func (s *service) LogEventViaRabbit(l LogPayload) (error, int, *JsonResponse) {
// 	var payload = JsonResponse{}
//
// 	err := s.pushToQueue(l.Name, l.Data)
// 	if err != nil {
// 		return errors.New("Error pushing to RabbitMQ"), http.StatusInternalServerError, &payload
// 	}
//
// 	payload.Error = false
// 	payload.Message = "logged via rabbitMQ"
//
// 	return nil, http.StatusAccepted, &payload
// }
//
// func (s *service) pushToQueue(name, msg string) error {
// 	payload := LogPayload{
// 		Name: name,
// 		Data: msg,
// 	}
//
// 	j, _ := json.MarshalIndent(&payload, "", "\t")
// 	// log.INFO is the severity, so we could break various info per severity
// 	// like WARNING etc
// 	err := s.emit.Push(string(j), "log.INFO")
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }

// old method which was logging logs directly to DB
// func (s *service) LogItem(rp LogPayload) (error, int, *JsonResponse) {
//
// 	var payload JsonResponse
//
// 	jsonData, _ := json.MarshalIndent(rp, "", "\t")
//
// 	logServiceURL := "http://logger-service/log"
//
// 	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return err, http.StatusNotFound, &payload
// 	}
//
// 	request.Header.Set("Content-Type", "application/json")
//
// 	client := &http.Client{}
//
// 	response, err := client.Do(request)
// 	if err != nil {
// 		return err, http.StatusNotFound, &payload
// 	}
// 	defer response.Body.Close()
//
// 	if response.StatusCode != http.StatusAccepted {
// 		return errors.New("Error connecting with the logger"), response.StatusCode, &payload
// 	}
//
// 	payload.Error = false
// 	payload.Message = "logged"
//
// 	return nil, http.StatusAccepted, &payload
//
// }

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
