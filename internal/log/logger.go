package log

import (
	"net/http"
)

type Logger interface {
	Print(data ...any)
}

type ResponseWriterFactory interface {
	New(req *http.Request, reqBody []byte, writer http.ResponseWriter) ResponseWriter
}

type ResponseWriter interface {
	Write(statusCode int, headers map[string][]string, content []byte)
}

type ResponseWriterFactoryInstance struct {
	Logger Logger
}

func (factory *ResponseWriterFactoryInstance) New(req *http.Request, reqBody []byte, writer http.ResponseWriter) ResponseWriter {
	return newLoggingResponseWriter(factory.Logger, req, reqBody, writer)
}

type ResponseWriterInstance struct {
	logger  Logger
	req     *http.Request
	reqBody []byte
	writer  http.ResponseWriter
}

func newLoggingResponseWriter(logger Logger, req *http.Request, reqBody []byte, writer http.ResponseWriter) ResponseWriter {
	return &ResponseWriterInstance{
		logger:  logger,
		req:     req,
		reqBody: reqBody,
		writer:  writer,
	}
}

func (loggingWriter *ResponseWriterInstance) Write(statusCode int, headers map[string][]string, content []byte) {
	err := PrintResponse(loggingWriter.logger, loggingWriter.req, loggingWriter.reqBody, statusCode, headers, content)
	if err != nil {
		panic(err)
	}

	for header, headerValues := range headers {
		for _, headerValue := range headerValues {
			loggingWriter.writer.Header().Set(header, headerValue)
		}
	}

	loggingWriter.writer.WriteHeader(statusCode)
	loggingWriter.writer.Write(content)
}
