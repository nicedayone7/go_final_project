package auth

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		secret := os.Getenv("TODO_SECRET")
		
		if len(pass) > 0 {
			var jwtStr string
			cookieToken, err := r.Cookie("token")
			if err == nil {
				jwtStr = cookieToken.Value
			}

			token, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil {
				http.Error(w, fmt.Sprintf("Error parsing token: %v", err), http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		}
	})
}