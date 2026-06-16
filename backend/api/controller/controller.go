package controller

import (
	"database/sql"
	"encoding/json"
	"html"
	"net/http"
	"strconv"

	"forum-projet/backend/api/data"
)

type Categorie struct {
	ID    int    `json:"id"`
	Titre string `json:"titre"`
}

type PostAffiche struct {
	ID         int    `json:"id"`
	Titre      string `json:"titre"`
	Corps      string `json:"corps"`
	Auteur     string `json:"auteur"`
	Date       string `json:"date"`
	NbLikes    int    `json:"nb_likes"`
	NbDislikes int    `json:"nb_dislikes"`
}

type NouveauSujet struct {
	Titre         string `json:"titre"`
	Corps         string `json:"corps"`
	IDCategorie   int    `json:"id_categorie"`
	IDUtilisateur int    `json:"id_utilisateur"`
}

type InfosLogin struct {
	Pseudo string `json:"pseudo"`
	Mdp    string `json:"mdp"`
}

type ProfilUtilisateur struct {
	ID        int    `json:"id"`
	Pseudo    string `json:"pseudo"`
	ProfilPic string `json:"profil_pic"`
	Bio       string `json:"bio"`
}

type NouveauMessage struct {
	IDExpediteur       int    `json:"id_expediteur"`
	PseudoDestinataire string `json:"pseudo_destinataire"`
	Contenu            string `json:"contenu"`
}

type MessageRecu struct {
	IDMsg      int    `json:"id_msg"`
	Expediteur string `json:"expediteur"`
	Contenu    string `json:"contenu"`
	DateEnvoi  string `json:"date_envoi"`
}

// ==========================================
// OUTIL DE SÉCURITÉ CORS
// ==========================================
func configurerCORS(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Si c'est juste le navigateur qui vérifie si la route est sûre (OPTIONS), on lui dit OK et on coupe ici
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return true
	}
	return false
}

// ==========================================
// ROUTES : CATÉGORIES ET SUJETS
// ==========================================

func GetCategories(w http.ResponseWriter, r *http.Request) {
	if configurerCORS(w, r) {
		return
	}

	lignes, _ := data.DB.Query("SELECT id_categorie, nom_categorie FROM categorie")
	defer lignes.Close()

	var categories []Categorie
	for lignes.Next() {
		var c Categorie
		lignes.Scan(&c.ID, &c.Titre)
		categories = append(categories, c)
	}
	json.NewEncoder(w).Encode(categories)
}

func ToggleLike(w http.ResponseWriter, r *http.Request) {
	if configurerCORS(w, r) {
		return
	}

	uid := r.URL.Query().Get("uid")
	pid := r.URL.Query().Get("pid")
	typeReq := r.URL.Query().Get("type") // On récupère le 1 ou -1

	var idLike, typeActuel int
	err := data.DB.QueryRow("SELECT id_like, type_reaction FROM likes WHERE id_utilisateur = ? AND id_post = ?", uid, pid).Scan(&idLike, &typeActuel)

	typeReqInt, _ := strconv.Atoi(typeReq)

	if err == sql.ErrNoRows {
		data.DB.Exec("INSERT INTO likes (id_utilisateur, id_post, type_reaction) VALUES (?, ?, ?)", uid, pid, typeReqInt)
	} else {
		if typeActuel == typeReqInt {
			data.DB.Exec("DELETE FROM likes WHERE id_like = ?", idLike)
		} else {
			data.DB.Exec("UPDATE likes SET type_reaction = ? WHERE id_like = ?", typeReqInt, idLike)
		}
	}
	w.WriteHeader(http.StatusOK)
}

