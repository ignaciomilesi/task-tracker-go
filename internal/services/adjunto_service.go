package services

import (
	"context"
	"strings"

	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"
)

// Interface para el acceso a datos de adjuntos. Se usa como contrato para el servicio
// Los métodos reciben context.Context para permitir trazabilidad/cancelación desde capas superiores.
type adjuntoManagerDbInterface interface {
	// Parámetros:
	//      - Contexto
	//      - Adjunto a crear (models.Adjunto)
	// Salida:
	//      - Id creado (int)
	// Errores que puede devuelve:
	//      - appErrors.ParametroDeCargaVacio
	//      - appErrors.PendienteNoEncontrado
	Cargar(context.Context, *models.Adjunto) (int, error)

	// Parámetros:
	//      - Contexto
	//      - Id del adjunto (int)
	// Salida:
	//      - *models.Adjunto
	// Errores que puede devuelve:
	//      - appErrors.AdjuntoNoEncontrado
	ObtenerDetalle(context.Context, int) (*models.Adjunto, error)

	// Parámetros:
	//      - Contexto
	//      - Id del adjunto (int)
	// Salida:
	//      - error
	// Errores que puede devuelve:
	//      - appErrors.AdjuntoNoEncontrado
	Eliminar(context.Context, int) error

	// Parámetros:
	//      - Contexto
	//      - Id del pendiente (int)
	// Salida:
	//      - []models.Adjunto
	// Errores que puede devuelve:
	//      - appErrors.ParametroDeBusquedaVacio
	FiltrarPorPendiente(context.Context, int) ([]models.Adjunto, error)

	// Parámetros:
	//      - Contexto
	//      - Id del adjunto (int)
	//      - Nuevo path (string)
	// Salida:
	//      - error
	// Errores que puede devuelve:
	//      - appErrors.ParametroDeBusquedaVacio
	//      - appErrors.ParametroDeCargaVacio
	//      - appErrors.AdjuntoNoEncontrado
	ActualizarPath(context.Context, int, string) error

	// Parámetros:
	//      - Contexto
	//      - Id del adjunto (int)
	//      - Nueva descripción (string)
	// Salida:
	//      - error
	// Errores que puede devuelve:
	//      - appErrors.ParametroDeBusquedaVacio
	//      - appErrors.ParametroDeCargaVacio
	//      - appErrors.AdjuntoNoEncontrado
	ActualizarDescripcion(context.Context, int, string) error
}

// Service
type adjuntoService struct {
	repo adjuntoManagerDbInterface
}

// NewAdjuntoService crea el servicio con la implementación del repo que cumpla la interfaz
func NewAdjuntoService(repo adjuntoManagerDbInterface) *adjuntoService {
	return &adjuntoService{repo: repo}
}

// Cargar crea un nuevo adjunto validando parámetros mínimos
func (s *adjuntoService) Cargar(ctx context.Context, pendienteID int, descripcion string, archivoPath string) (int, error) {
	if pendienteID <= 0 {
		return 0, appErrors.ParametroDeCargaVacio
	}
	if strings.TrimSpace(archivoPath) == "" {
		return 0, appErrors.ParametroDeCargaVacio
	}

	a := &models.Adjunto{PendienteID: pendienteID, Descripcion: descripcion, ArchivoPath: archivoPath}

	return s.repo.Cargar(ctx, a)
}

// ObtenerDetalle devuelve el detalle correspondiente al adjunto provisto
func (s *adjuntoService) ObtenerDetalle(ctx context.Context, id int) (*models.Adjunto, error) {
	if id <= 0 {
		return nil, appErrors.ParametroDeBusquedaVacio
	}
	return s.repo.ObtenerDetalle(ctx, id)
}

// Eliminar elimina un adjunto por id
func (s *adjuntoService) Eliminar(ctx context.Context, id int) error {
	if id <= 0 {
		return appErrors.ParametroDeBusquedaVacio
	}
	return s.repo.Eliminar(ctx, id)
}

// FiltrarPorPendiente devuelve adjuntos asociados a un pendiente
func (s *adjuntoService) FiltrarPorPendiente(ctx context.Context, pendienteID int) ([]models.Adjunto, error) {
	if pendienteID <= 0 {
		return nil, appErrors.ParametroDeBusquedaVacio
	}
	return s.repo.FiltrarPorPendiente(ctx, pendienteID)
}

// ActualizarPath actualiza el path de un adjunto
func (s *adjuntoService) ActualizarPath(ctx context.Context, id int, path string) error {
	if id <= 0 {
		return appErrors.ParametroDeBusquedaVacio
	}
	if strings.TrimSpace(path) == "" {
		return appErrors.ParametroDeCargaVacio
	}
	return s.repo.ActualizarPath(ctx, id, path)
}

// ActualizarDescripcion actualiza la descripción de un adjunto
func (s *adjuntoService) ActualizarDescripcion(ctx context.Context, id int, descripcion string) error {
	if id <= 0 {
		return appErrors.ParametroDeBusquedaVacio
	}
	if strings.TrimSpace(descripcion) == "" {
		return appErrors.ParametroDeCargaVacio
	}
	return s.repo.ActualizarDescripcion(ctx, id, descripcion)
}
