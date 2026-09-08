package repo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"task-tracker-go/internal/pendiente/dto"
	"task-tracker-go/internal/pendiente/model"

	"github.com/mattn/go-sqlite3"
)

type pendienteDB struct {
	ID            int       `db:"id"`
	Titulo        string    `db:"titulo"`
	Descripcion   string    `db:"descripcion"`
	SolicitanteID int       `db:"solicitante_id"`
	FechaPedido   time.Time `db:"fecha_pedido"`

	AsignadoID    *int       `db:"asignado_id"` //colaborador
	FechaAsignado *time.Time `db:"fecha_asignado"`

	Finalizado  bool       `db:"finalizado"`
	Cierre      *string    `db:"cierre"`
	FechaCierre *time.Time `db:"fecha_cierre"`
}

func (p pendienteDB) generarDesdeDomain(av model.Pendiente) pendienteDB {
	p.ID = av.ID
	p.Titulo = av.Titulo
	p.Descripcion = av.Descripcion
	p.SolicitanteID = av.SolicitanteID
	p.FechaPedido = av.FechaPedido
	p.AsignadoID = av.AsignadoID
	p.FechaAsignado = av.FechaAsignado
	p.Finalizado = av.Finalizado
	p.Cierre = av.Cierre
	p.FechaCierre = av.FechaCierre
	return p
}

func (p pendienteDB) devolverDomain() model.Pendiente {
	return model.Pendiente{
		ID:            p.ID,
		Titulo:        p.Titulo,
		Descripcion:   p.Descripcion,
		SolicitanteID: p.SolicitanteID,
		FechaPedido:   p.FechaPedido,
		AsignadoID:    p.AsignadoID,
		FechaAsignado: p.FechaAsignado,
		Finalizado:    p.Finalizado,
		Cierre:        p.Cierre,
		FechaCierre:   p.FechaCierre,
	}

}

type pendienteDetalleNombreDB struct {
	ID                int       `db:"id"`
	Titulo            string    `db:"titulo"`
	Descripcion       string    `db:"descripcion"`
	SolicitanteNombre string    `db:"solicitante_nombre"`
	FechaPedido       time.Time `db:"fecha_pedido"`

	AsignadoNombre *string    `db:"asignado_nombre"`
	FechaAsignado  *time.Time `db:"fecha_asignado"`

	Finalizado  bool       `db:"finalizado"`
	Cierre      *string    `db:"cierre"`
	FechaCierre *time.Time `db:"fecha_cierre"`
}

func (r *Repo) Crear(ctx context.Context, nuevoPendiente *model.Pendiente) (int, error) {

	if nuevoPendiente == nil {
		return -1, model.ErrModelDeCargaNil
	}

	var nuevaSolicitudPendiente pendienteDB
	nuevaSolicitudPendiente.generarDesdeDomain(*nuevoPendiente)

	query := `INSERT INTO pendientes (
			titulo, descripcion, solicitante_id, fecha_pedido
		) VALUES (:titulo, :descripcion, :solicitante_id, :fecha_pedido)`

	id, err := r.db.Cargar(ctx, &nuevaSolicitudPendiente, query)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {

			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintForeignKey {
				return -1, model.ErrUsuarioNoEncontrado
			}

		}
		return -1, err

	}

	return id, nil
}