func GetPostsParCategorie(w http.ResponseWriter, r *http.Request) {
	if configurerCORS(w, r) {
		return
	}

	idCat := r.URL.Query().Get("id")
	// Requête SQL intelligente qui sépare les +1 et les -1
	requete := `
		SELECT p.id_post, p.titre, p.corps, u.pseudo, p.date_, 
		       COALESCE(SUM(CASE WHEN l.type_reaction = 1 THEN 1 ELSE 0 END), 0) as nb_likes,
		       COALESCE(SUM(CASE WHEN l.type_reaction = -1 THEN 1 ELSE 0 END), 0) as nb_dislikes
		FROM post p 
		JOIN utilisateur u ON p.id_utilisateur = u.id_utilisateur 
		LEFT JOIN likes l ON p.id_post = l.id_post 
		WHERE p.id_categorie = ? 
		GROUP BY p.id_post 
		ORDER BY p.date_ DESC`

	lignes, err := data.DB.Query(requete, idCat)
	if err != nil {
		http.Error(w, "Erreur base de données", http.StatusInternalServerError)
		return
	}
	defer lignes.Close()

	var posts []PostAffiche
	for lignes.Next() {
		var p PostAffiche
		lignes.Scan(&p.ID, &p.Titre, &p.Corps, &p.Auteur, &p.Date, &p.NbLikes, &p.NbDislikes)
		posts = append(posts, p)
	}
	json.NewEncoder(w).Encode(posts)
}

