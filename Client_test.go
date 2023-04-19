package authorizer

import (
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestClient_GetAuthUrl(t *testing.T) {
	clientName := "test-client"
	baseHost := "http://authorizer/"
	secret := "testSuperSecret123!"
	userId := "12399"
	redirectUri := "https://example.com/redirect"

	t.Run("success", func(t *testing.T) {
		expectedPost := "client=" + clientName + "&client_user_id=" + userId + "&redirect_uri=" + url.QueryEscape(redirectUri)
		expectedOauthUrl := "https://auth.kneu.edu.ua/oauth?response_type=code&client_id=0&redirect_uri=https%3A%2F%2Fpigeon.com%2Fcomplete&_state_"

		gock.New(baseHost).
			Post("/url").
			BasicAuth(Username, secret).
			MatchType("url").
			BodyString(expectedPost).
			Reply(200).
			JSON(GetAuthUrlResponse{
				AuthUrl: expectedOauthUrl,
			})

		client := Client{
			Host:       baseHost,
			Secret:     secret,
			ClientName: clientName,
		}

		actualAuthUrl, err := client.GetAuthUrl(userId, redirectUri)

		assert.Equal(t, expectedOauthUrl, actualAuthUrl)
		assert.NoError(t, err)
	})

	t.Run("error empty url", func(t *testing.T) {
		expectedPost := "client=" + clientName + "&client_user_id=" + userId

		gock.New(baseHost).
			Post("/url").
			BasicAuth(Username, secret).
			MatchType("url").
			BodyString(expectedPost).
			Reply(200).
			JSON("{}")

		client := Client{
			Host:       baseHost,
			Secret:     secret,
			ClientName: clientName,
		}

		actualAuthUrl, err := client.GetAuthUrl(userId, "https://example.com/redirect")

		assert.Error(t, err)
		assert.Equal(t, "fail to get auth url", err.Error())
		assert.Empty(t, actualAuthUrl)
	})

	t.Run("error json", func(t *testing.T) {
		expectedPost := "client=" + clientName + "&client_user_id=" + userId

		gock.New(baseHost).
			Post("/url").
			BasicAuth(Username, secret).
			MatchType("url").
			BodyString(expectedPost).
			Reply(500).
			BodyString("Server error!")

		client := Client{
			Host:       baseHost,
			Secret:     secret,
			ClientName: clientName,
		}

		actualAuthUrl, err := client.GetAuthUrl(userId, "https://example.com/redirect")

		assert.Error(t, err)
		assert.Equal(t, "Request failed: 500 Internal Server Error", err.Error())
		assert.Empty(t, actualAuthUrl)
	})

}
