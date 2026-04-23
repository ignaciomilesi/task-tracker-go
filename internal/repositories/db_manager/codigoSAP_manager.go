package dbmanager

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type codigoSAPRepository struct {
	db *sql.DB
}

func NewCodigoSAPRepository(db *sql.DB) *codigoSAPRepository {
	return &codigoSAPRepository{db: db}
}

// Cargar nuevo código SAP
func (r *codigoSAPRepository) Cargar(NuevoCodigoSap *models.CodigoSAP) error {

	if NuevoCodigoSap == nil {
		return appErrors.ParametroDeCargaVacio
	}

	if strings.TrimSpace(NuevoCodigoSap.Codigo) == "" {
		return appErrors.CodigoSAPVacio
	}

	_, err := r.db.Exec(
		"INSERT INTO codigo_SAP (codigo, descripcion) VALUES (?, ?)",
		NuevoCodigoSap.Codigo,
		NuevoCodigoSap.Descripcion,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {

			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return appErrors.CodigoSAPDuplicado
			}
		}
		return fmt.Errorf("error inesperado: %v", err)
	}

	return nil
}

// Obtener por código
func (r *codigoSAPRepository) ObtenerDetalle(codigo string) (*models.CodigoSAP, error) {

	var c models.CodigoSAP

	err := r.db.QueryRow(
		"SELECT codigo, descripcion FROM codigo_SAP WHERE codigo = ?",
		codigo,
	).Scan(&c.Codigo, &c.Descripcion)

	if err == sql.ErrNoRows {
		return nil, appErrors.CodigoSAPNoEncontrado
	}
	if err != nil {
		return nil, fmt.Errorf("Error inesperado, detalle: %v", err)
	}

	return &c, nil
}

// Buscar por descripción
func (r *codigoSAPRepository) BuscarPorDescripcion(parametro string) ([]models.CodigoSAP, error) {

	if parametro == "" {
		return nil, appErrors.ParametroDeBusquedaVacio
	}
	rows, err := r.db.Query(
		"SELECT codigo, descripcion FROM codigo_SAP WHERE descripcion LIKE ?",
		"%"+parametro+"%",
	)
	if err != nil {
		return nil, fmt.Errorf("Error inesperado, detalle: %v", err)
	}
	defer rows.Close()

	var lista []models.CodigoSAP

	for rows.Next() {
		var c models.CodigoSAP

		err := rows.Scan(&c.Codigo, &c.Descripcion)
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

// Modificar descripción de un código SAP
func (r *codigoSAPRepository) ModificarDescripcion(codigo string, nuevaDescripcion string) error {

	if strings.TrimSpace(codigo) == "" {
		return appErrors.CodigoSAPVacio
	}

	result, err := r.db.Exec(
		"UPDATE codigo_SAP SET descripcion = ? WHERE codigo = ?", nuevaDescripcion, codigo)
	if err != nil {
		return fmt.Errorf("error inesperado: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error inesperado: %v", err)
	}

	if rowsAffected == 0 {
		return appErrors.CodigoSAPNoEncontrado
	}

	return nil
}

func (r *codigoSAPRepository) Listar(limit, offset int) ([]models.CodigoSAP, error) {

	rows, err := r.db.Query(
		`SELECT codigo, descripcion FROM codigo_SAP ORDER BY codigo
        LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("Error inesperado, detalle: %v", err)
	}
	defer rows.Close()

	var lista []models.CodigoSAP

	for rows.Next() {
		var c models.CodigoSAP

		err := rows.Scan(&c.Codigo, &c.Descripcion)
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
