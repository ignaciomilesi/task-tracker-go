package dbmanager

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"

	"github.com/mattn/go-sqlite3"
)

type avanceRepository struct {
	db *sql.DB
}

func NewAvanceRepository(db *sql.DB) *avanceRepository {
	return &avanceRepository{db: db}
}

// Cargar nuevo avance
func (r *avanceRepository) Cargar(ctx context.Context, nuevoAvance *models.Avance) (int, error) {

	if nuevoAvance == nil {
		return 0, appErrors.ParametroDeCargaVacio
	}

	res, err := r.db.ExecContext(ctx,
		`INSERT INTO avance (pendiente_id, descripcion, fecha, mail_path)
		 VALUES (?, ?, ?, ?)`,
		nuevoAvance.PendienteID,
		nuevoAvance.Descripcion,
		nuevoAvance.Fecha,
		nuevoAvance.MailPath,
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
		return 0, fmt.Errorf("Error inesperado, detalle: %v", err)
	}

	return int(id), nil
}

// Obtener detalle por ID
func (r *avanceRepository) ObtenerDetalle(ctx context.Context, id int) (*models.Avance, error) {

	var a models.Avance

	err := r.db.QueryRowContext(ctx,
		`SELECT id, pendiente_id, descripcion, fecha, mail_path
		 FROM avance
		 WHERE id = ?`,
		id,
	).Scan(
		&a.ID,
		&a.PendienteID,
		&a.Descripcion,
		&a.Fecha,
		&a.MailPath,
	)

	if err == sql.ErrNoRows {
		return nil, appErrors.AvanceNoEncontrado
	}

	if err != nil {
		return nil, fmt.Errorf("error inesperado: %v", err)
	}

	return &a, nil
}

// Eliminar avance
func (r *avanceRepository) Eliminar(ctx context.Context, id int) error {

	result, err := r.db.ExecContext(ctx,
		`DELETE FROM avance WHERE id = ?`,
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
		return appErrors.AvanceNoEncontrado
	}

	return nil
}

// Filtrar avances por PendienteID
func (r *avanceRepository) FiltrarPorPendiente(ctx context.Context, pendienteID int) ([]models.Avance, error) {

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, pendiente_id, descripcion, fecha, mail_path
		 FROM avance
		 WHERE pendiente_id = ?
		 ORDER BY fecha DESC`,
		pendienteID,
	)
	if err != nil {
		return nil, fmt.Errorf("error inesperado: %v", err)
	}
	defer rows.Close()

	var lista []models.Avance

	for rows.Next() {
		var a models.Avance

		err := rows.Scan(
			&a.ID,
			&a.PendienteID,
			&a.Descripcion,
			&a.Fecha,
			&a.MailPath,
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
