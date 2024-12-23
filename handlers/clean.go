package handlers

import (
	"context"
	"log"
	"os"
)

type CleanupFileHandler struct {
	TempFile *os.File
	Done     chan bool
}

func NewCleanupFileHandler(tempFile *os.File) *CleanupFileHandler {
	return &CleanupFileHandler{
		TempFile: tempFile,
		Done:     make(chan bool),
	}
}

func (ch *CleanupFileHandler) MonitorContext(ctx context.Context) {
	go func() {
		select {
		case <-ctx.Done():
			log.Println("Client disconnected")
		case <-ch.Done:
			log.Println("Handler completed")
		}
	}()
}

func (ch *CleanupFileHandler) CleanupFile() {
	defer func() {
		ch.Done <- true
	}()

	if ch.TempFile != nil {
		fileName := ch.TempFile.Name()
		err := os.Remove(fileName)
		if err != nil {
			log.Printf("Error deleting temp file, may have already been deleted")
		} else {
			log.Printf("Temp file has been deleted")
		}
	}
}
