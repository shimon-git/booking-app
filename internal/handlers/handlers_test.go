package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	Key   string
	Value string
}

var theTests = []struct {
	name       string
	url        string
	method     string
	params     []postData
	statusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"majors-suite", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"make-reservation-GET", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"reservation-summary", "/reservation-summary", "GET", []postData{}, http.StatusOK},
	{"check-avilability-GET", "/check-avilability", "GET", []postData{}, http.StatusOK},

	{"make-reservation-POST", "/make-reservation", "POST", []postData{
		{Key: "first_name", Value: "Shimon"},
		{Key: "last_name", Value: "Yaniv"},
		{Key: "email", Value: "shimon0584064942@gmail.com"},
		{Key: "phone", Value: "0584064942"},
	}, http.StatusOK},
	{"check-avilability-POST", "/check-avilability", "POST", []postData{
		{Key: "start", Value: "13-06-2023"},
		{Key: "end", Value: "14-06-2023"},
	}, http.StatusOK},
	{"home", "/check-avilability-json", "POST", []postData{
		{Key: "start", Value: "13-06-2023"},
		{Key: "end", Value: "14-06-2023"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, h := range theTests {
		if h.method == "GET" {
			response, err := testServer.Client().Get(testServer.URL + h.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if response.StatusCode != h.statusCode {
				t.Errorf("For %s, expected status code is: %d, but got %d", h.name, h.statusCode, response.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, data := range h.params {
				values.Add(data.Key, data.Value)
			}
			response, err := testServer.Client().PostForm(testServer.URL+h.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if response.StatusCode != h.statusCode {
				t.Errorf("For %s, expected status code is: %d, but got %d", h.name, h.statusCode, response.StatusCode)
			}
		}
	}

}
