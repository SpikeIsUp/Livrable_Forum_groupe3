package router

import (
	"forum-projet/backend/api/controller"
	"net/http"
)

func InitialiserRoutes() {
	http.HandleFunc("/api/categories", controller.GetCategories)
	http.HandleFunc("/api/sujets", controller.CreerSujet)
	http.HandleFunc("/api/connexion", controller.Connexion)
	http.HandleFunc("/api/messages/envoyer", controller.EnvoyerMessage)
	http.HandleFunc("/api/messages/lire", controller.GetMessages)
	http.HandleFunc("/api/inscription", controller.Inscription)
	http.HandleFunc("/api/categories/posts", controller.GetPostsParCategorie)
	http.HandleFunc("/api/commentaires", controller.AjouterCommentaire)
	http.HandleFunc("/api/commentaires/lire", controller.GetCommentaires)
}
