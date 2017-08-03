package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusMsg(t *testing.T) {
	type args struct {
		w          http.ResponseWriter
		msg        string
		statuscode int
	}
	tests := []struct {
		name string
		args args
	}{
		{"SimpleTest", args{httptest.NewRecorder(), "Simple Test", 200}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StatusMsg(tt.args.w, tt.args.msg, tt.args.statuscode)
		})
	}
}

func TestStatusStatusInternalServerError(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		msg string
	}
	tests := []struct {
		name string
		args args
	}{
		{"SimpleTest", args{httptest.NewRecorder(), "Server Error"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StatusStatusInternalServerError(tt.args.w, tt.args.msg)
		})
	}
}

func TestStatusUnauthorized(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		msg string
	}
	tests := []struct {
		name string
		args args
	}{
		{"SimpleTest", args{httptest.NewRecorder(), "Unauthorized"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StatusUnauthorized(tt.args.w, tt.args.msg)
		})
	}
}

func TestStatusNotFound(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		msg string
	}
	tests := []struct {
		name string
		args args
	}{
		{"SimpleTest", args{httptest.NewRecorder(), "Not Found"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StatusNotFound(tt.args.w, tt.args.msg)
		})
	}
}

func TestStatusBadRequest(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		msg string
	}
	tests := []struct {
		name string
		args args
	}{
		{"SimpleTest", args{httptest.NewRecorder(), "Bad Request"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StatusBadRequest(tt.args.w, tt.args.msg)
		})
	}
}

func TestStatusOK(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		msg string
	}
	tests := []struct {
		name string
		args args
	}{
		{"SimpleTest", args{httptest.NewRecorder(), "OK"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StatusOK(tt.args.w, tt.args.msg)
		})
	}
}

func TestStatusMsgTemplate(t *testing.T) {
	type args struct {
		msg        string
		statuscode int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"SimpleTest", args{"MyTest", 200}, `{"StatusCode":"200", "Msg":"MyTest"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StatusMsgTemplate(tt.args.msg, tt.args.statuscode); got != tt.want {
				t.Errorf("StatusMsgTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}
