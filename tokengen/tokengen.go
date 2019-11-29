package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	iapClientDefaultTimeoutMsec = 15000
	jwtAuthorizationTimeoutMsec = 2000
	jwtUrnGrantType             = "urn:ietf:params:oauth:grant-type:jwt-bearer"
	jwtAuthorizationURL         = "https://oauth2.googleapis.com/token"
)

type customClaims struct {
	Email          string `json:"email"`
	TargetAudience string `json:"target_audience"`
	*jwt.StandardClaims
}

func readCredsFile(path string) (map[string]interface{}, error) {

	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var res map[string]interface{}
	json.Unmarshal([]byte(byteValue), &res)

	return res, nil
}

func createJWT(audienceString string, credentialsPath string, algorithm string) (string, error) {

	creds, err := readCredsFile(credentialsPath)
	if err != nil {
		return "", err
	}

	saEmail := creds["client_email"].(string)
	tokenURI := creds["token_uri"].(string)

	claims := &customClaims{
		Email:          saEmail,
		TargetAudience: audienceString,
		StandardClaims: &jwt.StandardClaims{
			Issuer:    saEmail,
			Subject:   saEmail,
			Audience:  tokenURI,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Second * 3600).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(algorithm), claims)

	switch algorithm {
	case jwt.SigningMethodRS256.Alg():
		privKey, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(creds["private_key"].(string)))
		return token.SignedString(privKey)
	case jwt.SigningMethodES256.Alg():
		privKey, _ := jwt.ParseECPrivateKeyFromPEM([]byte(creds["private_key"].(string)))
		return token.SignedString(privKey)
	}

	return "", errors.New("Cannot find JWT algorithm. Specify 'ES256' or 'RS256'")
}

func resignToken(accessToken string) (string, error) {

	var (
		res map[string]interface{}
	)

	signClient := &http.Client{
		Timeout: time.Millisecond * jwtAuthorizationTimeoutMsec,
	}

	d := url.Values{}
	d.Set("grant_type", jwtUrnGrantType)
	d.Add("assertion", accessToken)

	req, err := http.NewRequest("POST", jwtAuthorizationURL, strings.NewReader(d.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := signClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal([]byte(body), &res)
	if err != nil {
		return "", err
	}

	if res["id_token"] != nil {
		return res["id_token"].(string), nil
	}

	return "", nil
}

func main() {

	var (
		algorithm         string
		credsPath         string
		iapAudienceString string
		token             string
		resignedToken     string
		err               error
	)

	// Check if crdes path has been set
	if credsPath = os.Getenv("SERVICE_ACCOUNT_JSON_PATH"); credsPath == "" {
		panic("No SERVICE_ACCOUNT_JSON_PATH has been set")
	}

	// Check if audience string has been set
	if iapAudienceString = os.Getenv("IAP_AUDIENCE_STRING"); iapAudienceString == "" {
		panic("No IAP_AUDIENCE_STRING has been set")
	}

	// Check if algorithm string has been set
	if algorithm = os.Getenv("ALG"); algorithm == "" {
		panic("No ALG has been set")
	}
	// Generate token
	if token, err = createJWT(iapAudienceString, credsPath, algorithm); err != nil {
		panic(err.Error())
	}

	// Re-sign token
	if resignedToken, err = resignToken(token); err != nil {
		panic(err.Error())
	}

	fmt.Printf("%s\n", resignedToken)
}
