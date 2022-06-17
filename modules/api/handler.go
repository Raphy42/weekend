package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Raphy42/weekend/core/errors"
)

// todo needs better granularity
func handleError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	statusCode := http.StatusBadRequest
	// restrict by domain
	if errors.HasAnyFlag(err, errors.DEncoding, errors.DResource, errors.DService, errors.DIO, errors.DTransport) {
		if errors.HasFlag(err, errors.ANotFound) {
			statusCode = http.StatusNotFound
		} else if errors.HasFlag(err, errors.ATimeout) {
			statusCode = http.StatusRequestTimeout
		} else if errors.HasFlag(err, errors.ATooBig) && errors.HasFlag(err, errors.DTransport) {
			statusCode = http.StatusRequestEntityTooLarge
		}
	}
	reason := err.Error()
	kind := errors.Diagnostic(err).String()
	c.AbortWithStatusJSON(statusCode, map[string]interface{}{
		"reason": reason,
		"kind":   kind,
	})
}

func MakeJSONHandler[Rq any, Rs any](handler func(ctx context.Context, request *Rq) (*Rs, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Rq
		handleError(c, c.Bind(&req))

		res, err := handler(c, &req)
		handleError(c, err)

		c.JSON(http.StatusOK, res)
	}
}
