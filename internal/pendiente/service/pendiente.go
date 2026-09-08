package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"task-tracker-go/internal/pendiente/dto"
	"task-tracker-go/internal/pendiente/model"
)

type repoPendienteInterface interface {
	Crear(ctx context.Context, nuevoPendiente *model.Pendiente) (int, error)
	ObtenerDetalle(ctx context.Context, id int) (*dto.PendienteDetalleNombre, error)
	Actualizar(ctx context.Context, pendiente *model.Pendiente) error
	BuscarPorTituloDescripcion(ctx context.Context, texto string, finalizado bool) ([]dto.PendienteDetalleNombre, error)
	Listar(ctx context.Context, finalizado bool, asignadoID *int, limit, offset int) ([]dto.PendienteDetalleNombre, error)

	VincularDocumento(ctx context.Context, pendienteID int, documentoID string) error
	VincularCodigoSAP(ctx context.Context, pendienteID int, codigoSAP string) error
	VincularCodigoID(ctx context.Context, pendienteID int, codigoID string) error
}

type ServicePendiente struct {
	repoPend repoPendienteInterface
	repoAvan repoAvanceInterface
	repoAdj  repoAdjuntoInterface
}

func NewServicePendiente(repo repoPendienteInterface, repoAvan repoAvanceInterface, repoAdj repoAdjuntoInterface) *ServicePendiente {
	return &ServicePendiente{repoPend: repo, repoAvan: repoAvan, repoAdj: repoAdj}
}

// creo y registra el avance
func (s *ServicePendiente) Crear(ctx context.Context, nuevaSolicitudPendiente *dto.PendienteRequest) error {

	if nuevaSolicitudPendiente == nil {
		return fmt.Errorf("service: %w", model.ErrModelDeCargaNil)
	}

	if strings.TrimSpace(nuevaSolicitudPendiente.Titulo) == "" {
		return fmt.Errorf("service: %w, parametro: titulo", model.ErrParametrosNoValidos)
	}
	if strings.TrimSpace(nuevaSolicitudPendiente.Descripcion) == "" {
		return fmt.Errorf("service: %w, parametro: descripcion", model.ErrParametrosNoValidos)
	}
	if nuevaSolicitudPendiente.SolicitanteID <= 0 {
		return fmt.Errorf("service: %w, parametro: solicitante_id", model.ErrParametrosNoValidos)
	}
	if nuevaSolicitudPendiente.FechaPedido.After(time.Now()) {
		return fmt.Errorf("service: %w, parametro: fecha_pedido", model.ErrParametrosNoValidos)
	}
	if nuevaSolicitudPendiente.MailPath != nil && strings.TrimSpace(*nuevaSolicitudPendiente.MailPath) == "" {
		nuevaSolicitudPendiente.MailPath = nil
	}

	nuevopendiente := model.Pendiente{
		Titulo:        nuevaSolicitudPendiente.Titulo,
		Descripcion:   nuevaSolicitudPendiente.Descripcion,
		SolicitanteID: nuevaSolicitudPendiente.SolicitanteID,
		FechaPedido:   nuevaSolicitudPendiente.FechaPedido,
	}
	// creo el pendiente
	id, err := s.repoPend.Crear(ctx, &nuevopendiente)
	if err != nil {
		return fmt.Errorf("service: %w", err)
	}

	// registro el avance, el fallo no es crítico
	go func() {
		nuevoAvance := model.Avance{
			PendienteID: id,
			Descripcion: "pendiente creado",
			Fecha:       nuevopendiente.FechaPedido,
			MailPath:    nuevaSolicitudPendiente.MailPath,
		}

		if _, err := s.repoAvan.CargarAvance(ctx, &nuevoAvance); err != nil {
			fmt.Println("Error al cargar avance:", err)
		}
	}()

	return nil
}

func (s *ServicePendiente) Actualizar(ctx context.Context, solicitud *dto.ActualizacionPendienteRequest) error {

	if solicitud == nil {
		return fmt.Errorf("service: %w", model.ErrModelDeCargaNil)
	}

	if strings.TrimSpace(solicitud.Titulo) == "" {
		return fmt.Errorf("service: %w, parametro: titulo", model.ErrParametrosNoValidos)
	}

	if solicitud.SolicitanteID <= 0 {
		return fmt.Errorf("service: %w, parametro: solicitante_id", model.ErrParametrosNoValidos)
	}
	if solicitud.FechaPedido.After(time.Now()) {
		return fmt.Errorf("service: %w, parametro: fecha_pedido", model.ErrParametrosNoValidos)
	}
	if solicitud.AsignadoID != nil && *solicitud.AsignadoID <= 0 {
		return fmt.Errorf("service: %w, parametro: asignado_id", model.ErrParametrosNoValidos)
	}
	if solicitud.FechaAsignado != nil && solicitud.FechaAsignado.After(time.Now()) {
		return fmt.Errorf("service: %w, parametro: fecha_asignado", model.ErrParametrosNoValidos)
	}
	if solicitud.FechaCierre != nil && solicitud.FechaCierre.After(time.Now()) {
		return fmt.Errorf("service: %w, parametro: fecha_cierre", model.ErrParametrosNoValidos)
	}

	nuevoPendiente := model.Pendiente{
		ID:            solicitud.ID,
		Titulo:        solicitud.Titulo,
		Descripcion:   solicitud.Descripcion,
		SolicitanteID: solicitud.SolicitanteID,
		FechaPedido:   solicitud.FechaPedido,
		AsignadoID:    solicitud.AsignadoID,
		FechaAsignado: solicitud.FechaAsignado,
		Finalizado:    solicitud.Finalizado,
		Cierre:        solicitud.Cierre,
		FechaCierre:   solicitud.FechaCierre,
	}

	if err := s.repoPend.Actualizar(ctx, &nuevoPendiente); err != nil {
		return fmt.Errorf("service: %w", err)
	}

	return nil
}

