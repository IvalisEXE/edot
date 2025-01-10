package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"warehousesvc/pkg/common"
	"warehousesvc/pkg/response"
	"warehousesvc/pkg/validator"
)

// ErrorHandler is a middleware function that handles errors in the Echo framework.
// It processes different types of errors and sends an appropriate JSON response to the client.
func ErrorHandler(err error, c echo.Context) {
	// Get the service context from the Echo context
	ctx := c.Get(common.KEY_API_CONTEXT).(common.ServiceContextManager)

	// Set default status code and response
	code := http.StatusInternalServerError
	response := response.Response{
		Message: http.StatusText(code),
	}

	// Handle HTTP errors
	if e, ok := err.(*echo.HTTPError); ok {
		code = e.Code
		response.Message = e.Message.(string)
	}

	// Handle validation errors
	if e, ok := err.(*validator.ValidationError); ok {
		code = http.StatusUnprocessableEntity
		response.Message = e.Error()
		response.Errors = e.Errors
	}

	// Log internal server errors
	if code == http.StatusInternalServerError {
		var (
			message string = "unhandled error"
			skip    int    = 4
		)

		// Check if the error is a runtime error
		if strings.Contains(err.Error(), "runtime error") {
			message = "runtime error"
			skip = 6
		}

		// Log the error with appropriate caller skip
		ctx.Logger().With(
			zap.String("payload", ctx.GetAudiTrail().Payload),
			zap.Int("user_id", ctx.GetUserHeader().ID),
		).WithOptions(zap.AddCallerSkip(skip)).Error(message, zap.Error(err))
	}

	// Send JSON response to the client
	if err := c.JSON(code, response); err != nil {
		ctx.Logger().Error("Failed to write JSON response:", zap.Error(err))
	}
}
