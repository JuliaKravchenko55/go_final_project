package middleware

import (
	"github.com/JuliaKravchenko55/go_final_project/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := config.GetAppPassword()
		if len(pass) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil {
				if err == http.ErrNoCookie {
					http.Error(w, `{"error":"Аутентификация требуется"}`, http.StatusUnauthorized)
					return
				}
				http.Error(w, `{"error":"Ошибка получения куки"}`, http.StatusBadRequest)
				return
			}

			tokenStr := cookie.Value
			claims := &Claims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})
			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					http.Error(w, `{"error":"Неверный токен"}`, http.StatusUnauthorized)
					return
				}
				http.Error(w, `{"error":"Ошибка валидации токена"}`, http.StatusBadRequest)
				return
			}

			if !token.Valid {
				http.Error(w, `{"error":"Неверный токен"}`, http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
