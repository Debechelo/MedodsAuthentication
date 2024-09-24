package pkg

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
}

// GenerateTokensHandler генерация access и refresh токенов
func GenerateTokensHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := r.URL.Query().Get("user_id")
	ip := r.RemoteAddr

	log.Printf("Generating tokens for user ID: %s, IP: %s", userID, ip)

	accessToken, err := GenerateAccessToken(userID, ip)
	if err != nil {
		log.Printf("Error generating access token for user %s: %v", userID, err)
		http.Error(w, "Error generating access token", http.StatusInternalServerError)
		return
	}

	hashAccessToken := hashWithSHA512(accessToken)

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		log.Printf("Error generating access token for user %s: %v", userID, err)
		http.Error(w, "Error generating access token", http.StatusInternalServerError)
		return
	}

	err = saveRefreshToken(userID, ip, refreshToken, w)
	if err != nil {
		log.Printf("Error saving refresh token to DB for user %s: %v", userID, err)
		http.Error(w, "Error saving refresh token", http.StatusInternalServerError)
		return
	}

	log.Printf("Tokens generated for user %s", userID)

	// Ответ клиенту
	response := map[string]string{
		"access_token":  hashAccessToken,
		"refresh_Token": refreshToken,
	}
	json.NewEncoder(w).Encode(response)
}

// RefreshTokenHandler обновляет Access токен с использованием Refresh токена
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Чтение тела запроса
	var req RefreshTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.RefreshToken == "" || req.UserID == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	refreshToken := req.RefreshToken
	userID := req.UserID
	newIP := r.RemoteAddr

	// Получаем хеш Refresh токена и IP-адрес из базы данных для пользователя
	hashedTokenFromDB, ipAddressFromDB, err := GetRefreshTokenFromDB(userID)
	if err != nil {
		log.Printf("Error retrieving refresh token from DB for user %s: %v", userID, err)
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Проверяем валидность переданного Refresh токена
	err = bcrypt.CompareHashAndPassword([]byte(hashedTokenFromDB), []byte(refreshToken))
	if err != nil {
		log.Printf("Invalid refresh token for user %s: %v", userID, err)
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Сравниваем IP-адреса
	if ipAddressFromDB != newIP {

		log.Printf("IP address changed for user %s: old IP %s, new IP %s", userID, ipAddressFromDB, newIP)
		sendEmailWarning(userID, ipAddressFromDB, newIP) // Отправка предупреждения на email
	}

	// Генерация нового Access токена
	newAccessToken, err := GenerateAccessToken(userID, newIP)
	if err != nil {
		log.Printf("Error generating new access token for user %s: %v", userID, err)
		http.Error(w, "Error generating new access token", http.StatusInternalServerError)
		return
	}

	// Ответ клиенту с новым Access токеном
	response := map[string]string{
		"access_token": newAccessToken,
	}
	json.NewEncoder(w).Encode(response)
	log.Printf("New access token generated for user %s", userID)
}

// Моковая функция для отправки email
func sendEmailWarning(userID, oldIP, newIP string) {
	log.Printf("Warning email sent to user %s: IP address changed from %s to %s", userID, oldIP, newIP)
}

func saveRefreshToken(userID, ip, refreshToken string, w http.ResponseWriter) error {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost) //hashedToken
	if err != nil {
		log.Printf("Error hashing refresh token: %v", err)
		return err
	}

	err = SaveRefreshTokenToDB(userID, ip, string(hashedToken))
	if err != nil {
		log.Printf("Error saving refresh token to DB for user %s: %v", userID, err)
		http.Error(w, "Error saving refresh token", http.StatusInternalServerError)
		return err
	}

	return nil
}

func hashWithSHA512(token string) string {
	hasher := sha512.New()
	hasher.Write([]byte(token))
	hashedBytes := hasher.Sum(nil)

	return hex.EncodeToString(hashedBytes)
}
