package exampleservice

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"

	"google.golang.org/appengine/log"
)

const (
	jwtGAETokenHeader           = "x-goog-iap-jwt-assertion"
	standardAuthorizationHeader = "authorization"
	baererString                = "BEARER "
)

// CustomClaims for holding google-specific claims
type CustomClaims struct {
	Email string `json:"email"`
	*jwt.StandardClaims
}

// TokenExtractor is a custom token extractor interface
type TokenExtractor struct {
	request.Extractor
}

// ExtractToken implements extraction of GEA repacked auth token or a standard one
func (e *TokenExtractor) ExtractToken(r *http.Request) (string, error) {
	for headerKey, headerValue := range r.Header {
		if strings.ToLower(headerKey) == jwtGAETokenHeader && len(headerValue) > 0 {
			return headerValue[0], nil
		} else if strings.ToLower(headerKey) == standardAuthorizationHeader {
			if len(headerValue[0]) > 6 && strings.ToUpper(headerValue[0][0:7]) == baererString {
				return headerValue[0][7:], nil
			}
			return "", nil
		}
	}
	return "", nil
}

// Auth is middleware to authenticate http requests.
func (s *Service) Auth(next http.HandlerFunc) http.HandlerFunc {

	keyFunction, err := publicKeysFunction()
	if err != nil {
		// can't start without public keys
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var (
			token *jwt.Token
			err   error
		)

		// Parse claims
		claims := &CustomClaims{}
		extractor := &TokenExtractor{}

		if token, err = request.ParseFromRequest(r, extractor, keyFunction, request.WithClaims(claims)); err != nil {
			s.writeResponseData(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Validate jwt token and GAE specific claims
		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Log the requester
		log.Infof(r.Context(), "URL requested by: %s", claims.Email)

		// Call the next handler
		next(w, r)
	}
}
