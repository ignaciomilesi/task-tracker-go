package dbmanager

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type codigoIDRepository struct {
	db *sql.DB
}

func NewCodigoIDRepository(db *sql.DB) *codigoIDRepository {
	return &codigoIDRepository{db: db}
}

// Cargar nuevo código ID
func (r *codigoIDRepository) Cargar(nuevoCodigoID *models.CodigoID) error {

	if strings.TrimSpace(nuevoCodigoID.Codigo) == "" {
		return appErrors.CodigoIDVacio
	}

	_, err := r.db.Exec(
		`INSERT INTO codigo_ID (codigo, descripcion, estado, fecha_pedido, fecha_actualizacion)
		 VALUES (?, ?, ?, ?, ?)`,
		nuevoCodigoID.Codigo,
		nuevoCodigoID.Descripcion,
		nuevoCodigoID.Estado,
		nuevoCodigoID.FechaPedido,
		nuevoCodigoID.FechaActualizacion,
	)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return appErrors.CodigoIDDuplicado
			}
		}
		return fmt.Errorf("error inesperado: %v", err)
	}

	return nil
}

// Obtener por código
func (r *codigoIDRepository) ObtenerDetalle(codigo string) (*models.CodigoID, error) {

	if strings.TrimSpace(codigo) == "" {
		return nil, appErrors.CodigoIDVacio
	}

	var c models.CodigoID

	err := r.db.QueryRow(
		`SELECT codigo, descripcion, estado, fecha_pedido, fecha_actualizacion
		 FROM codigo_ID
		 WHERE codigo = ?`,
		codigo,
	).Scan(
		&c.Codigo,
		&c.Descripcion,
		&c.Estado,
		&c.FechaPedido,
		&c.FechaActualizacion,
	)

	if err == sql.ErrNoRows {
		return nil, appErrors.CodigoIDNoEncontrado
	}
	if err != nil {
		return nil, fmt.Errorf("error inesperado: %v", err)
	}

	return &c, nil
}

// Filtrar por estado
func (r *codigoIDRepository) FiltrarPorEstado(estado string) ([]models.CodigoID, error) {

	rows, err := r.db.Query(
		`SELECT codigo, descripcion, estado, fecha_pedido, fecha_actualizacion
		 FROM codigo_ID
		 WHERE estado = ?`,
		estado,
	)
	if err != nil {
		return nil, fmt.Errorf("error inesperado: %v", err)
	}
	defer rows.Close()

	var lista []models.CodigoID

	for rows.Next() {
		var c models.CodigoID

		err := rows.Scan(
			&c.Codigo,
			&c.Descripcion,
			&c.Estado,
			&c.FechaPedido,
			&c.FechaActualizacion,
		)
		if err != nil {
			return nil, fmt.Errorf("error inesperado: %v", err)
		}

		lista = append(lista, c)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando resultados: %v", err)
	}

	return lista, nil
}

// Actualizar estado y fecha de actualización
func (r *codigoIDRepository) ActualizarEstado(codigo string, nuevoEstado string, fecha time.Time) error {

	if strings.TrimSpace(codigo) == "" {
		return appErrors.CodigoIDVacio
	}

	result, err := r.db.Exec(
		`UPDATE codigo_ID
		 SET estado = ?, fecha_actualizacion = ?
		 WHERE codigo = ?`,
		nuevoEstado,
		fecha,
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
		return appErrors.CodigoIDNoEncontrado
	}

	return nil
}

func (r *codigoIDRepository) Listar(limit, offset int) ([]models.CodigoID, error) {

	rows, err := r.db.Query(
		`SELECT codigo, descripcion, estado, fecha_pedido, fecha_actualizacion
		 FROM codigo_ID
        LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("Error inesperado, detalle: %v", err)
	}
	defer rows.Close()

	var lista []models.CodigoID

	for rows.Next() {
		var c models.CodigoID

		err := rows.Scan(&c.Codigo, &c.Descripcion, &c.Estado, &c.FechaPedido, &c.FechaActualizacion)
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
