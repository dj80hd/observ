package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

var app = New(4)

func TestApp(t *testing.T) {

	var testcases = []struct {
		descr  string
		method string
		url    string
		query  string
		code   int
		app    *App
	}{
		{
			"Basic get no params",
			"GET",
			"/req",
			"",
			200,
			nil,
		},
		{
			"2 secs of failed work",
			"GET",
			"/req?worksecs=1&workfail=true",
			"",
			200,
			app,
		},
		{
			"HTTP code 333",
			"GET",
			"/req?httpcode=333",
			"",
			333,
			app,
		},
		{
			"GET metrics",
			"GET",
			"/metrics",
			"",
			200,
			app,
		},
		{
			"Invalid method",
			"PUT",
			"/req?does-not-matter",
			"",
			405,
			app,
		},
		{
			"Not Found",
			"GET",
			"/not-real",
			"",
			404,
			app,
		},
		{
			"Could not enqueue new work",
			"GET",
			"/req?worksecs=1",
			"",
			507,
			&App{}, //default has no queue
		},
	}

	for _, tt := range testcases {
		is := is.New(t)

		r, err := http.NewRequest(tt.method, tt.url, nil)
		is.NoErr(err)

		w := httptest.NewRecorder()
		tt.app.ServeHTTP(w, r)
		is.Equal(w.Code, tt.code)
	}
}
