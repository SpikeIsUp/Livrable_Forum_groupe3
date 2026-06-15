package router

import (
	"net/http"
	// Remplace par le bon chemin de ton module
	"forum-projet/backend/api/controller"
)

func InitialiserRoutes() {
	http.HandleFunc("/", controller.HomeHandler)

	http.HandleFunc("/creer_sujet", controller.CreateTopicHandler)

	http.HandleFunc("/sujet", controller.SujetHandler)

	fs := http.FileServer(http.Dir("./Frontend/Static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

}
