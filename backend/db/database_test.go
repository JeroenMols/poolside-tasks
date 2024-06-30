package db

import "testing"

func TestDatabase_Authorize(t *testing.T) {
	database := InMemoryDatabase()

	const validAccessToken = "f2d869a8-e5bc-4fbf-ad71-111111111111"
	const nonExistingAccessToken = "f2d869a8-e5bc-4fbf-ad71-2222222222222"

	database.AccessTokens[validAccessToken] = "valid_account"

	t.Run("valid token", func(t *testing.T) {
		_, err := database.Authorize(validAccessToken)
		if err != nil {
			t.Errorf("expected nil, got %s", err)
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		_, err := database.Authorize("not-a-uuid")
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("account doesnt exist", func(t *testing.T) {
		_, err := database.Authorize(nonExistingAccessToken)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}
