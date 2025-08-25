package deviceservice

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrorResponse struct {
	HTTPStatusCode int    `json:"statuscode"` // status code of request
	Title          string `json:"title"`      // user-level error
	Message        string `json:"message"`    // user-level error
}

func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrorResponse{
		HTTPStatusCode: 400,
		Title:          "Invalid request",
		Message:        err.Error(),
	}
}
func ErrMakeRequest(err error) render.Renderer {
	return &ErrorResponse{
		HTTPStatusCode: 501,
		Title:          "Can't make request",
		Message:        err.Error(),
	}
}

func ErrInvalidRespondMessage(err error) render.Renderer {
	return &ErrorResponse{
		HTTPStatusCode: 501,
		Title:          "Invalid respond message",
		Message:        err.Error(),
	}
}

func ErrCreateJWT(err error) render.Renderer {
	return &ErrorResponse{
		HTTPStatusCode: 501,
		Title:          "JWT error",
		Message:        err.Error(),
	}
}
