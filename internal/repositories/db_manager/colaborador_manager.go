package dbmanager

import (
	"database/sql"
	"errors"
	"fmt"

	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type colaboradorRepository struct {
	db *sql.DB
}

func NewColaboradorRepository(db *sql.DB) *colaboradorRepository {
	return &colaboradorRepository{db: db}
}

// Crear nuevo colaborador, devuelve el id creado
func (r *colaboradorRepository) Crear(NuevoColaborador *models.Colaborador) (int, error) {
	res, err := r.db.Exec(
		"INSERT INTO colaborador (nombre) VALUES (?)",
		NuevoColaborador.Nombre,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return 0, appErrors.ColaboradorDuplicado
			}
		}

		return 0, fmt.Errorf("error inesperado: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Error inesperado, detalle: %v", err)
	}

	return int(id), nil
}

func (r *colaboradorRepository) ObtenerIDPorNombre(nombre string) (int, error) {
	var id int

	err := r.db.QueryRow(
		"SELECT id FROM colaborador WHERE nombre = ?",
		nombre,
	).Scan(&id)

	if err == sql.ErrNoRows {
		return 0, appErrors.ColaboradorNoEncontrado
	}
	if err != nil {
		return 0, fmt.Errorf("Error inesperado, detalle: %v", err)
	}

	return id, nil
}

// Busca en la lista.
func (r *colaboradorRepository) Buscar(parametro string) ([]models.Colaborador, error) {
	if parametro == "" {
		return nil, appErrors.ParametroDeBusquedaVacio
	}
	rows, err := r.db.Query(
		"SELECT id, nombre FROM colaborador WHERE nombre LIKE ?",
		"%"+parametro+"%",
	)
	if err != nil {
		return nil, fmt.Errorf("Error inesperado, detalle: %v", err)
	}
	defer rows.Close()

	var lista []models.Colaborador

	for rows.Next() {
		var c models.Colaborador

		err := rows.Scan(&c.ID, &c.Nombre)
		if err != nil {
			return nil, fmt.Errorf("Error inesperado, detalle: %v", err)
		}

		lista = append(lista, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error inesperado, detalle: %v", err)
	}

	return lista, nil
}

func (r *colaboradorRepository) Listar(limit, offset int) ([]models.Colaborador, error) {

	rows, err := r.db.Query(
		`SELECT id, nombre FROM colaborador ORDER BY id
        LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("Error inesperado, detalle: %v", err)
	}
	defer rows.Close()

	var lista []models.Colaborador

	for rows.Next() {
		var c models.Colaborador

		err := rows.Scan(&c.ID, &c.Nombre)
		if err != nil {
			return nil, fmt.Errorf("Error inesperado, detalle: %v", err)
		}

		lista = append(lista, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error inesperado, detalle: %v", err)
	}

	return lista, nil
}
