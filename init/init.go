// inicializa la base de datos
package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./database/app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// leer archivo SQL
	queryBytes, err := os.ReadFile("./init/esquema.sql")
	if err != nil {
		log.Fatal(err)
	}

	// ejecutar todo el script
	_, err = db.Exec(string(queryBytes))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Esquema ejecutado correctamente")
}
