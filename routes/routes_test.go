package routes

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sinnott74/goblogserver/database"
	"github.com/stretchr/testify/assert"
)

// Test404Route checks that the server will return 404 when requesting a route that doesn't exist
func Test404Route(t *testing.T) {
	db, err := database.Init()
	assert.NoError(t, err, "Could not connect to  database")
	defer db.Close()

	ts := httptest.NewServer(Handler())
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s/not/a/route", ts.URL))
	assert.NoError(t, err, "Unexpected error requesting 404 route")
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

// TestGetHelloWorldRoute tests the route /api/helloworld
func TestGetHelloWorldRoute(t *testing.T) {

	ts := httptest.NewServer(Handler())
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s/helloworld", ts.URL))
	assert.NoError(t, err, "Could not GET /helloworld")
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	b, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err, "Could read response")

	assert.Equal(t, "\"Hello world\"\n", string(b))
}

// TestGetHelloWorldRoute tests the route /api/helloworld
func TestGetAllBlogpostdRoute(t *testing.T) {

	db, err := database.Init()
	assert.NoError(t, err, "Could not connect to  database")
	defer db.Close()

	ts := httptest.NewServer(Handler())
	defer ts.Close()

	res, err := http.Get(fmt.Sprintf("%s/api/blogposts", ts.URL))
	assert.NoError(t, err, "Could not GET /blogposts")
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	b, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err, "Could read response")

	assert.True(t, len(string(b)) > 100)
}
