package helper

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	logDump "github.com/hlmn/senyum-go-utils/logger/echo"
	"github.com/sirupsen/logrus"
)

type BRIAPI struct {
	URL          string
	ClientId     string
	ClientSecret string
	Proxy        string
}

type ResGetTokenBri struct {
	AccessToken string `json:"access_token"`
}

func (api *BRIAPI) GetToken() (token string, err error) {

	ctx := context.Background()

	// timeout, err := strconv.Atoi(h.Config.RequestTimeout)
	// if err != nil {
	// 	logDump.Info(nil, logrus.Fields{"error": err}, "timeout")

	// 	return token, err
	// }

	client := &http.Client{}
	// client := &http.Client{Timeout: time.Second * time.Duration(timeout)}

	if api.Proxy != "" {
		proxyUrl, err := url.Parse(api.Proxy)
		if err != nil {

			logDump.Info(nil, logrus.Fields{"error": err}, "proxy")

			return token, err
		}

		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	}

	data := url.Values{}
	data.Set("client_id", api.ClientId)
	data.Set("client_secret", api.ClientSecret)

	url := api.URL + "/oauth/client_credential/accesstoken?grant_type=client_credentials"
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		logDump.Info(nil, logrus.Fields{"tag": "Helper",
			"job":      "postGenerateTokenResponse1",
			"url":      url,
			"header":   req.Header,
			"body":     data.Encode(),
			"response": err.Error()}, "generate token")

		return token, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		logDump.Info(nil, logrus.Fields{
			"tag":      "Helper",
			"job":      "postGenerateTokenResponse2",
			"url":      url,
			"header":   req.Header,
			"body":     data.Encode(),
			"response": err.Error(),
		}, "generate token")

		return token, err
	}

	defer res.Body.Close()

	bytesGet, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logDump.Info(nil, logrus.Fields{
			"tag":      "Helper",
			"job":      "postGenerateTokenResponse3",
			"url":      url,
			"header":   res.Header,
			"body":     data.Encode(),
			"response": err.Error(),
		}, "generate token")

		return token, err
	}

	if res.StatusCode != 200 {
		if err != nil {
			logDump.Info(nil, logrus.Fields{
				"tag":      "Helper",
				"job":      "postGenerateTokenResponse4",
				"url":      url,
				"header":   res.Header,
				"body":     data.Encode(),
				"response": string(bytesGet),
			}, "generate token failed, non 200")

			return token, err
		}
	}

	logDump.Info(nil, logrus.Fields{
		"tag":      "Helper",
		"job":      "postGenerateTokenResponse4A",
		"url":      url,
		"header":   res.Header,
		"body":     data.Encode(),
		"response": string(bytesGet),
	}, "generate token")

	var resultData ResGetTokenBri
	json.Unmarshal(bytesGet, &resultData)

	return resultData.AccessToken, err

}

func (api *BRIAPI) signHmacBase64(payload string) (*string, error) {
	key := []byte(api.ClientSecret)

	mac := hmac.New(sha256.New, key)

	_, err := mac.Write([]byte(payload))
	if err != nil {
		return nil, err
	}

	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return &sign, nil
}

func (api *BRIAPI) GenerateSignature(requestBody string, token string, path string, verb string) (headers map[string]string) {
	headers = make(map[string]string)
	dateTime := time.Now().UTC().Format("2006-01-02T15:04:05.515Z")

	payload := fmt.Sprintf("path=%v&verb=%v&token=Bearer %v&timestamp=%v&body=%v",
		path,
		verb,
		token,
		dateTime,
		requestBody,
	)

	sign, _ := api.signHmacBase64(payload)

	headers["BRI-Signature"] = *sign
	headers["BRI-Timestamp"] = dateTime

	return headers
}
