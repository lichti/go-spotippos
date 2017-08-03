package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//Property info
type Property struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Price        int    `json:"price"`
	Description  string `json:"description"`
	Lat          int    `json:"lat"`
	Long         int    `json:"long"`
	Beds         int    `json:"beds"`
	Baths        int    `json:"baths"`
	SquareMeters int    `json:"squareMeters"`
}

//Properties json structure
type Properties struct {
	TotalProperties int        `json:"totalProperties"`
	Properties      []Property `json:"properties"`
}

//PropertyOutput .
type PropertyOutput struct {
	ID           int      `json:"id"`
	Title        string   `json:"title"`
	Price        int      `json:"price"`
	Description  string   `json:"description"`
	Lat          int      `json:"x"`
	Long         int      `json:"y"`
	Beds         int      `json:"beds"`
	Baths        int      `json:"baths"`
	Provinces    []string `json:"provinces"`
	SquareMeters int      `json:"squareMeters"`
}

//PropertiesFind json structure
type PropertiesFind struct {
	FoundProperties int              `json:"foundProperties"`
	Properties      []PropertyOutput `json:"properties"`
}

//PropertyInput .
type PropertyInput struct {
	Title        string `json:"title"`
	Price        int    `json:"price"`
	Description  string `json:"description"`
	Lat          int    `json:"x"`
	Long         int    `json:"y"`
	Beds         int    `json:"beds"`
	Baths        int    `json:"baths"`
	SquareMeters int    `json:"squareMeters"`
}

/* const */
var (
	errPropTotalAbsentOrZero = errors.New("TotalProperties absent or equal 0")
	errPropAbsent            = errors.New("Properties absent or null")
	errPropAmountDiff        = errors.New("TotalProperties has value not equal of amount the properties")
)

var (
	propertiesStore Properties
)

func propertiesPOSTHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	action := vars["action"]
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	switch action {
	case "populate":

		body, _ := ioutil.ReadAll(r.Body)
		if len(body) <= 0 {
			StatusBadRequest(w, "Noooooo! Body must contain a properties object")
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

		var p Properties
		err := json.Unmarshal(body, &p)
		switch err {
		case nil:
			if p.TotalProperties == 0 && len(p.Properties) == 0 {
				err = errInvalidFmt
				break
			}
			if p.TotalProperties == 0 {
				err = errPropTotalAbsentOrZero
				break
			}
			if len(p.Properties) == 0 {
				err = errPropAbsent
				break
			}
			if p.TotalProperties != len(p.Properties) {
				err = errPropAmountDiff
				break
			}
		}

		if err != nil {
			StatusBadRequest(w, err.Error())
			return
		}

		type ErrLog struct {
			ID     int    `json:"id,omitempty"`
			ErrMsg string `json:"err_msg,omitempty"`
		}

		procErrors := []ErrLog{}

		for _, prop := range p.Properties {
			err = prop.validate()
			iErr := 0
			if err != nil {
				procErrors = append(procErrors, ErrLog{prop.ID, err.Error()})
				iErr++
			}

			err = propertiesStore.existsLatLong(prop)
			if err != nil {
				procErrors = append(procErrors, ErrLog{prop.ID, err.Error()})
				iErr++
			}

			if propertiesStore.existsID(prop) {
				procErrors = append(procErrors, ErrLog{prop.ID, "existsID"})
				iErr++
			}

			if iErr == 0 {
				propertiesStore.append(prop)
			}

		}

		type ErrLogResponse struct {
			Imported int      `json:"imported,omitempty"`
			Errors   []ErrLog `json:"errors,omitempty"`
		}

		type JSONStatus struct {
			StatusCode int            `json:"StatusCode,omitempty"`
			Msg        ErrLogResponse `json:"Msg,omitempty"`
		}

		SuccessJSON, _ := json.Marshal(JSONStatus{http.StatusOK, ErrLogResponse{propertiesStore.TotalProperties, procErrors}})
		fmt.Fprintf(w, "%v", string(SuccessJSON))

	default:
		notFoundHandler(w, r)
	}
}

func newPropertyPOSTHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	body, _ := ioutil.ReadAll(r.Body)
	if len(body) <= 0 {
		StatusBadRequest(w, "Noooooo! Body must contain a properties object")
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

	var propIn PropertyInput
	err := propIn.unmarshal(body)
	if err != nil {
		StatusBadRequest(w, err.Error())
		return
	}

	prop := propIn.toProperty(propertiesStore.nextID())

	err = prop.validate()
	if err != nil {
		StatusBadRequest(w, err.Error())
		return
	}

	err = propertiesStore.existsLatLong(prop)
	if err != nil {
		StatusBadRequest(w, err.Error())
		return
	}

	propertiesStore.append(prop)
	StatusOK(w, fmt.Sprintf("Created ID: %v", prop.ID))

}

func propertiesGETHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strID := vars["id"]
	id, _ := strconv.Atoi(strID)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	body, _ := ioutil.ReadAll(r.Body)
	if len(body) > 0 {
		StatusBadRequest(w, "Noooooo! Body don't must contain a data")
		//return
	} else {
		if propertiesStore.existsID(Property{ID: id}) {
			prop, _ := propertiesStore.findByID(Property{ID: id})
			//TODO:Catch error
			out, _ := prop.outputJSON()
			fmt.Fprintf(w, "%v", string(out))
		} else {
			StatusNotFound(w, strID)
		}
	}
}

