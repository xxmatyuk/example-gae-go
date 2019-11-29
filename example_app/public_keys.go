package exampleservice

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	// RSAPublicKeysURL for getting google RSA public keys
	RSAPublicKeysURL = "https://www.googleapis.com/oauth2/v1/certs"
	// ECDSAPublicKeysURL  for getting ECDSA pubclic keys
	ECDSAPublicKeysURL = "https://www.gstatic.com/iap/verify/public_key"
)

func getPublicKeys(publicKeysURL string) (map[string]string, error) {

	var keys map[string]string
	client := &http.Client{Timeout: 10 * time.Second}

	r, err := client.Get(publicKeysURL)
	if err != nil {
		return keys, err
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&keys); err != nil {
		return keys, err
	}

	return keys, nil
}

func publicKeysFunction() (func(token *jwt.Token) (interface{}, error), error) {
	return func(token *jwt.Token) (interface{}, error) {

		alg := token.Method.Alg()
		kid := token.Header["kid"].(string)

		switch alg {
		case jwt.SigningMethodES256.Alg():
			publicKeys, err := getPublicKeys(ECDSAPublicKeysURL)
			if err != nil {
				return nil, err
			}

			var keys = map[string]*ecdsa.PublicKey{}
			for k, v := range publicKeys {
				if len(v) != 0 {
					parsedKey, err := jwt.ParseECPublicKeyFromPEM([]byte(v))
					if err != nil {
						return nil, err
					}
					keys[k] = parsedKey
				}
			}

			return keys[kid], nil
		default:
			publicKeys, err := getPublicKeys(RSAPublicKeysURL)
			if err != nil {
				return nil, err
			}

			var keys = map[string]*rsa.PublicKey{}
			for k, v := range publicKeys {
				if len(v) != 0 {
					parsedKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(v))
					if err != nil {
						return nil, err
					}
					keys[k] = parsedKey
				}
			}

			return keys[kid], nil
		}

	}, nil
}
