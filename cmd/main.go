package main

import (
	"MedodsAuthentication/pkg"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()

	// Роуты  для работы с токенами
	router.POST("/token", pkg.GenerateTokensHandler)
	//router.GET("/token", pkg.RefreshTokenHandler())

	// Запуск сервера
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
