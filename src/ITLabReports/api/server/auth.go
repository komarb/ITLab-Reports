package server

import (
	"errors"
	"fmt"
	"github.com/auth0-community/go-auth0"
	"gopkg.in/square/go-jose.v2"
	"log"
	"net/http"
)
var validator *auth0.JWTValidator


func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Current token: ", r.Header.Get("Authorization"))
		client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: cfg.Auth.KeyURL}, nil)
		audience := cfg.Auth.Audience
		configuration := auth0.NewConfiguration(client, []string{audience}, cfg.Auth.Issuer, jose.RS256)
		validator = auth0.NewValidator(configuration, nil)

		token, err := validator.ValidateRequest(r)
		if err != nil {
			fmt.Println("Token is not valid:", token)
			fmt.Println("Error:", err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Token is not valid\nError: "))
			w.Write([]byte(err.Error()))
			return
		}

		claims := map[string]interface{}{}
		err = validator.Claims(r, token, &claims)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Invalid claims")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid claims"))
			return
		}

		if !checkScope(cfg.Auth.Scope, claims) {
			fmt.Println(err)
			fmt.Println("Invalid scope")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid scope"))
			return
		}
		log.Println("Current claims: ", claims)
		next.ServeHTTP(w, r)


	})
}

func testAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header.Get("Authorization"))
		secret := []byte("test")
		secretProvider := auth0.NewKeyProvider(secret)
		audience := cfg.Auth.Audience
		configuration := auth0.NewConfigurationTrustProvider(secretProvider, []string{audience}, cfg.Auth.Issuer)
		validator = auth0.NewValidator(configuration, nil)
		token, err := validator.ValidateRequest(r)
		if err != nil {
			fmt.Println("Token is not valid:", token)
			fmt.Println("Error:", err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Token is not valid\nError: "))
			w.Write([]byte(err.Error()))
			return
		}
		claims := map[string]interface{}{}
		err = validator.Claims(r, token, &claims)
		log.Println("Current claims: ", claims)
		next.ServeHTTP(w, r)
	})
}
func checkScope(scopeStr string, claims map[string]interface{}) bool {
	var hasScope = false
	_, okScope := claims[scopeStr].(map[string]interface{})

	if !okScope || okScope {
		hasScope = true
	}
	return hasScope
}

func getClaim(r *http.Request, claim string) (string, error) {
	token, err := validator.ValidateRequest(r)
	if err != nil {
		return "", err
	}
	claims := map[string]interface{}{}
	err = validator.Claims(r, token, &claims)

	switch claim {
	case "sub":
		if _, ok := claims["sub"]; ok {
			return fmt.Sprintf("%v", claims["sub"]), nil
		} else {
			return "", errors.New("there is no Sub claim in token")
		}
	case "role":
		if _, ok := claims["role"]; ok {
			return fmt.Sprintf("%v", claims["role"]), nil
		} else {
			return "", errors.New("there is no Role claim in token")
		}
	default:
		return "", errors.New("requested claim is invalid")
	}
}