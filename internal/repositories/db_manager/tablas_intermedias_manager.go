package dbmanager

import (
	"database/sql"
	"errors"
	"fmt"

	"task-tracker-go/internal/appErrors"

	"github.com/mattn/go-sqlite3"
)

type tablaIntermediaRepository struct {
	db *sql.DB
}

func NewTablaIntermediaRepository(db *sql.DB) *tablaIntermediaRepository {
	return &tablaIntermediaRepository{db: db}
}

func (r *tablaIntermediaRepository) VincularPendienteDocumento(pendienteID, documentoID int) error {

	_, err := r.db.Exec(
		`INSERT INTO ti_pendientes_documento (pendiente_id, documento_id)
		 VALUES (?, ?)`,
		pendienteID,
		documentoID,
	)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {

			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintForeignKey {
				return appErrors.FkNoEncontrado
			}

			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return appErrors.RelacionDuplicada
			}
		}

		return fmt.Errorf("error inesperado: %w", err)
	}

	return nil
}

func (r *tablaIntermediaRepository) VincularPendienteCodigoSAP(pendienteID int, codigoSAP string) error {

	_, err := r.db.Exec(
		`INSERT INTO ti_pendientes_codigo_sap (pendiente_id, codigo_sap_codigo)
		 VALUES (?, ?)`,
		pendienteID,
		codigoSAP,
	)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {

			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintForeignKey {
				return appErrors.FkNoEncontrado
			}

			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return appErrors.RelacionDuplicada
			}
		}

		return fmt.Errorf("error inesperado: %w", err)
	}

	return nil
}

func (r *tablaIntermediaRepository) VincularPendienteCodigoID(pendienteID int, codigoID string) error {

	_, err := r.db.Exec(
		`INSERT INTO ti_pendientes_codigo_id (pendiente_id, codigo_id_codigo)
		 VALUES (?, ?)`,
		pendienteID,
		codigoID,
	)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {

			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintForeignKey {
				return appErrors.FkNoEncontrado
			}

			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return appErrors.RelacionDuplicada
			}
		}

		return fmt.Errorf("error inesperado: %w", err)
	}

	return nil
}

func (r *tablaIntermediaRepository) ListarDocumentosPorPendiente(pendienteID int) ([]int, error) {

	rows, err := r.db.Query(
		`SELECT documento_id
		 FROM ti_pendientes_documento
		 WHERE pendiente_id = ?`,
		pendienteID,
	)
	if err != nil {
		return nil, fmt.Errorf("error inesperado: %w", err)
	}
	defer rows.Close()

	var lista []int

	for rows.Next() {
		var id int

		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("error inesperado: %w", err)
		}

		lista = append(lista, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando resultados: %w", err)
	}

	return lista, nil
}

func (r *tablaIntermediaRepository) ListarCodigosSAPPorPendiente(pendienteID int) ([]string, error) {

	rows, err := r.db.Query(
		`SELECT codigo_sap_codigo
		 FROM ti_pendientes_codigo_sap
		 WHERE pendiente_id = ?`,
		pendienteID,
	)
	if err != nil {
		return nil, fmt.Errorf("error inesperado: %w", err)
	}
	defer rows.Close()

	var lista []string

	for rows.Next() {
		var codigo string

		if err := rows.Scan(&codigo); err != nil {
			return nil, fmt.Errorf("error inesperado: %w", err)
		}

		lista = append(lista, codigo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando resultados: %w", err)
	}

	return lista, nil
}

func (r *tablaIntermediaRepository) ListarCodigosIDPorPendiente(pendienteID int) ([]string, error) {

	rows, err := r.db.Query(
		`SELECT codigo_id_codigo
		 FROM ti_pendientes_codigo_id
		 WHERE pendiente_id = ?`,
		pendienteID,
	)
	if err != nil {
		return nil, fmt.Errorf("error inesperado: %w", err)
	}
	defer rows.Close()

	var lista []string

	for rows.Next() {
		var codigo string

		if err := rows.Scan(&codigo); err != nil {
			return nil, fmt.Errorf("error inesperado: %w", err)
		}

		lista = append(lista, codigo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando resultados: %w", err)
	}

	return lista, nil
}
