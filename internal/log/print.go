package log

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var ErrJSONMarshalReq = errors.New("failed to json marshal req")
var ErrJSONMarshalResp = errors.New("failed to json marshal response")

type RequestLog struct {
	ReqUrl     string              `json:"req_url"`
	ReqMethod  string              `json:"req_method"`
	ReqHeaders map[string][]string `json:"req_headers"`
	ReqBody    string              `json:"req_body"`
}

type ResponseLog struct {
	ResponseHeaders    map[string][]string `json:"resp_headers"`
	ResponseBody       string              `json:"resp_body"`
	ResponseStatusCode int                 `json:"resp_status_code"`
}

func PrintResponse(logger Logger, req *http.Request, reqBody []byte, respStatusCode int, respHeaders map[string][]string, respBody []byte) error {
	buffRequest := bytes.NewBuffer(nil)

	err := json.NewEncoder(buffRequest).Encode(&RequestLog{
		ReqHeaders: req.Header,
		ReqMethod:  req.Method,
		ReqBody:    string(reqBody),
		ReqUrl:     req.RequestURI,
	})

	if err != nil {
		return fmt.Errorf("%s : %s", ErrJSONMarshalReq, err)
	}

	logLine := "got request: \n" + buffRequest.String()

	buffResponse := bytes.NewBuffer(nil)

	err = json.NewEncoder(buffResponse).Encode(&ResponseLog{
		ResponseStatusCode: respStatusCode,
		ResponseHeaders:    respHeaders,
		ResponseBody:       string(respBody),
	})

	if err != nil {
		return fmt.Errorf("%s : %s", ErrJSONMarshalResp, err)
	}

	logLine += "->responding with " + buffResponse.String() + "---"
	logger.Print(logLine)

	return nil
}