func (s *ServicePendiente) BuscarPorTituloDescripcion(ctx context.Context, texto string, finalizado bool) ([]dto.PendienteDetalleNombre, error) {
	if strings.TrimSpace(texto) == "" {
		return nil, fmt.Errorf("service: %w, texto de búsqueda vacio", model.ErrParametrosNoValidos)
	}
	return s.repoPend.BuscarPorTituloDescripcion(ctx, texto, finalizado)
}

func (s *ServicePendiente) Listar(ctx context.Context, asignadoID *int, finalizado bool, limit, offset int) ([]dto.PendienteDetalleNombre, error) {
	if limit <= 0 || offset < 0 {
		return nil, fmt.Errorf("service: %w, parámetros de lista no valido", model.ErrParametrosNoValidos)
	}
	return s.repoPend.Listar(ctx, finalizado, asignadoID, limit, offset)
}

func (s *ServicePendiente) ObtenerDetalle(ctx context.Context, id int) (*dto.PendienteDetalleCompleto, error) {
	if id <= 0 {
		return nil, fmt.Errorf("service: %w, pendienteID no valido", model.ErrParametrosNoValidos)
	}

	var pendienteDetalleCompleto dto.PendienteDetalleCompleto
	var wg sync.WaitGroup
	wg.Add(2)

	//buscamos los avances
	go func() {
		defer wg.Done()
		avances, err := s.repoAvan.ObtenerListaAvances(ctx, id)
		if err != nil {
			fmt.Printf("service: error al obtener avances al detallar el services, %v", err)
		}
		pendienteDetalleCompleto.Avances = avances
	}()

	//buscamos los adjuntos
	go func() {
		defer wg.Done()
		adjuntos, err := s.repoAdj.ObtenerListaAdjuntos(ctx, id)
		if err != nil {
			fmt.Printf("service: error al obtener adjuntos al detallar el services, %v", err)
		}
		pendienteDetalleCompleto.Adjuntos = adjuntos

	}()

	detalle, err := s.repoPend.ObtenerDetalle(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service: %w", err)
	}
	pendienteDetalleCompleto.PendienteDetalleNombre = *detalle

	// esperamos que terminen todas las búsquedas de avances y adjuntos
	wg.Wait()
	return &pendienteDetalleCompleto, nil

}

// ---------------------- Vincular documento, codigoID y codigoSAP ----------------------

func (s *ServicePendiente) VincularDocumento(ctx context.Context, pendienteID int, documentoID string) error {
	if pendienteID <= 0 {
		return fmt.Errorf("service: %w, pendienteID inválido", model.ErrParametrosNoValidos)
	}
	if strings.TrimSpace(documentoID) == "" {
		return fmt.Errorf("service: %w, documentoID inválido", model.ErrParametrosNoValidos)
	}
	return s.repoPend.VincularDocumento(ctx, pendienteID, documentoID)
}

func (s *ServicePendiente) VincularCodigoSAP(ctx context.Context, pendienteID int, codigoSAP string) error {
	if pendienteID <= 0 {
		return fmt.Errorf("service: %w, pendienteID inválido", model.ErrParametrosNoValidos)
	}
	if strings.TrimSpace(codigoSAP) == "" {
		return fmt.Errorf("service: %w, codigoSAP inválido", model.ErrParametrosNoValidos)
	}
	return s.repoPend.VincularCodigoSAP(ctx, pendienteID, codigoSAP)
}

func (s *ServicePendiente) VincularCodigoID(ctx context.Context, pendienteID int, codigoID string) error {
	if pendienteID <= 0 {
		return fmt.Errorf("service: %w, pendienteID inválido", model.ErrParametrosNoValidos)
	}
	if strings.TrimSpace(codigoID) == "" {
		return fmt.Errorf("service: %w, codigoID inválido", model.ErrParametrosNoValidos)
	}
	return s.repoPend.VincularCodigoID(ctx, pendienteID, codigoID)
}
