package forms

import (
	"net/http"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r, err := http.NewRequest("POST", "/some-path", nil)
	if err != nil {
		t.Error(err)
	}
	form := New(r.PostForm)
	valid := form.Valid()
	if !valid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r, err := http.NewRequest("POST", "/URI", nil)
	if err != nil {
		t.Error(err)
	}
	form := New(r.PostForm)
	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "b")
	postData.Add("c", "c")

	r.PostForm = postData
	form = New(r.PostForm)

	if !form.Valid() {
		t.Error("Shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	data := url.Values{}
	form := New(data)
	if form.Has("name") {
		t.Error("form shows has field when it does not")
	}

	data.Add("name", "shimon")
	form = New(data)

	if !form.Has("name") {
		t.Error("shows form does not have field when it have")
	}
}

func TestForm_MinLength(t *testing.T) {
	data := url.Values{}
	form := New(data)

	if form.MinLength("notExist", 2) {
		t.Error("form shows min length for non-existing field")
	}

	data.Add("name", "s")
	form = New(data)

	if form.MinLength("notExist", 2) {
		t.Error("shows min length valid when data is shorter")
	}

	formError := form.Errors.Get("notExist")
	if formError == "" {
		t.Error("should have an error but did not get one")
	}

	data.Set("first_name", "shimon")
	form = New(data)

	if !form.MinLength("first_name", 1) {
		t.Error("shows min-length is invalid when data min-length is valid")
	}

	formError = form.Errors.Get("name")
	if formError != "" {
		t.Error("should not have an error but got one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	data := url.Values{}

	form := New(data)
	form.IsEmail("notExist")

	if form.Valid() {
		t.Error("form shows a valid email address for non-existent field")
	}

	data.Add("email", "shimon@dummyDomain")
	form = New(data)
	form.IsEmail("email")

	if form.Valid() {
		t.Error("form shows valid email address for invalid email address")
	}

	data.Set("email", "shimon@golang.com")
	form = New(data)
	form.IsEmail("email")

	if !form.Valid() {
		t.Error("form shows invalid email address for a valid email address")
	}

}
