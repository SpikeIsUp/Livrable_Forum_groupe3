package main

import (
	"fmt"
	"net/http"

	database "forum-projet/backend/api/data"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bienvenue sur l'API de notre Forum !")
}

func main() {
	// 1. On lance la connexion à la base de données en premier
	database.Connect()

	// 2. On configure nos routes web
	http.HandleFunc("/", homeHandler)

	// 3. On démarre le serveur
	fmt.Println("🚀 Le serveur démarre sur le port 8080...")
	http.ListenAndServe(":8080", nil)
}
