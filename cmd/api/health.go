package main

import "net/http"

func (app *application) checkhealth(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Good health"))
}