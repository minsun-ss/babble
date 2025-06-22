package handlers

import (
	"babble/backend/internal/models"
	"context"
)

// retrieveLibraries retrieves all items available under a specific key
func retrieveLibraries(team_filter string, project_filter string) {

}

// APIListHandler is the api handler for retrieving all libraries and versions from the database
func APIListHandler(ctx context.Context, input *models.ListInput) (*models.BabbleAPIResponse[models.ListOutput], error) {
	// check to see if there is anything in input
	data := models.ListOutput{
		Library: "sheesh2",
		Version: "1.3.0",
	}
	resp := &models.BabbleAPIResponse[models.ListOutput]{Body: data}
	return resp, nil
}
