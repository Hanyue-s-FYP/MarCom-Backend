package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"

	"github.com/Hanyue-s-FYP/Marcom-Backend/modules"
)

// creates all the necessary header and writes obj in JSON to the response
func ResponseJSON[T any](w http.ResponseWriter, obj *T, statusCode int) {
	jsonBytes, err := json.Marshal(*obj)
	if err != nil {
		// don't need special message for internal server error, at least don't need to be returned to the client
		ResponseError(w, HttpError{
			Code:       http.StatusInternalServerError,
			LogMessage: fmt.Sprintf("failed to marshal JSON: %v", err),
		})
		return
	}
	slog.Info(fmt.Sprintf("%d %s", statusCode, string(jsonBytes)))
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(jsonBytes))
}

type HttpError struct {
	// error code
	Code int
	// optional message that is sent back to the requester, if nil, will get the message through golang's net/http's StatusText() func
	Message string
	// message to be logged for more context on what exactly is the error, this message is not sent back to the requester
	LogMessage string
}

func (h HttpError) Error() string {
	return fmt.Sprintf("HttpError %d: %s ", h.Code, h.LogMessage) +
		If(
			h.Message != "",
			fmt.Sprintf("(response message: %s)", h.Message),
			"",
		)
}

// creates all the necessary header and writes e to the response
func ResponseError(w http.ResponseWriter, e HttpError) {
	slog.Error(e.Error())
	ResponseJSON(w, &struct{ Message string }{Message: e.Message}, e.Code)
}

// creates a http.Handler with my own function type (I prefer making handlers or controllers to return concrete object/struct if success and return error if anything goes wrong, complying with go's error handling way)
// go don't have default parameters, and since for my use case only C operations will return a diff status code for success cases (I am not trying to rebuild a whole new full-fledged net framework :))), will use variadic parameters and determine if a diff code should be used
// the bad thing is there's ntg stopping me from providing more :((
func MakeHttpHandler[T any](handler modules.ApiFunc[T], customCode ...int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		finalCode := http.StatusOK
		if len(customCode) > 0 {
			finalCode = customCode[0]
		}
		if obj, err := handler(w, r); err != nil {
			if errors.As(err, &HttpError{}) {
				ResponseError(w, err.(HttpError))
			} else {
				// should actually panic because not possible returned error is not HttpError, it would be programmer's mistake alrd or some not reversible error alrd, but for gracefulness, will just log, and response with internal server error and everything proceeds
				ResponseError(w, HttpError{
					Code:       http.StatusInternalServerError,
					LogMessage: fmt.Sprintf("expecting HttpError, got other error: %v", err),
				})
			}
		} else {
			ResponseJSON(w, obj, finalCode)
		}
	}
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

// returns every elements that are in `ori` but not in `compared`
func NotIn[T any](ori, compared []T, compareFunc func(a, b T) bool) []T {
	var difference []T
	for _, elem := range ori {
		if !slices.ContainsFunc(compared, func(e T) bool {
			return compareFunc(elem, e)
		}) {
			difference = append(difference, elem)
		}
	}

	return difference
}

// Either type is a workaround for union types in go, eg. getMe can return either investor or business
// can get which type is value by using type assertions
type Either[A, B any] struct {
	Val any
}
