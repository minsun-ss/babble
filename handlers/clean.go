package handlers

import (
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
