package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

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
