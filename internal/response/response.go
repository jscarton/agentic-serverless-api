package response

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Meta struct {
	RequestID string `json:"request_id"`
	Timestamp string `json:"timestamp"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Envelope struct {
	Data  any          `json:"data"`
	Error *ErrorDetail `json:"error"`
	Meta  Meta         `json:"meta"`
}

func OK(c *gin.Context, data any) {
	c.JSON(200, Envelope{
		Data:  data,
		Error: nil,
		Meta:  newMeta(c),
	})
}

func Created(c *gin.Context, data any) {
	c.JSON(201, Envelope{
		Data:  data,
		Error: nil,
		Meta:  newMeta(c),
	})
}

func Fail(c *gin.Context, status int, code, message string) {
	c.JSON(status, Envelope{
		Data:  nil,
		Error: &ErrorDetail{Code: code, Message: message},
		Meta:  newMeta(c),
	})
}

func newMeta(c *gin.Context) Meta {
	return Meta{
		RequestID: c.GetString("request_id"),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}
