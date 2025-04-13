package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// BabelAPIResponse is the generic output response wrapper for all client-facing output.
// Generics for the win!
type BabelAPIResponse[T any] struct {
	Body T `json:"body"`
}

type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

func GreetingOperation() huma.Operation {
	return huma.Operation{
		OperationID: "get-greeting",
		Method:      http.MethodGet,
		Path:        "/greeting/{name}",
		Summary:     "Get a greeting",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func Greeting(ctx context.Context, input *struct {
	Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
}) (*GreetingOutput, error) {
	resp := &GreetingOutput{}
	resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Name)
	return resp, nil
}

type ListItemsInput struct {
	Library string `query:"library"`
}

type ListItemsOutput struct {
	Library string `json:"library"`
	Version string `json:"version"`
}

func ListOperation() huma.Operation {
	return huma.Operation{
		OperationID: "get-list",
		Method:      http.MethodGet,
		Path:        "/list/",
		Summary:     "List all library documents",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func APIList(ctx context.Context, input *ListItemsInput) (*BabelAPIResponse[ListItemsOutput], error) {
	data := ListItemsOutput{
		Library: "sheesh2",
		Version: "1.3.0",
	}
	resp := &BabelAPIResponse[ListItemsOutput]{Body: data}
	return resp, nil
}

// listDocs generates a list from docs given a hash
func listDocs(hash_key string) {
	fmt.Println("Surprise!")
}

// updateDocs should update an existing doc or add to it if it's not there
func updateDocs(hash_key string, update_item io.Reader) {
	fmt.Println("Not surprise!")
}
