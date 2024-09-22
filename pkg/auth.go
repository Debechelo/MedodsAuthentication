package pkg

import (
	"time"

	"github.com/gbrlsnchs/jwt/v3"
)

var hs512 = jwt.NewHS512([]byte("secret_key"))

// Payload - структура данных для JWT токена
type Payload struct {
	UserID string `json:"user_id"`
	IP     string `json:"ip"`
	jwt.Payload
}

// GenerateAccessToken генерирует Access токен
func GenerateAccessToken(userID, ip string) (string, error) {
	now := time.Now()
	claims := Payload{
		UserID: userID,
		IP:     ip,
		Payload: jwt.Payload{
			IssuedAt:       jwt.NumericDate(now),
			ExpirationTime: jwt.NumericDate(now.Add(15 * time.Minute)),
		},
	}

	token, err := jwt.Sign(claims, hs512)
	if err != nil {
		return "", err
	}

	return string(token), nil
}
