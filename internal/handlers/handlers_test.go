package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shimon-git/booking-app/internal/models"
)

type postData struct {
	Key   string
	Value string
}

var theTests = []struct {
	name       string
	url        string
	method     string
	statusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", http.StatusOK},
	{"majors-suite", "/majors-suite", "GET", http.StatusOK},
	{"reservation-summary", "/reservation-summary", "GET", http.StatusOK},
	{"check-avilability-GET", "/check-avilability", "GET", http.StatusOK},

	/* {"make-reservation-POST", "/make-reservation", "POST", []postData{
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
	}, http.StatusOK}, */
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, h := range theTests {
		response, err := testServer.Client().Get(testServer.URL + h.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if response.StatusCode != h.statusCode {
			t.Errorf("For %s, expected status code is: %d, but got %d", h.name, h.statusCode, response.StatusCode)
		}

	}

}

func TestReposetory_Reservation(t *testing.T) {
	resrvation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	req := httptest.NewRequest("GET", "/", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", resrvation)
	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted: %d", rr.Code, http.StatusOK)
	}

	// test case when reservation is not in session (reset everything)
	req = httptest.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test with non-existing room
	req = httptest.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	resrvation.RoomID = 100
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestReposetory_PostReservation(t *testing.T) {
	resrvation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	reqBody := "start_date=20-6-2030"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=20-6-2030")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=shimon")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=yaniv")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=shimon@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req := httptest.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", resrvation)
	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation  handler returned wrong response code: got %d, wanted: %d", rr.Code, http.StatusSeeOther)
	}

	// test for missing post body
	req = httptest.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for missing post body: got %d, wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test for invalid start date
	reqBody = "start_date=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=20-6-2030")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=shimon")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=yaniv")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=shimon@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
			
	req = httptest.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalod start date: got %d, wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}


	// test for invalid end date
	reqBody = "start_date=20-6-2030"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=invalid")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=shimon")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=yaniv")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=shimon@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
			
	req = httptest.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalod start date: got %d, wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}


	// test for invalid room id
	reqBody = "start_date=20-6-2030"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=21-6-2030")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=shimon")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=yaniv")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=shimon@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=invalid")
			
	req = httptest.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalod start date: got %d, wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}	
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
