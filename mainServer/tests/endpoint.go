package tests

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestEndpoint is a helper function for the common steps of controller tests.
// Performs a http request to the router and makes assertions on the http status code and body.
// The body and expectBody parameters are converted to JSON inside this function,
// can be passed directly as go structs or primitives
func TestEndpoint(t *testing.T, r *gin.Engine, method string, url string, body any, expectCode int, expectBody any) {

	// create the request body
	var reader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		assert.NoError(t, err)
		reader = bytes.NewReader(bodyBytes)
	}

	// create the request
	req, err := http.NewRequest(method, url, reader)
	assert.NoError(t, err)

	// perform the request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// check the response code
	b, err := ioutil.ReadAll(w.Body)
	assert.NoError(t, err)
	if w.Code != expectCode {
		t.Error(w.Code, string(b))
	}

	// check the response body
	if expectBody != nil {
		responseBody := string(b)

		// convert the expected body to JSON interface map
		expectedBytes, err := json.Marshal(expectBody)
		assert.NoError(t, err)
		expectedJson := string(expectedBytes)

		// compare the two
		if responseBody != expectedJson {
			t.Errorf("Expected response body=%v but got actual=%v", expectedJson, responseBody)
		}
	}
}