func (r *Repo) ObtenerDetalle(ctx context.Context, id int) (*dto.PendienteDetalleNombre, error) {

	if id < 0 {
		return nil, model.ErrModelDeCargaNil
	}

	var detalle pendienteDetalleNombreDB

	query := `SELECT 
				p.id,
				p.titulo,
				p.descripcion,
				us.nombre AS solicitante_nombre,
				p.fecha_pedido,
				ua.nombre AS asignado_nombre,
				p.fecha_asignado,
				p.cierre,
				p.fecha_cierre,
				p.finalizado
			FROM pendientes p
			LEFT JOIN usuarios us 
				ON p.solicitante_id = us.id
			LEFT JOIN usuarios ua 
				ON p.asignado_id = ua.id
			WHERE p.id = ?`

	if err := r.db.ObtenerDetalle(ctx, &detalle, query, id); err != nil {

		return nil, err
	}

	return &dto.PendienteDetalleNombre{
		ID:                detalle.ID,
		Titulo:            detalle.Titulo,
		Descripcion:       detalle.Descripcion,
		SolicitanteNombre: detalle.SolicitanteNombre,
		FechaPedido:       detalle.FechaPedido,
		AsignadoNombre:    detalle.AsignadoNombre,
		FechaAsignado:     detalle.FechaAsignado,
		Finalizado:        detalle.Finalizado,
		Cierre:            detalle.Cierre,
		FechaCierre:       detalle.FechaCierre,
	}, nil
}

func (r *Repo) Actualizar(ctx context.Context, pendiente *model.Pendiente) error {
	if pendiente == nil {
		return model.ErrModelDeCargaNil
	}

	var solicitud pendienteDB
	solicitud.generarDesdeDomain(*pendiente)

	query := `UPDATE pendientes
		 SET titulo = :titulo, 
		 descripcion = :descripcion, 
		 solicitante_id = :solicitante_id, 
		 fecha_pedido = :fecha_pedido,
		 asignado_id = :asignado_id,
		 fecha_asignado = :fecha_asignado,
		 finalizado = :finalizado,
		 cierre = :cierre,
		 fecha_cierre = :fecha_cierre
		 WHERE id = :id`

	err := r.db.ModificarRegistro(ctx, query, solicitud)

	return err
}

func (r *Repo) BuscarPorTituloDescripcion(ctx context.Context, texto string, finalizado bool) ([]dto.PendienteDetalleNombre, error) {

	if strings.TrimSpace(texto) == "" {
		return nil, model.ErrParametrosNoValidos
	}

	var lista []pendienteDetalleNombreDB

	query := `SELECT 
				p.id,
				p.titulo,
				p.descripcion,
				us.nombre AS solicitante_nombre,
				p.fecha_pedido,
				ua.nombre AS asignado_nombre,
				p.fecha_asignado,
				p.cierre,
				p.fecha_cierre,
				p.finalizado
			FROM pendientes p
			LEFT JOIN usuarios us 
				ON p.solicitante_id = us.id
			LEFT JOIN usuarios ua 
				ON p.asignado_id = ua.id
			WHERE 
				(p.titulo LIKE ? OR p.descripcion LIKE ?)
				AND p.finalizado = ?
			ORDER BY p.fecha_pedido;`

	err := r.db.ObtenerLista(ctx, &lista, query, "%"+texto+"%", "%"+texto+"%", finalizado)
	if err != nil {
		return nil, fmt.Errorf("Error inesperado, detalle: %w", err)
	}

	var listaDto []dto.PendienteDetalleNombre
	for _, item := range lista {
		dtoItem := dto.PendienteDetalleNombre{
			ID:                item.ID,
			Titulo:            item.Titulo,
			Descripcion:       item.Descripcion,
			SolicitanteNombre: item.SolicitanteNombre,
			FechaPedido:       item.FechaPedido,
			AsignadoNombre:    item.AsignadoNombre,
			FechaAsignado:     item.FechaAsignado,
			Finalizado:        item.Finalizado,
			Cierre:            item.Cierre,
			FechaCierre:       item.FechaCierre,
		}
		listaDto = append(listaDto, dtoItem)
	}

	return listaDto, nil
}

