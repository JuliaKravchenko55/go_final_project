package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/JuliaKravchenko55/go_final_project/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

const tokenExpirationHours = 8

var jwtKey = []byte(config.GetSecret())

type Credentials struct {
	Password string `json:"password"`
}

type Claims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

func Signin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, `{"error":"Неверный формат запроса"}`, http.StatusBadRequest)
		return
	}

	expectedPassword := os.Getenv("TODO_PASSWORD")
	if expectedPassword == "" || creds.Password != expectedPassword {
		http.Error(w, `{"error":"Неверный пароль"}`, http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(tokenExpirationHours * time.Hour)
	claims := &Claims{
		PasswordHash: creds.Password,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Printf("Error creating token: %v", err)
		http.Error(w, `{"error":"Не удалось создать токен"}`, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
