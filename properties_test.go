package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

type propertiesHTTPHandlerArgs struct {
	req *http.Request
	err error
}

type propertiesTestParams struct {
	name      string
	args      propertiesHTTPHandlerArgs
	wantCode  int
	wantBody  string
	wantCType string
}

func propertiesTest(name, method, uri, body, ctype, wantBody, wantCType string, wantCode int) *propertiesTestParams {
	propertiesT := new(propertiesTestParams)
	propertiesT.name = name
	propertiesT.args.req, propertiesT.args.err = http.NewRequest(method, uri, strings.NewReader(body))
	propertiesT.args.req.Header.Set("Content-Type", ctype)
	propertiesT.wantCode = wantCode
	propertiesT.wantBody = wantBody
	propertiesT.wantCType = wantCType
	return propertiesT
}

func Test_propertiesPOSTHandler(t *testing.T) {

	//Tests table
	tests := []*propertiesTestParams{}

	// Test Not Found
	tests = append(tests,
		propertiesTest(
			"Not Found", //name
			"POST",      //method
			"/properties/vaicurinthia", //uri
			"",                                        //body
			"application/json",                        //ctype
			`{"StatusCode":"404", "Msg":"Not Found"}`, //wantBody
			"application/json; charset=UTF-8",         //wantCType
			http.StatusNotFound,
		),
	)

	// Test Zero Bytes
	tests = append(tests,
		propertiesTest(
			"Zero Bytes",           //name
			"POST",                 //method
			"/properties/populate", //uri
			"",                 //body
			"application/json", //ctype
			`{"StatusCode":"400", "Msg":"Noooooo! Body must contain a properties object"}`, //wantBody
			"application/json; charset=UTF-8",                                              //wantCType
			http.StatusBadRequest,                                                          //wantCode
		),
	)

	// Test is a json
	tests = append(tests,
		propertiesTest(
			"Not a json",           //name
			"POST",                 //method
			"/properties/populate", //uri
			"Hi!",              //body
			"application/json", //ctype
			`{"StatusCode":"400", "Msg":"Ohhh Noooooo! Body shoud be a json object"}`, //wantBody
			"application/json; charset=UTF-8",                                         //wantCType
			http.StatusBadRequest,                                                     //wantCode
		),
	)

	// Test Content-Type
	tests = append(tests,
		propertiesTest(
			"Wrong Content-Type",   //name
			"POST",                 //method
			"/properties/populate", //uri
			`{"Key":"Value"}`,      //body
			"plain/text",           //ctype
			`{"StatusCode":"400", "Msg":"Ohhh Nooo, Content-Type should be application/json"}`, //wantBody
			"application/json; charset=UTF-8",                                                  //wantCType
			http.StatusBadRequest,                                                              //wantCode
		),
	)

	// Test Invalid object
	tests = append(tests,
		propertiesTest(
			"Wrong object",                                                  //name
			"POST",                                                          //method
			"/properties/populate",                                          //uri
			`{"Key":"Value"}`,                                               //body
			"application/json",                                              //ctype
			StatusMsgTemplate(errInvalidFmt.Error(), http.StatusBadRequest), //wantBody
			"application/json; charset=UTF-8",                               //wantCType
			http.StatusBadRequest,                                           //wantCode
		),
	)

	// Test TotalProperties absent or equal 0
	tests = append(tests,
		propertiesTest(
			"errPropTotalAbsentOrZero", //name
			"POST",                 //method
			"/properties/populate", //uri
			`{ "properties": [ { "id": 1, "title": "Imóvel código 1, com 3 quartos e 2 banheiros.", "price": 643000, "description": "Laboris quis quis elit commodo eiusmod qui exercitation. In laborum fugiat quis minim occaecat id.", "lat": 1257, "long": 928, "beds": 3, "baths": 2, "squareMeters": 61 }]}`, //body
			"application/json", //ctype
			StatusMsgTemplate(errPropTotalAbsentOrZero.Error(), http.StatusBadRequest), //wantBody
			"application/json; charset=UTF-8",                                          //wantCType
			http.StatusBadRequest,                                                      //wantCode
		),
	)

	// Test Properties absent or null
	tests = append(tests,
		propertiesTest(
			"errPropAbsentErr",                                              //name
			"POST",                                                          //method
			"/properties/populate",                                          //uri
			`{ "totalProperties": 2, "properties": []}`,                     //body
			"application/json",                                              //ctype
			StatusMsgTemplate(errPropAbsent.Error(), http.StatusBadRequest), //wantBody
			"application/json; charset=UTF-8",                               //wantCType
			http.StatusBadRequest,                                           //wantCode
		),
	)

	// Test TotalProperties has value not equal of amount the properties
	tests = append(tests,
		propertiesTest(
			"errPropAmountDiffErr", //name
			"POST",                 //method
			"/properties/populate", //uri
			`{ "totalProperties": 2, "properties": [ { "id": 1, "title": "Imóvel código 1, com 3 quartos e 2 banheiros.", "price": 643000, "description": "Laboris quis quis elit commodo eiusmod qui exercitation. In laborum fugiat quis minim occaecat id.", "lat": 1257, "long": 928, "beds": 3, "baths": 2, "squareMeters": 61 }]}`, //body
			"application/json",                                                  //ctype
			StatusMsgTemplate(errPropAmountDiff.Error(), http.StatusBadRequest), //wantBody
			"application/json; charset=UTF-8",                                   //wantCType
			http.StatusBadRequest,                                               //wantCode
		),
	)

	// Test valid object
	tests = append(tests,
		propertiesTest(
			"Valid object",         //name
			"POST",                 //method
			"/properties/populate", //uri
			`{ "totalProperties": 2, "properties": [ { "id": 1, "title": "Imóvel código 1, com 3 quartos e 2 banheiros.", "price": 643000, "description": "Laboris quis quis elit commodo eiusmod qui exercitation. In laborum fugiat quis minim occaecat id.", "lat": 1257, "long": 928, "beds": 3, "baths": 2, "squareMeters": 61 }, { "id": 2, "title": "Imóvel código 2, com 4 quartos e 3 banheiros.", "price": 949000, "description": "Anim mollit aliqua adipisicing labore magna pariatur aute nulla. Amet veniam ut voluptate aliquip esse officia adipisicing ipsum.", "lat": 679, "long": 680, "beds": 4, "baths": 3, "squareMeters": 94 }]}`, //body
			"application/json",                        //ctype
			`{"StatusCode":200,"Msg":{"imported":2}}`, //wantBody
			"application/json; charset=UTF-8",         //wantCType
			http.StatusOK,                             //wantCode
		),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.err != nil {
				t.Fatal(tt.args.err)
			}

			resp := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/properties/{action}", propertiesPOSTHandler).Methods("POST")
			router.ServeHTTP(resp, tt.args.req)

			status := resp.Code
			body := resp.Body.String()
			ctype := resp.Header().Get("Content-Type")
			if status != tt.wantCode || body != tt.wantBody || ctype != tt.wantCType {
				t.Errorf("propertiesHandler(): \n got status %v want %v\n got Content-Type %v want %v\n got body %v want %v\n",
					status, tt.wantCode, ctype, tt.wantCType, body, tt.wantBody)
			}
		})
	}
}

