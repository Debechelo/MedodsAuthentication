package pkg

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// GenerateTokensHandler
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

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		log.Printf("Error generating access token for user %s: %v", userID, err)
		http.Error(w, "Error generating access token", http.StatusInternalServerError)
		return
	}

	log.Printf("Tokens generated for user %s", userID)

	// Ответ клиенту
	response := map[string]string{
		"access_token":  accessToken,
		"refresh_Token": refreshToken,
	}
	json.NewEncoder(w).Encode(response)
}
