package main

import "testing"

func TestIsJSON(t *testing.T) {
	type args struct {
		obj []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Valid Root Json", args{[]byte(`{"key":"value","number":123}`)}, true},
		{"Valid Array of Json", args{[]byte(`[{"key":"value","number":123}]`)}, true},
		{"Invalid 1", args{[]byte(`[{"key":"value","number":abc}]`)}, false},
		{"Ivalid 2", args{[]byte(`"key":"value"`)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsJSON(tt.args.obj); got != tt.want {
				t.Errorf("IsJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsJsonObjorList(t *testing.T) {
	type args struct {
		obj []byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Valid Root Json", args{[]byte(`{"key":"value","number":123}`)}, 1},
		{"Valid Array of Json", args{[]byte(`[{"key":"value","number":123}]`)}, 2},
		{"Invalid 1", args{[]byte(`[{"key":"value","number":abc}]`)}, 0},
		{"Ivalid 2", args{[]byte(`"key":"value"`)}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsJSONObjorList(tt.args.obj); got != tt.want {
				t.Errorf("IsJsonObjorList() = %v, want %v", got, tt.want)
			}
		})
	}
}
