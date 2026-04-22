package dbmanager

import (
	"database/sql"
	"errors"
	"fmt"

	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"

	"github.com/mattn/go-sqlite3"
)

type adjuntoRepository struct {
	db *sql.DB
}

func NewAdjuntoRepository(db *sql.DB) *adjuntoRepository {
	return &adjuntoRepository{db: db}
}

// Cargar nuevo adjunto
func (r *adjuntoRepository) Cargar(nuevoAdjunto *models.Adjunto) (int, error) {

	if nuevoAdjunto == nil {
		return 0, appErrors.ParametroDeCargaVacio
	}

	res, err := r.db.Exec(
		`INSERT INTO adjunto (pendiente_id, descripcion, archivo_path)
		 VALUES (?, ?, ?)`,
		nuevoAdjunto.PendienteID,
		nuevoAdjunto.Descripcion,
		nuevoAdjunto.ArchivoPath,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintForeignKey {
				return 0, appErrors.PendienteNoEncontrado
			}
		}
		return 0, fmt.Errorf("error inesperado: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error obteniendo ID: %v", err)
	}

	return int(id), nil
}

// Obtener detalle por ID
func (r *adjuntoRepository) ObtenerDetalle(id int) (*models.Adjunto, error) {

	var a models.Adjunto

	err := r.db.QueryRow(
		`SELECT id, pendiente_id, descripcion, archivo_path
		 FROM adjunto
		 WHERE id = ?`,
		id,
	).Scan(
		&a.ID,
		&a.PendienteID,
		&a.Descripcion,
		&a.ArchivoPath,
	)

	if err == sql.ErrNoRows {
		return nil, appErrors.AdjuntoNoEncontrado
	}

	if err != nil {
		return nil, fmt.Errorf("error inesperado: %v", err)
	}

	return &a, nil
}

// Eliminar adjunto
func (r *adjuntoRepository) Eliminar(id int) error {

	result, err := r.db.Exec(
		`DELETE FROM adjunto WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("error inesperado: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error obteniendo filas afectadas: %v", err)
	}

	if rowsAffected == 0 {
		return appErrors.AdjuntoNoEncontrado
	}

	return nil
}

// Filtrar por PendienteID
func (r *adjuntoRepository) FiltrarPorPendiente(pendienteID int) ([]models.Adjunto, error) {

	rows, err := r.db.Query(
		`SELECT id, pendiente_id, descripcion, archivo_path
		 FROM adjunto
		 WHERE pendiente_id = ?`,
		pendienteID,
	)
	if err != nil {
		return nil, fmt.Errorf("error inesperado: %v", err)
	}
	defer rows.Close()

	var lista []models.Adjunto

	for rows.Next() {
		var a models.Adjunto

		err := rows.Scan(
			&a.ID,
			&a.PendienteID,
			&a.Descripcion,
			&a.ArchivoPath,
		)
		if err != nil {
			return nil, fmt.Errorf("error inesperado: %v", err)
		}

		lista = append(lista, a)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando resultados: %v", err)
	}

	return lista, nil
}

// Actualizar Path
func (r *adjuntoRepository) ActualizarPath(id int, path string) error {

	result, err := r.db.Exec(
		`UPDATE adjunto SET archivo_path = ? WHERE id = ?`,
		path,
		id,
	)
	if err != nil {
		return fmt.Errorf("error inesperado: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error obteniendo filas afectadas: %v", err)
	}

	if rowsAffected == 0 {
		return appErrors.AdjuntoNoEncontrado
	}

	return nil
}

// Actualizar Descripción
func (r *adjuntoRepository) ActualizarDescripcion(id int, descripcion string) error {

	result, err := r.db.Exec(
		`UPDATE adjunto SET descripcion = ? WHERE id = ?`,
		descripcion,
		id,
	)
	if err != nil {
		return fmt.Errorf("error inesperado: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error obteniendo filas afectadas: %v", err)
	}

	if rowsAffected == 0 {
		return appErrors.AdjuntoNoEncontrado
	}

	return nil
}
