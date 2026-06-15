package data

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Categorie struct {
	ID          int
	Nom         string
	Description string
}

type Sujet struct {
	ID          int
	Titre       string
	Contenu     string
	AuteurID    int
	CategorieID int
}

type Message struct {
	ID       int
	Contenu  string
	AuteurID int
	SujetID  int
}

var DB *sql.DB

func Connect() {
	// Remplacez ces valeurs par les vôtres (utilisateur, mot de passe, nom de la base)
	// Format : "utilisateur:mot_de_passe@tcp(127.0.0.1:3306)/nom_de_la_base"
	dsn := "projectforum:password@tcp(127.0.0.1:3306)/forum_projet?parseTime=true"

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erreur lors de l'ouverture de la base de données :", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Impossible de joindre la base de données :", err)
	}

	fmt.Println("✅ Connexion à MySQL réussie !")
}

func GetCategories() ([]Categorie, error) {
	rows, err := DB.Query("SELECT id, nom, description FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Categorie
	for rows.Next() {
		var c Categorie
		if err := rows.Scan(&c.ID, &c.Nom, &c.Description); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil

}

func InsertTopic(titre, contenu string, auteurID, categorieID int) error {
	query := "INSERT INTO sujets (titre, contenu, auteur_id, categorie_id) VALUES (?, ?, ?, ?)"
	_, err := DB.Exec(query, titre, contenu, auteurID, categorieID)
	return err
}

func GetTopicByID(id int) (Sujet, error) {
	// Stub: à implémenter avec une vraie requête SQL si besoin.
	return Sujet{ID: id}, nil
}

func GetMessagesByTopicID(topicID int) ([]Message, error) {
	// Stub: à implémenter avec une vraie requête SQL si besoin.
	return []Message{}, nil
}