func CreerSujet(w http.ResponseWriter, r *http.Request) {
	if configurerCORS(w, r) {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Seul le POST est autorisé ici", http.StatusMethodNotAllowed)
		return
	}

	var sujet NouveauSujet
	if err := json.NewDecoder(r.Body).Decode(&sujet); err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	titrePropre := html.EscapeString(sujet.Titre)
	corpsPropre := html.EscapeString(sujet.Corps)

	_, err := data.DB.Exec("INSERT INTO post (titre, corps, id_categorie, id_utilisateur, date_) VALUES (?, ?, ?, ?, NOW())", titrePropre, corpsPropre, sujet.IDCategorie, sujet.IDUtilisateur)
	if err != nil {
		http.Error(w, "Erreur enregistrement", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// ==========================================
// ROUTES : AUTHENTIFICATION
// ==========================================

func Connexion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		return
	}

	var login InfosLogin
	json.NewDecoder(r.Body).Decode(&login)

	var user ProfilUtilisateur
	var pic sql.NullString

	err := data.DB.QueryRow("SELECT id_utilisateur, pseudo, profilPic_bonus_ FROM utilisateur WHERE pseudo = ? AND mdp = ?", login.Pseudo, login.Mdp).Scan(&user.ID, &user.Pseudo, &pic)

	if err != nil {
		http.Error(w, "Pseudo ou mot de passe incorrect", http.StatusUnauthorized)
		return
	}

	if pic.Valid {
		user.ProfilPic = pic.String
	}

	json.NewEncoder(w).Encode(user)
}

func Inscription(w http.ResponseWriter, r *http.Request) {
	if configurerCORS(w, r) {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Seul le POST est autorisé ici", http.StatusMethodNotAllowed)
		return
	}

	var newUser InfosLogin
	json.NewDecoder(r.Body).Decode(&newUser)

	var existant int
	err := data.DB.QueryRow("SELECT id_utilisateur FROM utilisateur WHERE pseudo = ?", newUser.Pseudo).Scan(&existant)
	if err == nil {
		http.Error(w, "Pseudo pris", http.StatusConflict)
		return
	}

	_, err = data.DB.Exec("INSERT INTO utilisateur (pseudo, mdp) VALUES (?, ?)", newUser.Pseudo, newUser.Mdp)
	if err != nil {
		http.Error(w, "Erreur bdd", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// ==========================================
// ROUTES : MESSAGES ET LIKES (Tendances)
// ==========================================

func EnvoyerMessage(w http.ResponseWriter, r *http.Request) {
	if configurerCORS(w, r) {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Seul le POST est autorisé ici", http.StatusMethodNotAllowed)
		return
	}

	var msg NouveauMessage
	json.NewDecoder(r.Body).Decode(&msg)

	var idDest int
	err := data.DB.QueryRow("SELECT id_utilisateur FROM utilisateur WHERE pseudo = ?", msg.PseudoDestinataire).Scan(&idDest)
	if err != nil {
		http.Error(w, "Destinataire introuvable", http.StatusNotFound)
		return
	}

	data.DB.Exec("INSERT INTO messages_prives (id_expediteur, id_destinataire, contenu, date_envoi) VALUES (?, ?, ?, NOW())", msg.IDExpediteur, idDest, html.EscapeString(msg.Contenu))
	w.WriteHeader(http.StatusCreated)
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	if configurerCORS(w, r) {
		return
	}
	idU := r.URL.Query().Get("id")

	rows, _ := data.DB.Query("SELECT m.id_msg, u.pseudo, m.contenu, m.date_envoi FROM messages_prives m JOIN utilisateur u ON m.id_expediteur = u.id_utilisateur WHERE m.id_destinataire = ? ORDER BY m.date_envoi DESC", idU)
	defer rows.Close()

	var msgs []MessageRecu
	for rows.Next() {
		var m MessageRecu
		rows.Scan(&m.IDMsg, &m.Expediteur, &m.Contenu, &m.DateEnvoi)
		msgs = append(msgs, m)
	}
	json.NewEncoder(w).Encode(msgs)
}

func GetTendances(w http.ResponseWriter, r *http.Request) {
	if configurerCORS(w, r) {
		return
	}

	requete := `SELECT p.id_post, p.titre, p.corps, u.pseudo, p.date_, COUNT(l.id_like) as nb_likes
				FROM post p 
				JOIN utilisateur u ON p.id_utilisateur = u.id_utilisateur
				LEFT JOIN likes l ON p.id_post = l.id_post
				GROUP BY p.id_post ORDER BY nb_likes DESC LIMIT 10`

	rows, _ := data.DB.Query(requete)
	defer rows.Close()

	var posts []PostAffiche
	for rows.Next() {
		var p PostAffiche
		rows.Scan(&p.ID, &p.Titre, &p.Corps, &p.Auteur, &p.Date, &p.NbLikes)
		posts = append(posts, p)
	}
	json.NewEncoder(w).Encode(posts)
}

type NouveauCommentaire struct {
	IDPost        int    `json:"id_post"`
	IDUtilisateur int    `json:"id_utilisateur"`
	Contenu       string `json:"contenu"`
}

type CommentaireAffiche struct {
	ID      int    `json:"id"`
	Auteur  string `json:"auteur"`
	Contenu string `json:"contenu"`
	Date    string `json:"date"`
}

func AjouterCommentaire(w http.ResponseWriter, r *http.Request) {
	if configurerCORS(w, r) {
		return
	}

	var c NouveauCommentaire
	json.NewDecoder(r.Body).Decode(&c)

	contenuPropre := html.EscapeString(c.Contenu)

	_, err := data.DB.Exec("INSERT INTO commentaire (id_post, id_utilisateur, contenu, date_) VALUES (?, ?, ?, NOW())", c.IDPost, c.IDUtilisateur, contenuPropre)
	if err != nil {
		http.Error(w, "Erreur bdd", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func GetCommentaires(w http.ResponseWriter, r *http.Request) {
	if configurerCORS(w, r) {
		return
	}
	idPost := r.URL.Query().Get("id_post")

	lignes, _ := data.DB.Query("SELECT c.id_commentaire, u.pseudo, c.contenu, c.date_ FROM commentaire c JOIN utilisateur u ON c.id_utilisateur = u.id_utilisateur WHERE c.id_post = ? ORDER BY c.date_ ASC", idPost)
	defer lignes.Close()

	var comms []CommentaireAffiche
	for lignes.Next() {
		var c CommentaireAffiche
		lignes.Scan(&c.ID, &c.Auteur, &c.Contenu, &c.Date)
		comms = append(comms, c)
	}
	json.NewEncoder(w).Encode(comms)
}
