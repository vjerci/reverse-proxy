package log_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/stretchr/testify/assert"
	"github.com/vjerci/reverse-proxy/internal/log"
)

type LoggerMock struct {
	Data string
}

func (logger *LoggerMock) Print(data ...any) {
	logger.Data = data[0].(string)
}

func TestLogger(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8000", http.NoBody)
	if err != nil {
		t.Fatalf("failed to create req %s", err)
	}
	req.Header.Set("requestHeader", "requestHeader")

	resp := &http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Response-Header": []string{"responseHeader"},
		},
	}

	testCases := []struct {
		testName string
		logger   *LoggerMock
		req      *http.Request
		reqBody  []byte
		resp     *http.Response
		respBody []byte
		recorder *httptest.ResponseRecorder
	}{
		{
			testName: "response",
			logger:   &LoggerMock{},
			req:      req,
			reqBody:  []byte("req body"),
			resp:     resp,
			respBody: []byte("resp body"),
			recorder: httptest.NewRecorder(),
		},
	}

	for _, test := range testCases {
		factory := log.ResponseWriterFactoryInstance{
			Logger: test.logger,
		}
		factory.New(test.req, test.reqBody, test.recorder).Write(test.resp.StatusCode, test.resp.Header, test.respBody)

		respBytes, err := io.ReadAll(test.recorder.Body)
		if err != nil {
			t.Fatalf("failed to read response body %s", err)
		}

		assert.EqualValues(t, respBytes, test.respBody, "expected body to be written to resp equal")

		assert.EqualValues(t, test.resp.Header, test.recorder.Header(), "expected headers to be written equal")

		assert.EqualValues(t, test.resp.StatusCode, test.recorder.Code, "expected status code to be written to resp equal")

		snapshotName := t.Name() + "-" + test.testName

		err = cupaloy.SnapshotWithName(snapshotName, test.logger.Data)
		if err != nil {
			t.Fatalf("for %s snapshot doesn't match", snapshotName)
		}
	}
}
