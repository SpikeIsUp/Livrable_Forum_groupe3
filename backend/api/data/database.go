package data

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
	dsn := "root:root@tcp(127.0.0.1:3306)/forum?parseTime=true"

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erreur préparation BDD :", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Impossible de joindre MySQL :", err)
	}

	fmt.Println("✅ Connexion à MySQL (forum) réussie !")
}
