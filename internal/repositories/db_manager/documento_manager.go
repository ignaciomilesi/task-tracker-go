package dbmanager

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"

	"github.com/mattn/go-sqlite3"
)

type documentoRepository struct {
	db *sql.DB
}

func NewDocumentoRepository(db *sql.DB) *documentoRepository {
	return &documentoRepository{db: db}
}

// Cargar nuevo documento
func (r *documentoRepository) Cargar(doc *models.Documento) error {

	if doc == nil {
		return appErrors.ParametroDeCargaVacio
	}

	if strings.TrimSpace(doc.Codigo) == "" {
		return appErrors.CodigoDocumentoVacio
	}

	_, err := r.db.Exec(
		`INSERT INTO documento (codigo, emision, titulo, tipo, ubicacion_path, backup_path)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		doc.Codigo,
		doc.Emision,
		doc.Titulo,
		doc.Tipo,
		doc.UbicacionPath,
		doc.BackupPath,
	)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return appErrors.DocumentoDuplicado
			}
		}
		return fmt.Errorf("error inesperado: %v", err)
	}

	return nil
}

// Obtener detalle por código
func (r *documentoRepository) ObtenerDetalle(codigo string) (*models.Documento, error) {

	if strings.TrimSpace(codigo) == "" {
		return nil, appErrors.ParametroDeBusquedaVacio
	}

	var d models.Documento

	err := r.db.QueryRow(
		`SELECT codigo, emision, titulo, tipo, ubicacion_path, backup_path
		 FROM documento
		 WHERE codigo = ?`,
		codigo,
	).Scan(
		&d.Codigo,
		&d.Emision,
		&d.Titulo,
		&d.Tipo,
		&d.UbicacionPath,
		&d.BackupPath,
	)

	if err == sql.ErrNoRows {
		return nil, appErrors.DocumentoNoEncontrado
	}
	if err != nil {
		return nil, fmt.Errorf("error inesperado: %v", err)
	}

	return &d, nil
}

// Filtrar por tipo
func (r *documentoRepository) FiltrarPorTipo(tipo string) ([]models.Documento, error) {

	if strings.TrimSpace(tipo) == "" {
		return nil, appErrors.ParametroDeBusquedaVacio
	}

	rows, err := r.db.Query(
		`SELECT codigo, emision, titulo, tipo, ubicacion_path, backup_path
		 FROM documento
		 WHERE tipo = ?`,
		tipo,
	)
	if err != nil {
		return nil, fmt.Errorf("error inesperado: %v", err)
	}
	defer rows.Close()

	var lista []models.Documento

	for rows.Next() {
		var d models.Documento

		err := rows.Scan(
			&d.Codigo,
			&d.Emision,
			&d.Titulo,
			&d.Tipo,
			&d.UbicacionPath,
			&d.BackupPath,
		)
		if err != nil {
			return nil, fmt.Errorf("error inesperado: %v", err)
		}

		lista = append(lista, d)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando resultados: %v", err)
	}

	return lista, nil
}

// Filtrar por título (LIKE)
func (r *documentoRepository) FiltrarPorTitulo(titulo string) ([]models.Documento, error) {

	if strings.TrimSpace(titulo) == "" {
		return nil, appErrors.ParametroDeBusquedaVacio
	}

	rows, err := r.db.Query(
		`SELECT codigo, emision, titulo, tipo, ubicacion_path, backup_path
		 FROM documento
		 WHERE titulo LIKE ?`,
		"%"+titulo+"%",
	)
	if err != nil {
		return nil, fmt.Errorf("error inesperado: %v", err)
	}
	defer rows.Close()

	var lista []models.Documento

	for rows.Next() {
		var d models.Documento

		err := rows.Scan(
			&d.Codigo,
			&d.Emision,
			&d.Titulo,
			&d.Tipo,
			&d.UbicacionPath,
			&d.BackupPath,
		)
		if err != nil {
			return nil, fmt.Errorf("error inesperado: %v", err)
		}

		lista = append(lista, d)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando resultados: %v", err)
	}

	return lista, nil
}

// Actualizar path
func (r *documentoRepository) ActualizarPath(codigo string, nuevoPath string) error {

	if strings.TrimSpace(codigo) == "" {
		return appErrors.ParametroDeBusquedaVacio
	}
	if strings.TrimSpace(nuevoPath) == "" {
		return appErrors.ParametroDeCargaVacio
	}

	result, err := r.db.Exec(
		`UPDATE documento
		 SET ubicacion_path = ?
		 WHERE codigo = ?`,
		nuevoPath,
		codigo,
	)
	if err != nil {
		return fmt.Errorf("error inesperado: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error obteniendo filas afectadas: %v", err)
	}

	if rowsAffected == 0 {
		return appErrors.DocumentoNoEncontrado
	}

	return nil
}

// Actualizar backup path
func (r *documentoRepository) ActualizarBackupPath(codigo string, nuevoPath string) error {

	if strings.TrimSpace(codigo) == "" {
		return appErrors.ParametroDeBusquedaVacio
	}
	if strings.TrimSpace(nuevoPath) == "" {
		return appErrors.ParametroDeCargaVacio
	}

	result, err := r.db.Exec(
		`UPDATE documento
		 SET backup_path = ?
		 WHERE codigo = ?`,
		nuevoPath,
		codigo,
	)
	if err != nil {
		return fmt.Errorf("error inesperado: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error obteniendo filas afectadas: %v", err)
	}

	if rowsAffected == 0 {
		return appErrors.DocumentoNoEncontrado
	}

	return nil
}