func Test_newPropertyPOSTHandler(t *testing.T) {
	//Tests table
	tests := []*propertiesTestParams{}

	// Test Zero Bytes
	tests = append(tests,
		propertiesTest(
			"Zero Bytes",       //name
			"POST",             //method
			"/properties",      //uri
			"",                 //body
			"application/json", //ctype
			`{"StatusCode":"400", "Msg":"Noooooo! Body must contain a properties object"}`, //wantBody
			"application/json; charset=UTF-8",                                              //wantCType
			http.StatusBadRequest,                                                          //wantCode
		),
	)

	// Test is a json
	tests = append(tests,
		propertiesTest(
			"Not a json",       //name
			"POST",             //method
			"/properties",      //uri
			"Hi!",              //body
			"application/json", //ctype
			`{"StatusCode":"400", "Msg":"Ohhh Noooooo! Body shoud be a json object"}`, //wantBody
			"application/json; charset=UTF-8",                                         //wantCType
			http.StatusBadRequest,                                                     //wantCode
		),
	)

	// Test Content-Type
	tests = append(tests,
		propertiesTest(
			"Wrong Content-Type", //name
			"POST",               //method
			"/properties",        //uri
			`{"Key":"Value"}`,    //body
			"plain/text",         //ctype
			`{"StatusCode":"400", "Msg":"Ohhh Nooo, Content-Type should be application/json"}`, //wantBody
			"application/json; charset=UTF-8",                                                  //wantCType
			http.StatusBadRequest,                                                              //wantCode
		),
	)

	// Test Invalid object
	tests = append(tests,
		propertiesTest(
			"Wrong object",                                                   //name
			"POST",                                                           //method
			"/properties",                                                    //uri
			`{"Key":"Value"}`,                                                //body
			"application/json",                                               //ctype
			StatusMsgTemplate(errWrongFormat.Error(), http.StatusBadRequest), //wantBody
			"application/json; charset=UTF-8",                                //wantCType
			http.StatusBadRequest,                                            //wantCode
		),
	)

	// Test valid object
	tests = append(tests,
		propertiesTest(
			"Valid object", //name
			"POST",         //method
			"/properties",  //uri
			`{"title": "Imóvel código 1, com 3 quartos e 2 banheiros.", "price": 643000, "description": "Laboris quis quis elit commodo eiusmod qui exercitation. In laborum fugiat quis minim occaecat id.", "x": 125, "y": 928, "beds": 3, "baths": 2, "squareMeters": 61 }`, //body
			"application/json",                            //ctype
			`{"StatusCode":"200", "Msg":"Created ID: 3"}`, //wantBody
			"application/json; charset=UTF-8",             //wantCType
			http.StatusOK,                                 //wantCode
		),
	)

	// Test valid object duplicated
	tests = append(tests,
		propertiesTest(
			"Valid object duplicated", //name
			"POST",        //method
			"/properties", //uri
			`{"title": "Imóvel código 1, com 3 quartos e 2 banheiros.", "price": 643000, "description": "Laboris quis quis elit commodo eiusmod qui exercitation. In laborum fugiat quis minim occaecat id.", "x": 125, "y": 928, "beds": 3, "baths": 2, "squareMeters": 61 }`, //body
			"application/json", //ctype
			`{"StatusCode":"400", "Msg":"Already exists a property in this cordinates"}`, //wantBody
			"application/json; charset=UTF-8",                                            //wantCType
			http.StatusBadRequest,                                                        //wantCode
		),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.err != nil {
				t.Fatal(tt.args.err)
			}

			resp := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/properties", newPropertyPOSTHandler).Methods("POST")
			router.ServeHTTP(resp, tt.args.req)

			status := resp.Code
			body := resp.Body.String()
			ctype := resp.Header().Get("Content-Type")
			if status != tt.wantCode || body != tt.wantBody || ctype != tt.wantCType {
				t.Errorf("propertiesHandler(): \n got status %v want %v\n got Content-Type %v want %v\n got body %v want %v\n",
					status, tt.wantCode, ctype, tt.wantCType, body, tt.wantBody)
			}
		})
	}
}

