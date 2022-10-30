package responseutil

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type Responseutil struct {
	log *zap.Logger
}

func New(log *zap.Logger) *Responseutil {
	return &Responseutil{
		log: log,
	}
}

type Response struct {
	log    *zap.Logger
	writer http.ResponseWriter
}

func (r *Responseutil) CreateResponse(writer http.ResponseWriter) *Response {
	return &Response{
		log:    r.log,
		writer: writer,
	}
}

func (r *Response) JSON(httpCode int, data interface{}) {
	r.writer.Header().Set("Content-Type", "application/json")
	r.writer.WriteHeader(httpCode)

	json.NewEncoder(r.writer).Encode(data)
}

type StructuredApiError struct {
	message interface{}
}

func (r *Response) Error4xx(httpCode int, message interface{}) {
	r.JSON(httpCode, &StructuredApiError{
		message: message,
	})
}

func (r *Response) Error5xx(err error) {
	r.log.Error("api", zap.String("error", fmt.Sprint(err)))
	r.JSON(http.StatusInternalServerError, &StructuredApiError{
		message: "An unknown error occurred during processing the request",
	})
}
