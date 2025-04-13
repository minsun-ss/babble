package models

// BabelAPIResponse is the generic output response wrapper for all client-facing output.
// Generics for the win!
type BabelAPIResponse[T any] struct {
	Body T `json:"body"`
}
