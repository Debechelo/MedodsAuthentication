package pkg

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
)

var hs512 = jwt.NewHS512([]byte("secret_key"))

// AccessToken - структура данных для JWT токена
type AccessToken struct {
	UserID string `json:"user_id"`
	IP     string `json:"ip"`
	jwt.Payload
}

type RefreshToken struct {
	jwt.Payload
}

// GenerateAccessToken генерирует Access токен
func GenerateAccessToken(userID, ip string) (string, error) {
	now := time.Now()
	claims := AccessToken{
		UserID: userID,
		IP:     ip,
		Payload: jwt.Payload{
			IssuedAt:       jwt.NumericDate(now),
			ExpirationTime: jwt.NumericDate(now.Add(15 * time.Minute)),
		},
	}

	token, err := jwt.Sign(claims, hs512)
	if err != nil {
		log.Printf("Error signing access token: %v", err)
		return "", err
	}

	log.Printf("Access token generated for user %s", userID)
	return string(token), nil
}

// GenerateRefreshToken генерирует Refresh токен
func GenerateRefreshToken() (string, error) {
	now := time.Now()
	claims := RefreshToken{
		Payload: jwt.Payload{
			IssuedAt:       jwt.NumericDate(now),
			ExpirationTime: jwt.NumericDate(now.Add(15 * time.Minute)),
		},
	}

	token, err := jwt.Sign(claims, hs512)
	if err != nil {
		log.Printf("Error signing refresh token: %v", err)
		return "", err
	}

	_, err = makeHashToken(token) //hashedToken
	if err != nil {
		log.Printf("Error hashing refresh token: %v", err)
		return "", err
	}

	log.Println("Refresh token generated")
	return string(token), nil
}

func makeHashToken(token []byte) (string, error) {
	hashedToken, err := bcrypt.GenerateFromPassword(token, bcrypt.DefaultCost)
	return string(hashedToken), err
}
