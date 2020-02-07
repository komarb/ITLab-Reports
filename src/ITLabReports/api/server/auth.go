package server

import (
	"encoding/json"
	"errors"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"strings"
)

type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	N string `json:"n"`
	E string `json:"e"`
	X5c []string `json:"x5c"`
	Alg string	`json:"alg"`
}

type CustomClaims struct {
	jwt.StandardClaims
	Scope []string `json:"scope"`
	Sub 	string	`json:"sub"`
}

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		if(cfg.App.TestMode) {
			return nil, nil
		}
		aud := cfg.Auth.Audience
		checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, true)
		if !checkAud {
			return token, errors.New("Invalid audience")
		}

		// Verify 'iss' claim
		iss := cfg.Auth.Issuer
		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, true)
		if !checkIss {
			return token, errors.New("Invalid issuer")
		}

		cert, err := getPemCert(cfg.Auth.KeyURL,token)
		if err != nil {
			log.Println(err.Error())
		}

		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		if !checkScope(cfg.Auth.Scope, token.Raw) {
			return token, errors.New("Invalid scope")
		}

		return result, nil



	},
	SigningMethod: jwt.SigningMethodRS256,
})

var mySigningKey = []byte("test")
var testJwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

func getPemCert(key string, token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get(key)

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}

func checkScope(scope string, tokenString string) bool {
	token, _ := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func (token *jwt.Token) (interface{}, error) {
		cert, err := getPemCert(cfg.Auth.KeyURL, token)
		if err != nil {
			return nil, err
		}
		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	})

	claims, ok := token.Claims.(*CustomClaims)
	hasScope := false
	if ok && token.Valid {
		for i := range claims.Scope {
			if claims.Scope[i] == scope {
				hasScope = true
			}
		}
	}

	return hasScope
}

func getSubClaim(r *http.Request) string {
	authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
	tokenString := authHeaderParts[1]

	token, _ := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func (token *jwt.Token) (interface{}, error) {
		cert, err := getPemCert(cfg.Auth.KeyURL, token)
		if err != nil {
			return nil, err
		}
		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	})

	claims, _ := token.Claims.(*CustomClaims)
	return claims.Sub
}
