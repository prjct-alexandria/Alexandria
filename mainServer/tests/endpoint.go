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
	"reflect"
	"testing"
)

// TestEndpoint is a helper function for the common steps of controller tests
// performs a http request to the router and makes assertions on the http status code and body.
// Body and expectBody parameters can be passed as their model or entity structs directly, no need to convert to JSON string.
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
		// convert the expected body to JSON interface map
		var expectedInterface map[string]interface{}
		expected, _ := json.Marshal(expectBody)
		err = json.Unmarshal(expected, &expectedInterface)
		assert.NoError(t, err)

		// convert the response body to JSON interface map
		responseType := reflect.TypeOf(expectBody)
		responseInterface := reflect.Zero(responseType).Interface()
		err = json.Unmarshal(b, &responseInterface)
		assert.NoError(t, err)

		// compare the two
		if !reflect.DeepEqual(responseInterface, expectedInterface) {
			t.Errorf("Expected response body=%v but got actual=%v", responseInterface, expectedInterface)
		}
	}
}
