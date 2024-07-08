package utils

import (
	"fmt"
	"net/http"

	"github.com/Hanyue-s-FYP/Marcom-Backend/modules"
)

// creates all the necessary header and writes obj in JSON to the response
func ResponseJSON[T any](obj T, w http.ResponseWriter) {

}

type HttpError struct {
	// error code
	Code int
	// optional message that is sent back in JSON to the requester
	Message string
}

func (h HttpError) Error() string {
	return fmt.Sprintf("%d: %s", h.Code, h.Message)
}

// creates all the necessary header and writes e to the response
func ResponseError(e error, w http.ResponseWriter) {

}

// creates a http.Handler with my own function type (I prefer making handlers or controllers to return concrete object/struct if success and return error if anything goes wrong, complying with go's error handling way)
func MakeHttpHandler[T any](handler modules.ApiFunc[T]) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {}
}

// Solely for QoL development (like adding ternary to Go :))
// Ternary util
func If[T any](cond bool, valTrue, valFalse T) T {
	if cond {
		return valTrue
	} else {
		return valFalse
	}
}
