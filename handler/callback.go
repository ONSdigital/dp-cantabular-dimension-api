package handler

import (
	"errors"
	"net/http"
)

// statusCode is a callback function that allows you to extract
// a status code from an error, or returns 500 as a default
func statusCode(err error) int {
	var cerr coder
	if errors.As(err, &cerr) {
		if code := cerr.Code(); code != 0 {
			return code
		}
	}

	return http.StatusInternalServerError
}
