package main

import (
	"bytes"
	"encoding/json"
)

//IsJSON .
// This function return true if obj is a valid json
func IsJSON(obj []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(obj, &js) == nil

}

//IsJSONObjorList .
// This function return a int for each json type
// 0: non valiud json
// 1: root json object
// 2: root list of json objects
// 3: Unknow
//TODO Unknow must return a error and no integer
func IsJSONObjorList(obj []byte) int {
	obj = bytes.TrimSpace(obj)
	ret := 0
	if IsJSON(obj) {
		switch string(obj)[0:1] {
		case "{":
			ret = 1
		case "[":
			ret = 2
		}
	}
	return ret
}
