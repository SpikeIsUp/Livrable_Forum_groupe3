package main

import (
	"fmt"
	"net/http"

	"forum-projet/backend/api/data"
	"forum-projet/backend/api/router"
)

func main() {
	// 1. On connecte la base de données
	data.Connect()

	// 2. On charge toutes les routes (URL) de l'API
	router.InitialiserRoutes()

	// 3. On démarre le serveur
	fmt.Println("🚀 Le serveur démarre sur le port 8080...")
	http.ListenAndServe(":8080", nil)
}
