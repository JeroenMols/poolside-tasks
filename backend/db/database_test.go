package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDatabase_GetAccessToken(t *testing.T) {
	database := InMemoryDatabase()

	const validAccessToken = "f2d869a8-e5bc-4fbf-ad71-111111111111"
	const nonExistingAccessToken = "f2d869a8-e5bc-4fbf-ad71-2222222222222"

	fakeToken := AccessToken{AccountNumber: "valid_account", Token: validAccessToken}
	database.AccessTokens[validAccessToken] = fakeToken

	t.Run("valid token", func(t *testing.T) {
		accessToken, err := database.GetAccessToken(validAccessToken)
		assert.Nil(t, err)
		assert.Equal(t, &fakeToken, accessToken)
	})

	t.Run("invalid token", func(t *testing.T) {
		accessToken, err := database.GetAccessToken("not-a-uuid")
		assert.NotNil(t, err)
		assert.Nil(t, accessToken)
	})

	t.Run("account doesnt exist", func(t *testing.T) {
		accessToken, err := database.GetAccessToken(nonExistingAccessToken)
		assert.NotNil(t, err)
		assert.Nil(t, accessToken)
	})
}
