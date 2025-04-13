package handlers

import (
	"babel/backend/internal/models"
	"context"
	"log/slog"
)

type ListInput struct {
	Library string `query:"library" example:"traderpythonlib" doc:"Specific library to lookup"`
}

type ListOutput struct {
	Library string `json:"library"`
	Version string `json:"version"`
}

// APIListHandler is the api handler for retrieving all libraries and versions from the database
func APIListHandler(ctx context.Context, input *ListInput) (*models.BabelAPIResponse[ListOutput], error) {
	// check to see if there is anything in input
	slog.Error("Data from input", input)
	data := ListOutput{
		Library: "sheesh2",
		Version: "1.3.0",
	}
	resp := &models.BabelAPIResponse[ListOutput]{Body: data}
	return resp, nil
}
