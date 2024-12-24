/* The handlers package sets up all the handlers required by the webserver.
 */
package utils

import (
	"context"
	"log/slog"
	"os"
)

type CleanupFileHandler struct {
	TempFile  *os.File
	RequestId string
	Done      chan bool
}

func NewCleanupFileHandler(tempFile *os.File, requestId string) *CleanupFileHandler {
	return &CleanupFileHandler{
		TempFile:  tempFile,
		RequestId: requestId,
		Done:      make(chan bool),
	}
}

func (ch *CleanupFileHandler) MonitorContext(ctx context.Context) {
	go func() {
		select {
		case <-ctx.Done():
			slog.Info("client disconnected", "correlationId", ch.RequestId)
			ch.CleanupFile()
		case <-ch.Done:
			slog.Info("handler completed", "correlationId", ch.RequestId)
		}
	}()
}

func (ch *CleanupFileHandler) CleanupFile() {
	slog.Info("File cleanup triggered", "correlationId", ch.RequestId)
	defer func() {
		ch.Done <- true
	}()

	if ch.TempFile != nil {
		fileName := ch.TempFile.Name()
		err := os.Remove(fileName)
		if err != nil {
			slog.Error("Temp file deletion error", "correlationId", ch.RequestId,
				"filename", fileName, "error", err)
		} else {
			slog.Info("Temp file deleted", "correlationId", ch.RequestId,
				"filename", fileName)
		}
	}
}
