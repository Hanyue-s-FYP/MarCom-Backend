package modules

import "net/http"

// the generic type that all my api controllers (or handlers) will adhere to, returning structs/objects if no errors and errors if any
// this type will be taken by the util function and wrapped to create a http handler that go router wants
// incase response writer is to be used, since provided by go's http handler so why not just take in
type ApiFunc[T any] func(http.ResponseWriter, *http.Request) (*T, error)

// generic response structure for those responses that only need to give user a message for success operation
type ExecResponse struct{ Message string }

// wrapper for slice returns as ApiFunc for some reason is unable to infer slice return types
type SliceWrapper[T any] struct {
	Data []T
}
