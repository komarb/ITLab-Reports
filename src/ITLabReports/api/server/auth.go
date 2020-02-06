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
}

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {

		aud := "itlab"
		checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, true)
		if !checkAud {
			return token, errors.New("Invalid audience")
		}

		// Verify 'iss' claim
		iss := "https://dev.identity.rtuitlab.ru"
		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, true)
		if !checkIss {
			return token, errors.New("Invalid issuer")
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
	resp, err := http.Get("https://pastebin.com/raw/D7UL1cbH")

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
		cert, err := getPemCert(token)
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

func getTokenAndCheckScope(w http.ResponseWriter, r *http.Request){
	authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
	token := authHeaderParts[1]

	if !checkScope("rtuitlab.reports", token) {
		message := "Insufficient scope."
		w.WriteHeader(401)
		w.Write([]byte(message))
		return
	}
}