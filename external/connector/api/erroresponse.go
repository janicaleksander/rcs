package api

import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrorResponse defines a standard JSON error payload
// that your API returns when something goes wrong.
type ErrorResponse struct {
	HTTPStatusCode int    `json:"statusCode"` // HTTP status code
	Title          string `json:"title"`      // Short, human-readable title
	Message        string `json:"message"`    // Detailed description of the error
}

// Render implements the render.Renderer interface.
// It sets the appropriate HTTP status code before sending the response.
func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// NewErrorResponse is a helper to quickly build an ErrorResponse.
func NewErrorResponse(status int, title, message string) render.Renderer {
	return &ErrorResponse{
		HTTPStatusCode: status,
		Title:          title,
		Message:        message,
	}
}

// Predefined error response helpers:

func ErrInvalidRequestBody(err error) render.Renderer {
	return NewErrorResponse(
		http.StatusBadRequest,
		"Invalid request body",
		err.Error(),
	)
}

func ErrActorMakeRequest(err error) render.Renderer {
	return NewErrorResponse(
		http.StatusInternalServerError,
		"Actor request failed",
		err.Error(),
	)
}

func ErrInvalidCredentials(err error) render.Renderer {
	return NewErrorResponse(
		http.StatusUnauthorized,
		"Invalid credentials",
		err.Error(),
	)
}

func ErrCreateJWT(err error) render.Renderer {
	return NewErrorResponse(
		http.StatusInternalServerError,
		"JWT creation failed",
		err.Error(),
	)
}

func ErrJwtMiddleware(err error) render.Renderer {
	return NewErrorResponse(
		http.StatusUnauthorized,
		"JWT validation failed",
		err.Error(),
	)
}
func ErrUpdateLocation(err error) render.Renderer {
	return NewErrorResponse(
		http.StatusBadRequest,
		"Can't update location",
		err.Error(),
	)
}

func ErrBadQueryParam(err error) render.Renderer {
	return NewErrorResponse(
		http.StatusBadRequest,
		"Bad query param",
		err.Error(),
	)
}
func ErrGetUserTask(err error) render.Renderer {
	return NewErrorResponse(
		http.StatusBadRequest,
		"Can't fetch user task",
		err.Error(),
	)
}
