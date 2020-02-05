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
	Scope string `json:"scope"`
	jwt.StandardClaims
}

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {

		aud := "itlab"
		checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
		if !checkAud {
			return token, errors.New("Invalid audience")
		}

		// Verify 'iss' claim
		iss := "https://dev.identity.rtuitlab.ru"
		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
		if !checkIss {
			return token, errors.New("Invalid issuer")
		}

		claims, ok  := token.Claims.(*CustomClaims)

		hasScope := false

		if ok && token.Valid {
			result := strings.Split(claims.Scope, " ")
			for i := range result {
				if result[i] == "rtuitlab.reports" {
					hasScope = true
				}
			}
			if !hasScope {
				return token, errors.New("Scope does not exist")
			}
		}


		cert, err := getPemCert(token)
		if err != nil {
			log.Println(err.Error())
		}


		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))

		return result, nil



	},
	SigningMethod: jwt.SigningMethodRS256,
})

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://pastebin.com/raw/872PpC6G")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)
	log.Println(jwks)
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
