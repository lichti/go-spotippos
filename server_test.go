package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_indexHandler(t *testing.T) {
	type args struct {
		r *http.Request
		e error
	}

	type test struct {
		name string
		args args
		want int
	}

	tests := []*test{}
	mytest := new(test)
	mytest.name = "Index"
	mytest.args.r, mytest.args.e = http.NewRequest("GET", "/", nil)
	mytest.want = http.StatusUnauthorized
	tests = append(tests, mytest)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.e != nil {
				t.Fatal(tt.args.e)
			}
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(indexHandler)
			handler.ServeHTTP(w, tt.args.r)
			status := w.Code
			if status != tt.want {
				t.Errorf("indexHandler(): got %v want %v",
					status, tt.want)
			}
		})
	}
}

func Test_notFoundHandler(t *testing.T) {
	type args struct {
		r *http.Request
		e error
	}

	type test struct {
		name string
		args args
		want int
	}

	tests := []*test{}
	mytest := new(test)
	mytest.name = "teste1"
	mytest.args.r, mytest.args.e = http.NewRequest("GET", "/404", nil)
	mytest.want = http.StatusNotFound
	tests = append(tests, mytest)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.e != nil {
				t.Fatal(tt.args.e)
			}
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(notFoundHandler)
			handler.ServeHTTP(w, tt.args.r)
			status := w.Code
			if status != tt.want {
				t.Errorf("notFoundHandler(): got %v want %v",
					status, tt.want)
			}
		})
	}
}
