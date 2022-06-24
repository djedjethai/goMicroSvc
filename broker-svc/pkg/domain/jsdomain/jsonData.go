package jsdomain

type JsonResponse struct {
	Error   bool
	Message string
	Data    interface{}
}
