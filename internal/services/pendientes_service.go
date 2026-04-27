package services

import (
	"context"
	"strings"
	"time"

	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"
)

// Interface para el acceso a datos de pendientes. Se usa como contrato para el servicio
// Los métodos reciben context.Context para permitir trazabilidad/cancelación desde capas superiores.
type pendientesManagerDbInterface interface {
	// Parámetros:
	//      - Contexto
	//      - Pendiente a crear (models.Pendientes)
	// Salida:
	//      - Id creado (int)
	// Errores que puede devuelve:
	//      - appErrors.ParametroDeCargaVacio
	//      - appErrors.SolicitanteNoEncontrado
	Crear(context.Context, *models.Pendientes) (int, error)

	// Parámetros:
	//      - Contexto
	//      - Id del pendiente (int)
	//      - Id del colaborador asignado (int)
	//      - Fecha de asignación (time.Time)
	// Salida:
	//      - error
	// Errores que puede devuelve:
	//      - appErrors.ColaboradorNoEncontrado
	//      - appErrors.PendienteNoEncontrado
	Asignar(context.Context, int, int, time.Time) error

	// Parámetros:
	//      - Contexto
	//      - Id del pendiente (int)
	//      - Cierre (opcional *string)
	//      - Fecha de cierre (time.Time)
	// Salida:
	//      - error
	// Errores que puede devuelve:
	//      - appErrors.PendienteNoEncontrado
	Terminar(context.Context, int, *string, time.Time) error

	// Parámetros:
	//      - Contexto
	//      - Id del pendiente (int)
	//      - Nueva descripción (string)
	// Salida:
	//      - error
	// Errores que puede devuelve:
	//      - appErrors.PendienteNoEncontrado
	ModificarDescripcion(context.Context, int, string) error

	// Parámetros:
	//      - Contexto
	//      - Id del pendiente (int)
	//      - Nueva identificación de tabla pendiente (string)
	// Salida:
	//      - error
	// Errores que puede devuelve:
	//      - appErrors.PendienteNoEncontrado
	ModificarIdentificacionTablaPendiente(context.Context, int, string) error

	// Parámetros:
	//      - Contexto
	//      - Texto a buscar en título o descripción (string)
	//      - finalizado (bool)
	// Salida:
	//      - []models.Pendientes
	// Errores que puede devuelve:
	//      - appErrors.ParametroDeBusquedaVacio
	BuscarPorTituloDescripcion(context.Context, string, bool) ([]models.Pendientes, error)

	// Parámetros:
	//      - Contexto
	//      - Id del colaborador asignado (int)
	//      - finalizado (bool)
	// Salida:
	//      - []models.Pendientes
	ListarPorAsignado(context.Context, int, bool) ([]models.Pendientes, error)

	// Parámetros:
	//      - Contexto
	//      - finalizado (bool)
	//      - limit
	//      - offset
	// Salida:
	//      - []models.Pendientes
	// Errores que puede devuelve:
	//      - appErrors.ParametrosDeListaInvalidos
	Listar(context.Context, bool, int, int) ([]models.Pendientes, error)
}

// Service
type pendientesService struct {
	repo pendientesManagerDbInterface
}

// NewPendientesService crea el servicio con la implementación del repo que cumpla la interfaz
func NewPendientesService(repo pendientesManagerDbInterface) *pendientesService {
	return &pendientesService{repo: repo}
}

// Crear crea un nuevo pendiente validando parámetros mínimos
func (s *pendientesService) Crear(ctx context.Context, titulo, descripcion string, solicitanteID int, fechaPedido time.Time) (int, error) {
	if strings.TrimSpace(titulo) == "" || strings.TrimSpace(descripcion) == "" {
		return 0, appErrors.ParametroDeCargaVacio
	}
	if solicitanteID <= 0 {
		return 0, appErrors.ParametroDeCargaVacio
	}
	if fechaPedido.After(time.Now()) {
		return 0, appErrors.FechaNoValida
	}

	p := &models.Pendientes{Titulo: titulo, Descripcion: descripcion, SolicitanteID: solicitanteID, FechaPedido: fechaPedido}

	return s.repo.Crear(ctx, p)
}

// Asignar asigna un colaborador a un pendiente
func (s *pendientesService) Asignar(ctx context.Context, id int, asignadoID int, fecha time.Time) error {
	if id <= 0 || asignadoID <= 0 {
		return appErrors.ParametroDeBusquedaVacio
	}
	if fecha.After(time.Now()) {
		return appErrors.FechaNoValida
	}
	return s.repo.Asignar(ctx, id, asignadoID, fecha)
}

// Terminar marca un pendiente como finalizado
func (s *pendientesService) Terminar(ctx context.Context, id int, cierre string, fecha time.Time) error {
	if id <= 0 {
		return appErrors.ParametroDeBusquedaVacio
	}

	if fecha.After(time.Now()) {
		return appErrors.FechaNoValida
	}

	var cierreCarga *string
	if strings.TrimSpace(cierre) != "" {
		cierreCarga = &cierre
	}
	return s.repo.Terminar(ctx, id, cierreCarga, fecha)
}

// ModificarDescripcion actualiza la descripción de un pendiente
func (s *pendientesService) ModificarDescripcion(ctx context.Context, id int, descripcion string) error {
	if id <= 0 {
		return appErrors.ParametroDeBusquedaVacio
	}
	if strings.TrimSpace(descripcion) == "" {
		return appErrors.ParametroDeCargaVacio
	}
	return s.repo.ModificarDescripcion(ctx, id, descripcion)
}

// ModificarIdentificacionTablaPendiente actualiza la identificación de tabla del pendiente
func (s *pendientesService) ModificarIdentificacionTablaPendiente(ctx context.Context, id int, identificacion string) error {
	if id <= 0 {
		return appErrors.ParametroDeBusquedaVacio
	}
	if strings.TrimSpace(identificacion) == "" {
		return appErrors.ParametroDeCargaVacio
	}
	return s.repo.ModificarIdentificacionTablaPendiente(ctx, id, identificacion)
}

// BuscarPorTituloDescripcion busca pendientes por título o descripción
func (s *pendientesService) BuscarPorTituloDescripcion(ctx context.Context, texto string, finalizado bool) ([]models.Pendientes, error) {
	if strings.TrimSpace(texto) == "" {
		return nil, appErrors.ParametroDeBusquedaVacio
	}
	return s.repo.BuscarPorTituloDescripcion(ctx, texto, finalizado)
}

// ListarPorAsignado lista pendientes por asignado
func (s *pendientesService) ListarPorAsignado(ctx context.Context, asignadoID int, finalizado bool) ([]models.Pendientes, error) {
	if asignadoID <= 0 {
		return nil, appErrors.ParametroDeBusquedaVacio
	}
	return s.repo.ListarPorAsignado(ctx, asignadoID, finalizado)
}

// Listar devuelve un conjunto paginado de pendientes
func (s *pendientesService) Listar(ctx context.Context, finalizado bool, limit, offset int) ([]models.Pendientes, error) {
	if limit <= 0 || offset < 0 {
		return nil, appErrors.ParametrosDeListaInvalidos
	}
	return s.repo.Listar(ctx, finalizado, limit, offset)
}
