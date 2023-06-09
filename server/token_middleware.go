package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"node/config"
	"strings"

	"github.com/golang-jwt/jwt"
)

const (
	AUTH_HEADER = "authorization"
)

var (
	ErrExtractClaims = errors.New("could not extract claims")
	ErrJWTParse      = errors.New("could not parse jwt")
	ErrNoAuthHeader  = errors.New("invalid value for header: 'Authorization'")
	ErrNoAuthToken   = errors.New("token missing in 'Authorization' header")
)

func VerifyToken(token, secretToken string) (map[string]interface{}, error) {
	tok, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrJWTParse
		}
		return []byte(secretToken), nil
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if ok {
		return claims, nil
	}

	return nil, ErrExtractClaims
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(AUTH_HEADER)
		if authHeader == "" {
			err := ErrNoAuthHeader
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(ErrorToResponse(err))
			return
		}

		bearer_n_token := strings.Split(authHeader, " ")
		if len(bearer_n_token) < 2 {
			err := ErrNoAuthToken
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(ErrorToResponse(err))
			return
		}
		jwtToken := bearer_n_token[1]
		if jwtToken == "" {
			err := ErrNoAuthToken
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(ErrorToResponse(err))
			return
		}

		secretToken := config.Configuration.SecretToken
		claims, err := VerifyToken(jwtToken, secretToken)
		if err != nil {
			log.Println("error while decoding jwt token:", err)
			log.Println(r.Header.Get("User-Agent"))
			log.Println(jwtToken)
			http.Error(w, err.Error(), http.StatusInsufficientStorage)
			return
		}

		newReq := r.WithContext(
			context.WithValue(r.Context(), "claims", claims),
		)
		next.ServeHTTP(w, newReq)
	})
}
