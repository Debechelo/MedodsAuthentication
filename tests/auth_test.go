package tests

import (
	"MedodsAuthentication/pkg"
	"github.com/gbrlsnchs/jwt/v3"
	"testing"
)

func TestGenerateAccessToken(t *testing.T) {
	userID := "test-user"
	ip := "192.168.1.1"
	hs512 := jwt.NewHS512([]byte("secret_key"))

	t.Run("Проверка генерации токена", func(t *testing.T) {
		token, err := pkg.GenerateAccessToken(userID, ip)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if token == "" {
			t.Errorf("Expected a token, got empty string")
		}

		// Проверяем валидность сгенерированного токена
		var claims pkg.AccessToken
		_, err = jwt.Verify([]byte(token), hs512, &claims)
		if err != nil {
			t.Errorf("Failed to verify token: %v", err)
		}

		if claims.UserID != userID {
			t.Errorf("Expected user ID %s, got %s", userID, claims.UserID)
		}
		if claims.IP != ip {
			t.Errorf("Expected IP %s, got %s", ip, claims.IP)
		}
	})

	t.Run("GenerateAccessToken с пустым userID", func(t *testing.T) {
		token, err := pkg.GenerateAccessToken("", ip)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if token == "" {
			t.Errorf("Expected a token, got empty string")
		}
	})

	t.Run("GenerateAccessToken с пустым IP", func(t *testing.T) {
		token, err := pkg.GenerateAccessToken(userID, "")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if token == "" {
			t.Errorf("Expected a token, got empty string")
		}
	})

	t.Run("GenerateAccessToken с пустым IP и пустым userID", func(t *testing.T) {
		token, err := pkg.GenerateAccessToken("", "")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if token == "" {
			t.Errorf("Expected a token, got empty string")
		}
	})
}

func TestGenerateRefreshToken(t *testing.T) {
	t.Run("Проверка генерации токена", func(t *testing.T) {
		token, err := pkg.GenerateRefreshToken()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if token == "" {
			t.Errorf("Expected a token, got empty string")
		}

		// Проверка длины токена, если это важно
		if len(token) != 64 { // Поскольку вы используете 32 байта, закодированных в hex, итоговая длина будет 64
			t.Errorf("Expected token length 64, got %d", len(token))
		}
	})
}
