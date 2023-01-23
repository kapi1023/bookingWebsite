package forms

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestValid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	isValid := form.Valid()
	if !isValid {
		t.Error("invalid form")
	}
}

func TestRequired(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	form.Required("a", "b", "c", "d")
	if form.Valid() {
		t.Error("form should be valid when required missing")
	}
	postedData := url.Values{}
	postedData.Add("A", "B")
	postedData.Add("B", "B")
	postedData.Add("C", "B")
	postedData.Add("D", "B")
	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("A", "B", "C")

	if !form.Valid() {
		t.Error("shows does not have required fields")
	}

}

func TestHas(t *testing.T) {

	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	form.Required("D")
	if form.Valid() {
		t.Error("form should be valid when required missing")
	}

	has := form.Has("whatever")
	if has {
		log.Println(has, "elo")
		t.Error("Empty required fields")
	}

	postedData := url.Values{}
	postedData.Add("D", "DDDDD")
	r, _ = http.NewRequest("POST", "/whatever", nil)

	form = New(postedData)
	has = form.Has("D")
	if !has {
		log.Println(has, form)
		t.Error("Empty required fields")
	}

}

func TestMinLegth(t *testing.T) {
	length := 3
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	form.MinLength("somethink", length)

	if form.Valid() {
		log.Println(form)
		t.Error("min lenght error")
	}

	isError := form.Error.Get("D")
	if isError != "" {
		t.Error("should have an error but did not get one 1")
	}

	postedValues := url.Values{}
	postedValues.Add("D", "DDDDD")

	form = New(postedValues)
	form.MinLength("D", 100)
	if form.Valid() {
		t.Error("show 100 minlength")
	}

	postedValues = url.Values{}
	postedValues.Add("C", "ABC")
	form = New(postedValues)
	form.MinLength("C", 1)

	if !form.Valid() {
		t.Error("shows minlength 1 is not when it is")
	}

	isError = form.Error.Get("C")
	if isError != "" {
		t.Error("Should not have error but got one 2")
	}

}

func TestIsEmail(t *testing.T) {
	postedValues := url.Values{}
	form := New(postedValues)
	form.IsEmail("x")
	if form.Valid() {
		log.Println(form)
		t.Error("Form shows valid email for not existing email")
	}

	form = New(postedValues)
	postedValues = url.Values{}
	postedValues.Add("D", "xxxxx@op.pl")
	form = New(postedValues)
	form.IsEmail("D")
	if !form.Valid() {
		log.Println(form)
		t.Error("Form shows not valid for good email")
	}

	form = New(postedValues)
	postedValues = url.Values{}
	postedValues.Add("D", "xxxxx@")
	form = New(postedValues)
	form.IsEmail("D")
	if form.Valid() {
		log.Println(form)
		t.Error("Form shows valid for bad email")
	}

}
