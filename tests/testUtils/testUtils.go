package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
)

type Request struct {
	URL              string
	Method           string
	Payload          map[string]interface{}
	ExpectedCode     int
	ExpectedResponse map[string]interface{}
	Headers          []map[string]string
}

func (request *Request) ExecuteRequest(t *testing.T, Router *mux.Router) {

	payload, err := json.Marshal(request.Payload)
	if err != nil {
		t.Error("Request error: ", err)
	}

	httpRequest, err := http.NewRequest(request.Method, request.URL, bytes.NewReader(payload))
	if err != nil {
		t.Error("Request error: ", err)
	}
	if len(request.Headers) > 0 {
		for _, header := range request.Headers {
			httpRequest.Header.Add(header["key"], header["value"])
		}
	}
	response := httptest.NewRecorder()
	Router.ServeHTTP(response, httpRequest)

	if request.ExpectedCode != response.Code {
		t.Errorf("Expected response code %d. Got %d\n", request.ExpectedCode, response.Code)
	}

	var respBodyM map[string]interface{}
	err = json.Unmarshal(response.Body.Bytes(), &respBodyM)
	if err != nil {
		t.Error("Request error: ", err)
	}

	if eq := reflect.DeepEqual(respBodyM, request.ExpectedResponse); !eq {
		t.Errorf("Expected response to be '%v'. Got '%v'", request.ExpectedResponse, respBodyM)
	}

}
