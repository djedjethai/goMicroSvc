package storage

import (
	"context"
	"github.com/djedjethai/logger/pkg/logs"
)

type LogServer struct {
	// logs.UnimplementedLogServiceServer
	logs.UnimplementedLogServiceServer
	Models Models // necessary methods to write to mongo from helper.go
}

// func which is refered inside the proto file
func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry() // from our proto file(which has been compiled)

	// write the logs
	logEntry := LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	// write to mongoDB
	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{Result: "Failed"}
		return res, err
	}

	// return response
	res := &logs.LogResponse{Result: "Logged"}

	return res, nil
}
