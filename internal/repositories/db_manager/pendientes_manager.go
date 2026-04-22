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
)

type pendientesRepository struct {
	db *sql.DB
}

func NewPendientesRepository(db *sql.DB) *pendientesRepository {
	return &pendientesRepository{db: db}
}

func (r *pendientesRepository) Crear(p *models.Pendientes) (int, error) {

	result, err := r.db.Exec(
		`INSERT INTO pendientes (
			titulo, descripcion, solicitante_id, fecha_pedido
		) VALUES (?, ?, ?, ?)`,
		p.Titulo,
		p.Descripcion,
		p.SolicitanteID,
		p.FechaPedido,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {

			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintForeignKey {
				return 0, appErrors.SolicitanteNoEncontrado
			}

		}
		return 0, fmt.Errorf("error inesperado: %w", err)

	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Error inesperado, detalle: %v", err)
	}

	return int(id), nil
}

func (r *pendientesRepository) Asignar(id int, asignadoID int, fecha time.Time) error {

	result, err := r.db.Exec(
		`UPDATE pendientes
		 SET asignado_id = ?, fecha_asignado = ?
		 WHERE id = ?`,
		asignadoID,
		fecha,
		id,
	)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {

			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintForeignKey {
				return appErrors.ColaboradorNoEncontrado
			}

		}
		return fmt.Errorf("error inesperado: %w", err)

	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error inesperado: %w", err)
	}
	if rows == 0 {
		return appErrors.PendienteNoEncontrado
	}

	return nil
}

func (r *pendientesRepository) Terminar(id int, cierre *string, fecha time.Time) error {

	result, err := r.db.Exec(
		`UPDATE pendientes
		 SET cierre = ?, fecha_cierre = ?, finalizado = 1
		 WHERE id = ?`,
		cierre,
		fecha,
		id,
	)
	if err != nil {
		return fmt.Errorf("error inesperado: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error inesperado: %w", err)
	}
	if rows == 0 {
		return appErrors.PendienteNoEncontrado
	}

	return nil
}

func (r *pendientesRepository) ModificarDescripcion(id int, descripcion string) error {

	result, err := r.db.Exec(
		`UPDATE pendientes
		 SET descripcion = ?
		 WHERE id = ?`,
		descripcion,
		id,
	)
	if err != nil {
		return fmt.Errorf("error inesperado: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error inesperado: %w", err)
	}
	if rows == 0 {
		return appErrors.PendienteNoEncontrado
	}

	return nil
}

func (r *pendientesRepository) ModificarIdentificacionTablaPendiente(id int, identificacion string) error {

	result, err := r.db.Exec(
		`UPDATE pendientes
		 SET identificacion_tabla_pendiente = ?
		 WHERE id = ?`,
		identificacion,
		id,
	)
	if err != nil {
		return fmt.Errorf("error inesperado: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error inesperado: %w", err)
	}
	if rows == 0 {
		return appErrors.PendienteNoEncontrado
	}

	return nil
}

func (r *pendientesRepository) BuscarPorTituloDescripcion(texto string, finalizado bool) ([]models.Pendientes, error) {

	if strings.TrimSpace(texto) == "" {
		return nil, appErrors.ParametroDeBusquedaVacio
	}

	rows, err := r.db.Query(
		`SELECT id, titulo, descripcion,
		        solicitante_id, fecha_pedido,
		        asignado_id, fecha_asignado,
		        cierre, fecha_cierre, finalizado,
		        identificacion_tabla_pendiente
		 FROM pendientes
		 WHERE (titulo LIKE ? OR descripcion LIKE ?) AND finalizado = ?
		 ORDER BY fecha_pedido`,
		"%"+texto+"%",
		"%"+texto+"%",
		finalizado,
	)
	if err != nil {
		return nil, fmt.Errorf("error inesperado: %v", err)
	}
	defer rows.Close()

	var lista []models.Pendientes

	for rows.Next() {
		var p models.Pendientes

		err := rows.Scan(
			&p.ID,
			&p.Titulo,
			&p.Descripcion,
			&p.SolicitanteID,
			&p.FechaPedido,
			&p.AsignadoID,
			&p.FechaAsignado,
			&p.Cierre,
			&p.FechaCierre,
			&p.Finalizado,
			&p.IdentificacionTablaPendiente,
		)
		if err != nil {
			return nil, fmt.Errorf("error inesperado: %v", err)
		}

		lista = append(lista, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando resultados: %v", err)
	}

	return lista, nil
}

func (r *pendientesRepository) ListarPorAsignado(asignadoID int, finalizado bool) ([]models.Pendientes, error) {

	rows, err := r.db.Query(
		`SELECT id, titulo, descripcion,
		        solicitante_id, fecha_pedido,
		        asignado_id, fecha_asignado,
		        cierre, fecha_cierre, finalizado,
		        identificacion_tabla_pendiente
		 FROM pendientes
		 WHERE asignado_id = ? AND finalizado = ?
		 ORDER BY fecha_pedido`,
		asignadoID,
		finalizado,
	)
	if err != nil {
		return nil, fmt.Errorf("error inesperado: %v", err)
	}
	defer rows.Close()

	var lista []models.Pendientes

	for rows.Next() {
		var p models.Pendientes

		err := rows.Scan(
			&p.ID,
			&p.Titulo,
			&p.Descripcion,
			&p.SolicitanteID,
			&p.FechaPedido,
			&p.AsignadoID,
			&p.FechaAsignado,
			&p.Cierre,
			&p.FechaCierre,
			&p.Finalizado,
			&p.IdentificacionTablaPendiente,
		)
		if err != nil {
			return nil, fmt.Errorf("error inesperado: %v", err)
		}

		lista = append(lista, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando resultados: %v", err)
	}

	return lista, nil
}

func (r *pendientesRepository) Listar(finalizado bool, limit, offset int) ([]models.Pendientes, error) {

	rows, err := r.db.Query(
		`SELECT id, titulo, descripcion,
		        solicitante_id, fecha_pedido,
		        asignado_id, fecha_asignado,
		        cierre, fecha_cierre, finalizado,
		        identificacion_tabla_pendiente
		 FROM pendientes
		 WHERE finalizado = ?
		 ORDER BY fecha_pedido LIMIT ? OFFSET ?`,
		finalizado, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error inesperado: %v", err)
	}
	defer rows.Close()

	var lista []models.Pendientes

	for rows.Next() {
		var p models.Pendientes

		err := rows.Scan(
			&p.ID,
			&p.Titulo,
			&p.Descripcion,
			&p.SolicitanteID,
			&p.FechaPedido,
			&p.AsignadoID,
			&p.FechaAsignado,
			&p.Cierre,
			&p.FechaCierre,
			&p.Finalizado,
			&p.IdentificacionTablaPendiente,
		)
		if err != nil {
			return nil, fmt.Errorf("error inesperado: %v", err)
		}

		lista = append(lista, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando resultados: %v", err)
	}

	return lista, nil
}
