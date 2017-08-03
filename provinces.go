package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

//Point .
type Point struct {
	X int
	Y int
}

//Boundary ..
type Boundary struct {
	UpperLeft   Point
	BottomRight Point
}

//Province .
type Province struct {
	Boundaries Boundary
}

//Provinces .
type Provinces map[string]Province

//Global store of Provinces
var (
	provincesStore Provinces
)

//Populate provincesStore
func provinceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	action := vars["action"]
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	switch action {
	case "populate":
		body, _ := ioutil.ReadAll(r.Body)
		if len(body) <= 0 {
			StatusBadRequest(w, "Noooooo! Body must contain a province object")
			return
		}
		if !IsJSON(body) {
			StatusBadRequest(w, "Ohhh Noooooo! Body shoud be a json object")
			return
		}
		if r.Header.Get("Content-Type") != "application/json" {
			StatusBadRequest(w, "Ohhh Nooo, Content-Type should be application/json")
			return
		}

		var p Provinces
		err := json.Unmarshal(body, &p)
		if err == nil {
			provincesStore = p
			StatusOK(w, "Success")
		} else {
			StatusBadRequest(w, "Wrong format")
		}

	default:
		notFoundHandler(w, r)
	}
}

// Find provinces by Lat and Long of properties
func (p Provinces) findByLatLong(prop *Property) (provName []string) {
	//Targets
	tX := prop.Lat
	tY := prop.Long
	//TODO: add Boundary validate
	for name, prov := range p {
		aX := prov.Boundaries.UpperLeft.X
		aY := prov.Boundaries.UpperLeft.Y
		bX := prov.Boundaries.BottomRight.X
		bY := prov.Boundaries.BottomRight.Y
		if (tX >= aX && tX <= bX) && (tY <= aY && tY >= bY) {
			provName = append(provName, name)
		}
	}
	//TODO: add NotFound Validation
	return provName
}
