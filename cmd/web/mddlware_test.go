package main

import (
	"net/http"
	"testing"
)

func TestNosurf(t *testing.T) {
	var myH myHandler
	h := Nosurf(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("type is not http.Handler, But is %T", v)
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler
	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("type is not http.Handler, But is %T", v)
	}
}
