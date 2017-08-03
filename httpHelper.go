package main

import (
	"fmt"
	"net/http"
)

//StatusMsgTemplate .
//TODO: Doc
func StatusMsgTemplate(msg string, statuscode int) string {
	return fmt.Sprintf(`{"StatusCode":"%v", "Msg":"%v"}`, statuscode, msg)
}

//StatusMsg .
//TODO: Doc
func StatusMsg(w http.ResponseWriter, msg string, statuscode int) {
	w.WriteHeader(statuscode)
	fmt.Fprintf(w, StatusMsgTemplate(msg, statuscode))
}

//StatusStatusInternalServerError .
//TODO: Doc
func StatusStatusInternalServerError(w http.ResponseWriter, msg string) {
	StatusMsg(w, msg, http.StatusInternalServerError)
}

//StatusUnauthorized .
//TODO: Doc
func StatusUnauthorized(w http.ResponseWriter, msg string) {
	StatusMsg(w, msg, http.StatusUnauthorized)
}

//StatusNotFound .
//TODO: Doc
func StatusNotFound(w http.ResponseWriter, msg string) {
	StatusMsg(w, msg, http.StatusNotFound)
}

//StatusBadRequest .
//TODO: Doc
func StatusBadRequest(w http.ResponseWriter, msg string) {
	StatusMsg(w, msg, http.StatusBadRequest)
}

//StatusOK .
//TODO: Doc
func StatusOK(w http.ResponseWriter, msg string) {
	StatusMsg(w, msg, http.StatusOK)
}
