package main

import (
	"encoding/json"
	"net/http"
)

type JobPayload struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (app *application) addJobHandler(w http.ResponseWriter, r *http.Request) {
	var payload JobPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jobData, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return	
	}
	if err := app.redis.EnqueueJob(r.Context(),"job_queue",jobData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Job added to queue",
	})
}