package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	we "webengine"
)

// This function is used to do setup before executing the test functions
func TestMain(m *testing.M) {
	// Run the other tests
	os.Exit(m.Run())
}

// Helper function to create a router during testing
func getRouter(withTemplates bool) *we.Engine {
	r := we.Default()
	if withTemplates {
	}
	return r
}

func testHTTPResponse(t *testing.T, r *we.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Create the service and process the above request.
	r.DispatchHandler.ServeHTTP(w, req)

	if !f(w) {
		t.Fail()
	}
}




