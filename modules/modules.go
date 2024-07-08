package modules

import "net/http"

// the generic type that all my api controllers (or handlers) will adhere to, returning structs/objects if no errors and errors if any
// this type will be taken by the util function and wrapped to create a http handler that go router wants
// incase response writer is to be used, since provided by go's http handler so why not just take in
type ApiFunc[T any] func(http.ResponseWriter, http.Request) (T, error)