func searchPropertiesGETHandler(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	strAX := vars["ax"][0]
	strAY := vars["ay"][0]
	strBX := vars["bx"][0]
	strBY := vars["by"][0]
	aX, _ := strconv.Atoi(strAX)
	aY, _ := strconv.Atoi(strAY)
	bX, _ := strconv.Atoi(strBX)
	bY, _ := strconv.Atoi(strBY)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	body, _ := ioutil.ReadAll(r.Body)
	if len(body) > 0 {
		StatusBadRequest(w, "Noooooo! Body don't must contain a data")
		//return
	} else {

		//TODO:Catch error
		out, _ := propertiesStore.findByRangeLatLong(aX, aY, bX, bY)
		fmt.Fprintf(w, "%v", string(out))
	}

}

func (p *Properties) findByRangeLatLong(aX, aY, bX, bY int) ([]byte, error) {
	var pOut PropertiesFind
	for _, prop := range p.Properties {
		tX := prop.Lat
		tY := prop.Long
		if (tX >= aX && tX <= bX) && (tY <= aY && tY >= bY) {
			propOut := PropertyOutput{
				ID:           prop.ID,
				Title:        prop.Title,
				Price:        prop.Price,
				Description:  prop.Description,
				Lat:          prop.Lat,
				Long:         prop.Long,
				Beds:         prop.Beds,
				Provinces:    provincesStore.findByLatLong(&prop),
				Baths:        prop.Baths,
				SquareMeters: prop.SquareMeters,
			}
			pOut.append(propOut)
		}
	}
	return json.Marshal(pOut)
}

func (prop *Property) outputJSON() ([]byte, error) {
	out := PropertyOutput{
		ID:           prop.ID,
		Title:        prop.Title,
		Price:        prop.Price,
		Description:  prop.Description,
		Lat:          prop.Lat,
		Long:         prop.Long,
		Beds:         prop.Beds,
		Provinces:    provincesStore.findByLatLong(prop),
		Baths:        prop.Baths,
		SquareMeters: prop.SquareMeters,
	}

	return json.Marshal(out)
}

func (p *Properties) findByID(propID Property) (Property, error) {
	for _, prop := range p.Properties {
		if prop.ID == propID.ID {
			return prop, nil
		}
	}
	return Property{}, errors.New("NotFound")
}

func (p *Properties) append(prop Property) {
	p.Properties = append(p.Properties, prop)
	p.TotalProperties = len(p.Properties)
}

func (p *PropertiesFind) append(prop PropertyOutput) {
	p.Properties = append(p.Properties, prop)
	p.FoundProperties = len(p.Properties)
}

func (p *Properties) existsID(prop Property) bool {
	id := prop.ID

	for _, prop := range p.Properties {
		if prop.ID == id {
			return true
		}
	}

	return false
}

func (p *Properties) existsLatLong(prop Property) error {
	lat := prop.Lat
	long := prop.Long

	//for i := 0; i < p.TotalProperties; i++ {
	for _, prop := range p.Properties {
		if prop.Lat == lat && prop.Long == long {
			return errors.New("Already exists a property in this cordinates")
		}
	}

	return nil
}

func (p *Property) validate() error {

	if p.ID <= 0 {
		return errors.New("Invalid ID")
	}

	if len(p.Title) <= 0 {
		return errors.New("Invalid title")
	}

	if p.Price <= 0 {
		return errors.New("Invalid price")
	}

	if len(p.Description) <= 0 {
		return errors.New("Invalid Description")
	}

	if p.Lat < 0 || p.Lat > 1400 {
		return errors.New("Invalid Lat")
	}

	if p.Long < 0 || p.Long > 1000 {
		return errors.New("Invalid Long")
	}

	if p.Beds < 1 || p.Beds > 5 {
		return errors.New("Invalid Beds")
	}

	if p.Baths < 1 || p.Baths > 4 {
		return errors.New("Invalid Baths")
	}

	if p.SquareMeters < 20 || p.SquareMeters > 240 {
		return errors.New("Invalid SquareMeters")
	}

	return nil

}

func (prop *Property) unmarshal(data []byte) error {
	var f interface{}
	err := json.Unmarshal(data, &f)
	if err != nil {
		return err
	}

	tmp := f.(map[string]interface{})
	keys := []string{"id", "title", "price", "description", "lat", "long", "beds", "baths", "squareMeters"}

	for _, key := range keys {
		_, valid := tmp[key]
		if valid == false {
			return errWrongFormat
		}
	}

	err = json.Unmarshal(data, &prop)
	if err != nil {
		return err
	}
	return nil
}

func (propIn *PropertyInput) unmarshal(data []byte) error {
	var f interface{}
	err := json.Unmarshal(data, &f)
	if err != nil {
		return err
	}

	tmp := f.(map[string]interface{})
	keys := []string{"title", "price", "description", "x", "y", "beds", "baths", "squareMeters"}

	for _, key := range keys {
		_, valid := tmp[key]
		if valid == false {
			return errWrongFormat
		}
	}

	err = json.Unmarshal(data, &propIn)
	if err != nil {
		return err
	}
	return nil
}

func (p *Properties) nextID() int {
	var biggest int
	for _, prop := range p.Properties {
		if prop.ID > biggest {
			biggest = prop.ID
		}
	}
	biggest++
	return biggest
}

func (propIn *PropertyInput) toProperty(id int) (prop Property) {
	out := Property{
		ID:           id,
		Title:        propIn.Title,
		Price:        propIn.Price,
		Description:  propIn.Description,
		Lat:          propIn.Lat,
		Long:         propIn.Long,
		Beds:         propIn.Beds,
		Baths:        propIn.Baths,
		SquareMeters: propIn.SquareMeters,
	}
	return out
}
