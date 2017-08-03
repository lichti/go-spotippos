package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

type provinceHTTPHandlerArgs struct {
	req *http.Request
	err error
}

type provinceTestParams struct {
	name      string
	args      provinceHTTPHandlerArgs
	wantCode  int
	wantBody  string
	wantCType string
}

func provinceTest(name, method, uri, body, ctype, wantBody, wantCType string, wantCode int) *provinceTestParams {
	provinceT := new(provinceTestParams)
	provinceT.name = name
	provinceT.args.req, provinceT.args.err = http.NewRequest(method, uri, strings.NewReader(body))
	provinceT.args.req.Header.Set("Content-Type", ctype)
	provinceT.wantCode = wantCode
	provinceT.wantBody = wantBody
	provinceT.wantCType = wantCType
	return provinceT
}

func Test_provinceHandler(t *testing.T) {

	//Tests table
	tests := []*provinceTestParams{}

	// Test Not Found
	tests = append(tests,
		provinceTest(
			"Not Found", //name
			"POST",      //method
			"/provinces/vaicurinthia", //uri
			"",                                        //body
			"application/json",                        //ctype
			`{"StatusCode":"404", "Msg":"Not Found"}`, //wantBody
			"application/json; charset=UTF-8",         //wantCType
			http.StatusNotFound,
		),
	)

	// Test Zero Bytes
	tests = append(tests,
		provinceTest(
			"Zero Bytes",          //name
			"POST",                //method
			"/provinces/populate", //uri
			"",                 //body
			"application/json", //ctype
			`{"StatusCode":"400", "Msg":"Noooooo! Body must contain a province object"}`, //wantBody
			"application/json; charset=UTF-8",                                            //wantCType
			http.StatusBadRequest,                                                        //wantCode
		),
	)

	// Test is a json
	tests = append(tests,
		provinceTest(
			"Not a json",          //name
			"POST",                //method
			"/provinces/populate", //uri
			"Hi!",              //body
			"application/json", //ctype
			`{"StatusCode":"400", "Msg":"Ohhh Noooooo! Body shoud be a json object"}`, //wantBody
			"application/json; charset=UTF-8",                                         //wantCType
			http.StatusBadRequest,                                                     //wantCode
		),
	)

	// Test Content-Type
	tests = append(tests,
		provinceTest(
			"Wrong Content-Type",  //name
			"POST",                //method
			"/provinces/populate", //uri
			`{"Key":"Value"}`,     //body
			"plain/text",          //ctype
			`{"StatusCode":"400", "Msg":"Ohhh Nooo, Content-Type should be application/json"}`, //wantBody
			"application/json; charset=UTF-8",                                                  //wantCType
			http.StatusBadRequest,                                                              //wantCode
		),
	)

	// Test wrong object
	tests = append(tests,
		provinceTest(
			"Wrong object",                               //name
			"POST",                                       //method
			"/provinces/populate",                        //uri
			`{"Key":"Value"}`,                            //body
			"application/json",                           //ctype
			`{"StatusCode":"400", "Msg":"Wrong format"}`, //wantBody
			"application/json; charset=UTF-8",            //wantCType
			http.StatusBadRequest,                        //wantCode
		),
	)

	// Test valid object
	tests = append(tests,
		provinceTest(
			"Valid object",        //name
			"POST",                //method
			"/provinces/populate", //uri
			`{"Gode" : {"boundaries" : {"upperLeft" : {"x" : 0,"y" : 1000},"bottomRight" : {"x" : 600,"y" : 500}}}}`, //body
			"application/json",                      //ctype
			`{"StatusCode":"200", "Msg":"Success"}`, //wantBody
			"application/json; charset=UTF-8",       //wantCType
			http.StatusOK,                           //wantCode
		),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.err != nil {
				t.Fatal(tt.args.err)
			}

			resp := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/provinces/{action}", provinceHandler).Methods("POST")
			router.ServeHTTP(resp, tt.args.req)

			status := resp.Code
			body := resp.Body.String()
			ctype := resp.Header().Get("Content-Type")
			if status != tt.wantCode || body != tt.wantBody || ctype != tt.wantCType {
				t.Errorf("provinceHandler(): \n got status %v want %v\n got Content-Type %v want %v\n got body %v want %v\n",
					status, tt.wantCode, ctype, tt.wantCType, body, tt.wantBody)
			}
		})
	}
}

func TestProvinces_findByLatLong(t *testing.T) {
	type args struct {
		prop *Property
	}

	type test struct {
		name         string
		p            Provinces
		args         args
		wantProvName []string
	}

	var tProp Property
	tProp.Lat = 252
	tProp.Long = 868
	var tests []test
	var test1 test
	test1.name = "Fisrt Case"
	test1.p = provincesStore
	test1.args.prop = &tProp
	//test1.wantProvName = []string{"Gode", "Ruja"}
	test1.wantProvName = []string{"Gode"}

	tests = append(tests, test1)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotProvName := tt.p.findByLatLong(tt.args.prop); !reflect.DeepEqual(gotProvName, tt.wantProvName) {
				t.Errorf("Provinces.findByLatLong() = %v, want %v", gotProvName, tt.wantProvName)
			}
		})
	}
}
