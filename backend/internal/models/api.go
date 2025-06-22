package models

// BabbleAPIResponse is the generic output response wrapper for all client-facing output.
// Generics for the win!
type BabbleAPIResponse[T any] struct {
	Body T `json:"body"`
}

type ListInput struct {
	Library string `query:"library" example:"traderpythonlib" doc:"Specific library to lookup"`
}

type ListOutput struct {
	Library string `json:"library"`
	Version string `json:"version"`
}
