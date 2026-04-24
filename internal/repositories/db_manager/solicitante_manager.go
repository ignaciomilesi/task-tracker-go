package dbmanager

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type solicitanteRepository struct {
	db *sql.DB
}

func NewsolicitanteRepository(db *sql.DB) *solicitanteRepository {
	return &solicitanteRepository{db: db}
}

// Crear nuevo solicitante, devuelve el id creado
func (r *solicitanteRepository) Crear(ctx context.Context, nuevoSolicitante *models.Solicitante) (int, error) {

	if nuevoSolicitante == nil {
		return 0, appErrors.ParametroDeCargaVacio
	}

	res, err := r.db.ExecContext(
		ctx,
		"INSERT INTO solicitante (nombre) VALUES (?)",
		nuevoSolicitante.Nombre,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return 0, appErrors.SolicitanteDuplicado
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

func (r *solicitanteRepository) ObtenerIDPorNombre(ctx context.Context, nombre string) (int, error) {
	var id int

	err := r.db.QueryRowContext( // toma solo el primer resultado
		ctx,
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

// Busca en la lista. El parámetro no puede estar vacío
func (r *solicitanteRepository) Buscar(ctx context.Context, parametro string) ([]models.Solicitante, error) {

	if parametro == "" {
		return nil, appErrors.ParametroDeBusquedaVacio
	}
	rows, err := r.db.QueryContext(ctx, "SELECT id, nombre FROM solicitante WHERE nombre LIKE ?", "%"+parametro+"%")
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

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error inesperado, detalle: %v", err)
	}

	return lista, nil
}

func (r *solicitanteRepository) Listar(ctx context.Context, limit, offset int) ([]models.Solicitante, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, nombre FROM solicitante
		ORDER BY id LIMIT ? OFFSET ?`, limit, offset)
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

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error inesperado, detalle: %v", err)
	}

	return lista, nil
}
