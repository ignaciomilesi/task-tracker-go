package dbmanager

import (
	"database/sql"
	"fmt"

	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type solicitanteRepository struct {
	db *sql.DB
}

func NewsolicitanteRepository(db *sql.DB) *solicitanteRepository {
	return &solicitanteRepository{db: db}
}

// Crear nuevo solicitante, devuelve el id creado o, en caso que ya exista,
// devuelve el id del solicitante existente
func (r *solicitanteRepository) Crear(nombre string) (int64, error) {
	res, err := r.db.Exec(
		"INSERT OR IGNORE INTO solicitante (nombre) VALUES (?)",
		nombre,
	)
	if err != nil {
		return 0, fmt.Errorf("Error inesperado, detalle: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Error inesperado, detalle: %v", err)
	}

	// si no insertó (ya existía), buscamos el ID
	if id == 0 {
		existingID, err := r.ObtenerIDPorNombre(nombre)
		if err != nil {
			return 0, err
		}
		return int64(existingID), nil
	}

	return id, nil
}

func (r *solicitanteRepository) ObtenerIDPorNombre(nombre string) (int, error) {
	var id int

	err := r.db.QueryRow( // toma solo el primer resultado
		"SELECT id FROM solicitante WHERE nombre = ?",
		nombre,
	).Scan(&id)

	if err == sql.ErrNoRows {
		return 0, appErrors.SolicitanteNoEncontrado
	}
	if err != nil {
		return 0, fmt.Errorf("Error inesperado, detalle: %v", err)
	}

	return id, nil
}

// Busca en la lista. Si el parámetro es vacío, trae la lista completa
func (r *solicitanteRepository) Buscar(parametro string) ([]models.Solicitante, error) {
	rows, err := r.db.Query("SELECT id, nombre FROM solicitante WHERE nombre LIKE ?", "%"+parametro+"%")
	if err != nil {
		return nil, fmt.Errorf("Error inesperado, detalle: %v", err)
	}
	defer rows.Close()

	var lista []models.Solicitante

	for rows.Next() {
		var s models.Solicitante

		err := rows.Scan(&s.ID, &s.Nombre)
		if err != nil {
			return nil, fmt.Errorf("Error inesperado, detalle: %v", err)
		}

		lista = append(lista, s)
	}

	return lista, nil
}
