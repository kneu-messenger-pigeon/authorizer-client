package authorizer

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

const Username = "pigeon"

type ClientInterface interface {
	GetAuthUrl(userId string) (string, error)
}

type Client struct {
	Host       string
	Secret     string
	ClientName string
}

type GetAuthUrlResponse struct {
	AuthUrl string `json:"authUrl" binding:"required"`
}

func (client *Client) GetAuthUrl(userId string) (string, error) {
	getAuthUrlResponse := GetAuthUrlResponse{}

	req, _ := http.NewRequest(
		http.MethodPost,
		client.Host+"/url",
		strings.NewReader(
			"client="+client.ClientName+"&client_user_id="+userId,
		),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(Username, client.Secret)

	response, err := http.DefaultClient.Do(req)

	if err == nil && response.StatusCode != http.StatusOK {
		err = errors.New("Request failed: " + response.Status)
	}

	if err == nil {
		err = json.NewDecoder(response.Body).Decode(&getAuthUrlResponse)
	}

	if err == nil && getAuthUrlResponse.AuthUrl == "" {
		err = errors.New("fail to get auth url")
	}

	return getAuthUrlResponse.AuthUrl, err
}