func Test_propertiesGETHandler(t *testing.T) {
	//Tests table
	tests := []*propertiesTestParams{}

	// Test valid
	tests = append(tests,
		propertiesTest(
			"Valid",         //name
			"GET",           //method
			"/properties/1", //uri
			"",              //body
			"",              //ctype
			`{"id":1,"title":"Imóvel código 1, com 3 quartos e 2 banheiros.","price":643000,"description":"Laboris quis quis elit commodo eiusmod qui exercitation. In laborum fugiat quis minim occaecat id.","x":1257,"y":928,"beds":3,"baths":2,"provinces":null,"squareMeters":61}`, //wantBody
			"application/json; charset=UTF-8", //wantCType
			http.StatusOK,                     //wantCode
		),
	)

	tests = append(tests,
		propertiesTest(
			"Zero Bytes",    //name
			"GET",           //method
			"/properties/1", //uri
			"xx",            //body
			"",              //ctype
			`{"StatusCode":"400", "Msg":"Noooooo! Body don't must contain a data"}`, //wantBody
			"application/json; charset=UTF-8",                                       //wantCType
			http.StatusBadRequest,                                                   //wantCode
		),
	)

	tests = append(tests,
		propertiesTest(
			"Zero Bytes",        //name
			"GET",               //method
			"/properties/10000", //uri
			"",                  //body
			"",                  //ctype
			`{"StatusCode":"404", "Msg":"10000"}`, //wantBody
			"application/json; charset=UTF-8",     //wantCType
			http.StatusNotFound,                   //wantCode
		),
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.err != nil {
				t.Fatal(tt.args.err)
			}

			resp := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/properties/{id:[0-9]+}", propertiesGETHandler).Methods("GET")
			router.ServeHTTP(resp, tt.args.req)

			status := resp.Code
			body := resp.Body.String()
			ctype := resp.Header().Get("Content-Type")
			if status != tt.wantCode || body != tt.wantBody || ctype != tt.wantCType {
				t.Errorf("propertiesHandler(): \n got status %v want %v\n got Content-Type %v want %v\n got body %v want %v\n",
					status, tt.wantCode, ctype, tt.wantCType, body, tt.wantBody)
			}
		})
	}
}
