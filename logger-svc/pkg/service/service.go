package service

import (
	"github.com/djedjethai/logger/pkg/storage"
	"github.com/djedjethai/toolbox"
)

type Service interface {
	InsertPayload(pl JSONPayload, tools *toolbox.Tools) (toolbox.JSONResponse, error)
}

type service struct {
	repo storage.LogEntry
}

func NewService(r storage.LogEntry) *service {
	return &service{
		repo: r,
	}
}

func (s *service) InsertPayload(pl JSONPayload, tools *toolbox.Tools) (toolbox.JSONResponse, error) {
	// add pl en db
	var logEntry storage.LogEntry
	logEntry.Name = pl.Name
	logEntry.Data = pl.Data

	var t toolbox.JSONResponse

	err := s.repo.Insert(logEntry)
	if err != nil {
		return t, err
	}

	t.Error = false
	t.Message = "logged"

	return t, nil
}
