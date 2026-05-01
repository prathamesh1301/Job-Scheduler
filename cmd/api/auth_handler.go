package main

import (
	"auth/internals/store"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	var payload LoginPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	user, err := app.store.User.GetUserByUsername(r.Context(), payload.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if user != nil {
		err = bcrypt.CompareHashAndPassword(
			user.Password,
			[]byte(payload.Password),
		)
		if err != nil {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}
		jwtToken, err := app.jwt.GenerateToken(user.Username)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}
		refreshToken, err := app.jwt.GenerateRefreshToken()
		if err != nil {
			http.Error(w, "Error generating refresh token", http.StatusInternalServerError)
			return
		}
		hashedToken := app.jwt.HashToken(refreshToken)
		if err = app.store.RefreshToken.InsertRefreshToken(r.Context(), &store.RefreshToken{
			UserID:    user.Username,
			Token:     hashedToken,
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
		}); err != nil {
			http.Error(w, "Error inserting refresh token", http.StatusInternalServerError)
			return
		}
		
		if err = app.store.RefreshToken.EnforceSessionLimit(r.Context(), user.Username, 5); err != nil {
			log.Println("Failed to enforce session limit:", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"access_token": jwtToken,
			"refresh_token": refreshToken,
		})
	}
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	var payload LoginPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(payload.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	_, err = app.store.User.CreateUser(
		r.Context(),
		&store.User{
			Username: payload.Username,
			Password: hashedPassword,
		},
	)
	refreshToken, err := app.jwt.GenerateRefreshToken()
	if err != nil {
		http.Error(w, "Error generating refresh token", http.StatusInternalServerError)
		return
	}
	hashedToken := app.jwt.HashToken(refreshToken)
	if err = app.store.RefreshToken.InsertRefreshToken(r.Context(), &store.RefreshToken{
		UserID:    payload.Username,
		Token:     hashedToken,
		ExpiresAt: time.Now().Add(2 * time.Minute),
	}); err != nil {
		http.Error(w, "Error inserting refresh token", http.StatusInternalServerError)
		return
	}
	
	if err = app.store.RefreshToken.EnforceSessionLimit(r.Context(), payload.Username, 5); err != nil {
		log.Println("Failed to enforce session limit:", err)
	}

	jwtToken, err := app.jwt.GenerateToken(payload.Username)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": jwtToken,
		"refresh_token": refreshToken,
	})
}

func (app *application) ValidateTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "No token provided", http.StatusUnauthorized)
			return
		}

		// Extract token (Bearer <token>)
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		valid, err := app.jwt.Validate(tokenString)
		if err != nil || !valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}