package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/kapi1023/bookingWebsite/internal/models"
)

type postData struct {
	key   string
	value string
}

var theThest = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"generalsQuarters", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"majorsSuite", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	// 	{"reservationSummary", "/reservation-summary", "GET", []postData{}, http.StatusOK},
	// 	{"PostSearchAvailability", "/search-availability", "POST", []postData{{key: "start", value: "2020-01-01"}, {key: "end", value: "2020-01-04"}}, http.StatusOK},
	// 	{"PostSearchAvailabilityJson", "/search-availability-json", "POST", []postData{{key: "start", value: "2020-01-01"}, {key: "end", value: "2020-01-04"}}, http.StatusOK},
	// 	{"PostMakeReservation", "/make-reservation", "POST", []postData{
	// 		{key: "first_name", value: "John"},
	// 		{key: "last_name", value: "Smith"},
	// 		{key: "email", value: "Smith@op.pl"},
	// 		{key: "phone", value: "434232312"},
	// 	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()
	for _, e := range theThest {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}

}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomId: 1,
		Room: models.Room{
			Id:       1,
			RoomName: "Generals Quarters",
		},
	}
	req, err := http.NewRequest("GET", "/make-reservation", nil)
	if err != nil {
		log.Println(err)
	}
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Error("Reservation handler return wrong response code: ", rr.Code)
	}
	//test case where reservation is not in session (reset everything)

	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Error("Reservation handler return wrong response code: ", rr.Code)
	}

	//test with non-existing room

	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomId = 100
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Error("Reservation handler return wrong response code: ", rr.Code)
	}

}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
