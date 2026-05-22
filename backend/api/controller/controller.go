package controller

import (
	"html"
	"html/template"
	"net/http"
	"strconv"

	"forum-projet/backend/api/data"
)

var tmpl = template.New("tmpl")

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := data.GetCategories()
	if err != nil {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// On suppose que 'tmpl' est ton template pré-chargé (ex: Index.html)
	tmpl.ExecuteTemplate(w, "Index.html", categories)
}

func CreateTopicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()

	// Nettoyage XSS : on échappe les caractères HTML dangereux (<, >, etc.)
	titreBrut := r.FormValue("titre")
	contenuBrut := r.FormValue("contenu")

	titreClean := html.EscapeString(titreBrut)
	contenuClean := html.EscapeString(contenuBrut)

	// ID de l'auteur (en dur pour l'instant, sera remplacé par la session de l'Étudiant 2)
	auteurID := 1
	categorieID := 1 // À récupérer via le formulaire également

	err := data.InsertTopic(titreClean, contenuClean, auteurID, categorieID)
	if err != nil {
		http.Error(w, "Erreur lors de la création", http.StatusInternalServerError)
		return
	}

	// Redirection vers l'accueil après création
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Structure pour combiner les données à envoyer au template
type SujetPageData struct {
    Sujet    data.Sujet
    Messages []data.Message
}

func SujetHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID dans l'URL
	idStr := r.URL.Query().Get("id")
	sujetID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// Récupérer le sujet complet
    sujet, _ := data.GetTopicByID(sujetID)
    
    // Récupérer les messages liés
    messages, _ := data.GetMessagesByTopicID(sujetID)
	// Préparer les données pour le template
	pageData := SujetPageData{
		Sujet:    sujet,
		Messages: messages,
	}

	tmpl.ExecuteTemplate(w, "Sujet.html", pageData)
}
