package main

import (
	"encoding/json"
	"net/http"
)

type FailedJobPayload struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	LastError string `json:"error"`
}

func (app *application) getFailedJobs(w http.ResponseWriter, r *http.Request) {
	jobs,err := app.redis.Client.LRange(r.Context(),"failed_jobs",0,-1).Result()
	if err != nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
	}
	var parsed []FailedJobPayload
	for _,j:=range jobs{
		var job FailedJobPayload
        if err := json.Unmarshal([]byte(j), &job); err == nil {
            parsed = append(parsed, job)
        }
	}
	json.NewEncoder(w).Encode(parsed)
	w.WriteHeader(http.StatusOK)
}