func (r *Repo) Listar(ctx context.Context, finalizado bool, asignadoID *int, limit, offset int) ([]dto.PendienteDetalleNombre, error) {

	var lista []pendienteDetalleNombreDB

	query := `
		SELECT 
			p.id,
			p.titulo,
			p.descripcion,
			us.nombre AS solicitante_nombre,
			p.fecha_pedido,
			ua.nombre AS asignado_nombre,
			p.fecha_asignado,
			p.cierre,
			p.fecha_cierre,
			p.finalizado
		FROM pendientes p
		LEFT JOIN usuarios us 
			ON p.solicitante_id = us.id
		LEFT JOIN usuarios ua 
			ON p.asignado_id = ua.id
		WHERE p.finalizado = ?
	`

	args := []any{finalizado}

	if asignadoID != nil {
		query += ` AND p.asignado_id = ?`
		args = append(args, *asignadoID)
	}

	query += `
		ORDER BY p.fecha_pedido
		LIMIT ? OFFSET ?
	`

	args = append(args, limit, offset)

	err := r.db.ObtenerLista(ctx, &lista, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error inesperado, detalle: %w", err)
	}

	listaDto := make([]dto.PendienteDetalleNombre, 0, len(lista))

	for _, item := range lista {
		dtoItem := dto.PendienteDetalleNombre{
			ID:                item.ID,
			Titulo:            item.Titulo,
			Descripcion:       item.Descripcion,
			SolicitanteNombre: item.SolicitanteNombre,
			FechaPedido:       item.FechaPedido,
			AsignadoNombre:    item.AsignadoNombre,
			FechaAsignado:     item.FechaAsignado,
			Finalizado:        item.Finalizado,
			Cierre:            item.Cierre,
			FechaCierre:       item.FechaCierre,
		}

		listaDto = append(listaDto, dtoItem)
	}

	return listaDto, nil
}

//--------------------------------------------------------------

func (r *Repo) VincularDocumento(ctx context.Context, pendienteID int, documentoID string) error {

	strct := struct {
		PendienteID int    `db:"pendiente_id"`
		DocumentoID string `db:"documento_id"`
	}{
		PendienteID: pendienteID,
		DocumentoID: documentoID,
	}

	query := `INSERT INTO ti_pendientes_documento (pendiente_id, documento_id)
		 VALUES (:pendiente_id, :documento_id)`

	_, err := r.db.Cargar(ctx, &strct, query)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {

			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintForeignKey {
				return model.ErrFkNoEncontrado
			}

			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return model.ErrRelacionDuplicada
			}
		}

		return fmt.Errorf("error inesperado: %w", err)
	}

	return nil
}

func (r *Repo) VincularCodigoSAP(ctx context.Context, pendienteID int, codigoSAP string) error {

	strct := struct {
		PendienteID int    `db:"pendiente_id"`
		CodigoSAP   string `db:"codigo_sap_codigo"`
	}{
		PendienteID: pendienteID,
		CodigoSAP:   codigoSAP,
	}

	query := `INSERT INTO ti_pendientes_codigo_sap (pendiente_id, codigo_sap_codigo)
		 VALUES (:pendiente_id, :codigo_sap_codigo)`

	_, err := r.db.Cargar(ctx, &strct, query)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {

			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintForeignKey {
				return model.ErrFkNoEncontrado
			}

			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return model.ErrRelacionDuplicada
			}
		}

		return fmt.Errorf("error inesperado: %w", err)
	}

	return nil
}

func (r *Repo) VincularCodigoID(ctx context.Context, pendienteID int, codigoID string) error {

	strct := struct {
		PendienteID int    `db:"pendiente_id"`
		CodigoID    string `db:"codigo_id_codigo"`
	}{
		PendienteID: pendienteID,
		CodigoID:    codigoID,
	}

	query := `INSERT INTO ti_pendientes_codigo_id (pendiente_id, codigo_id_codigo)
		 VALUES (:pendiente_id, :codigo_id_codigo)`

	_, err := r.db.Cargar(ctx, &strct, query)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {

			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintForeignKey {
				return model.ErrFkNoEncontrado
			}

			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return model.ErrRelacionDuplicada
			}
		}

		return fmt.Errorf("error inesperado: %w", err)
	}

	return nil
}
