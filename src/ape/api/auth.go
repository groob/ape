package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

var secret = []byte("insecure-signing-key")

// verify using jwt token
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		accept := acceptHeader(r)
		token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			respondError(rw, http.StatusUnauthorized, accept,
				fmt.Errorf("Authentication error: %v", err))
			return
		}
		// successful, pass
		next.ServeHTTP(rw, r)
	})
}

func handleBasicAuth() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// accept := acceptHeader(r)
		username, password, ok := r.BasicAuth()
		if !ok {
			log.Println("not ok")
		}
		fmt.Println(username, password)
	}
}
