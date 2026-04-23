// declaramos funciones auxiliares que usaremos en los test
package dbmanager

import (
	"database/sql"
	"math/rand"
)

func randomString(n int) string {

	const letras = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, n)
	for i := range b {
		b[i] = letras[rand.Intn(len(letras))]
	}
	return string(b)
}

func cleanDB(db *sql.DB, tablaALimpiar string) {

	_, err := db.Exec("DELETE FROM " + tablaALimpiar)
	if err != nil {
		panic("error limpiando tabla " + tablaALimpiar + ": " + err.Error())
	}

}

func generarGestorDbLimpio(tablaLimpia string) *sql.DB {

	db := NewGestotorDb("../../../database/app_test.db")
	cleanDB(db, tablaLimpia)

	return db
}
