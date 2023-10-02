package internal

import "github.com/gofiber/fiber/v2"

const (
	OK                  = 200
	NoContent           = 204
	BadRequest          = 400
	Unauthorized        = 401
	Forbidden           = 403
	NotFound            = 404
	InternalServerError = 500
)

const (
	InternalOK = 200
)

type HTTPResponder struct{}

type HTTPResponse struct {
	StatusCode         int         `json:"statusCode"`
	InternalStatusCode int         `json:"internalStatusCode"`
	Status             string      `json:"status"`
	Data               interface{} `json:"data,omitempty"`
	Error              interface{} `json:"error,omitempty"`
}

func NewHTTPResponder() *HTTPResponder {
	return &HTTPResponder{}
}

func (HTTPResponder) Success(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(HTTPResponse{
		StatusCode: OK,
		Status:     "OK",
		Data:       data,
	})
}

func (HTTPResponder) NoContent(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusNoContent).JSON(HTTPResponse{
		StatusCode: NoContent,
		Status:     "No Content",
		Error:      err.Error(),
	})
}

func (HTTPResponder) BadRequest(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusBadRequest).JSON(HTTPResponse{
		StatusCode: BadRequest,
		Status:     "Bad Request",
		Error:      err.Error(),
	})
}

func (HTTPResponder) Unauthorized(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(HTTPResponse{
		StatusCode: Unauthorized,
		Status:     "Unauthorized",
		Error:      err.Error(),
	})
}

func (HTTPResponder) Forbidden(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusForbidden).JSON(HTTPResponse{
		StatusCode: Forbidden,
		Status:     "Forbidden",
		Error:      err.Error(),
	})
}

func (HTTPResponder) NotFound(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusNotFound).JSON(HTTPResponse{
		StatusCode: NotFound,
		Status:     "Not Found",
		Error:      err.Error(),
	})
}

func (HTTPResponder) InternalServerError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusInternalServerError).JSON(HTTPResponse{
		StatusCode: InternalServerError,
		Status:     "Internal Server Error",
		Error:      err.Error(),
	})
}
