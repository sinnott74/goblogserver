package routes

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sinnott74/goblogserver/env"
)

// TestGetHelloWorld tests the helloworld Get handler
func TestGetHelloWorld(t *testing.T) {

	req, err := http.NewRequest("GET", "localhost:"+env.Port()+"/api/helloworld", nil)
	assert.NoError(t, err, "Could not create request")
	rec := httptest.NewRecorder()
	getHelloWorld(rec, req)

	res := rec.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	b, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	assert.NoError(t, err, "Could read response body")

	assert.Equal(t, "\"Hello world\"\n", string(b))
}
