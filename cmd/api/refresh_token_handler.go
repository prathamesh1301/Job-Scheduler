package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func (app *application) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if payload.RefreshToken == "" {
		http.Error(w, "Refresh token is required", http.StatusBadRequest)
		return
	}

	hashedToken := app.jwt.HashToken(payload.RefreshToken)
	refreshToken, err := app.store.RefreshToken.GetRefreshTokenByToken(r.Context(), hashedToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}
	if time.Now().After(refreshToken.ExpiresAt) {
		http.Error(w, "Refresh token expired", http.StatusUnauthorized)
		return
	}

	newAccessToken, err := app.jwt.GenerateToken(refreshToken.UserID)
	if err != nil {
		http.Error(w, "Failed to generate new access token", http.StatusInternalServerError)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": payload.RefreshToken,
	})
}