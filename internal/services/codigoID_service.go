package services

import (
	"context"
	"strings"
	"time"

	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"
)

// Interface para el acceso a datos de códigos ID. Se usa como contrato para el servicio
// Los métodos reciben context.Context para permitir trazabilidad/cancelación desde capas superiores.
type codigoIDManagerDbInterface interface {
	// Parámetros:
	//      - Contexto
	//      - Nuevo código ID (models.CodigoID)
	// Salida:
	//      - error
	// Errores que puede devuelve:
	//      - appErrors.ParametroDeCargaVacio
	//      - appErrors.CodigoIDVacio
	//      - appErrors.CodigoIDDuplicado
	Cargar(context.Context, *models.CodigoID) error

	// Parámetros:
	//      - Contexto
	//      - Código a buscar (string)
	// Salida:
	//      - *models.CodigoID
	// Errores que puede devuelve:
	//      - appErrors.CodigoIDVacio
	//      - appErrors.CodigoIDNoEncontrado
	ObtenerDetalle(context.Context, string) (*models.CodigoID, error)

	// Parámetros:
	//      - Contexto
	//      - Estado a filtrar (string)
	// Salida:
	//      - []models.CodigoID
	// Errores que puede devuelve:
	//      - appErrors.ParametroDeBusquedaVacio
	FiltrarPorEstado(context.Context, string) ([]models.CodigoID, error)

	// Parámetros:
	//      - Contexto
	//      - Código a actualizar (string)
	//      - Nuevo estado (string)
	//      - Fecha de actualización (time.Time)
	// Salida:
	//      - error
	// Errores que puede devuelve:
	//      - appErrors.CodigoIDVacio
	//      - appErrors.CodigoIDNoEncontrado
	ActualizarEstado(context.Context, string, string, time.Time) error

	// Parámetros:
	//      - Contexto
	//      - limit
	//      - offset
	// Salida:
	//      - []models.CodigoID
	// Errores que puede devuelve:
	//      - appErrors.ParametrosDeListaInvalidos
	Listar(context.Context, int, int) ([]models.CodigoID, error)
}

// Service
type codigoIDService struct {
	repo codigoIDManagerDbInterface
}

// NewCodigoIDService crea el servicio con la implementación del repo que cumpla la interfaz
func NewCodigoIDService(repo codigoIDManagerDbInterface) *codigoIDService {
	return &codigoIDService{repo: repo}
}

// Cargar crea un nuevo código ID validando parámetros mínimos
func (s *codigoIDService) Cargar(ctx context.Context, codigo string, descripcion string, estado string, fechaPedido time.Time) error {
	if strings.TrimSpace(codigo) == "" {
		return appErrors.ParametroDeCargaVacio
	}

	if strings.TrimSpace(estado) == "" {
		return appErrors.ParametroDeCargaVacio
	}

	if fechaPedido.After(time.Now()) {
		return appErrors.FechaNoValida
	}

	nuevo := &models.CodigoID{Codigo: codigo, Estado: estado, FechaPedido: fechaPedido}

	if strings.TrimSpace(descripcion) != "" {
		nuevo.Descripcion = &descripcion
	}

	return s.repo.Cargar(ctx, nuevo)
}

// ObtenerDetalle devuelve el detalle correspondiente al código provisto
func (s *codigoIDService) ObtenerDetalle(ctx context.Context, codigo string) (*models.CodigoID, error) {
	if strings.TrimSpace(codigo) == "" {
		return nil, appErrors.ParametroDeBusquedaVacio
	}
	return s.repo.ObtenerDetalle(ctx, codigo)
}

// FiltrarPorEstado filtra códigos ID por estado
func (s *codigoIDService) FiltrarPorEstado(ctx context.Context, estado string) ([]models.CodigoID, error) {
	if strings.TrimSpace(estado) == "" {
		return nil, appErrors.ParametroDeBusquedaVacio
	}
	return s.repo.FiltrarPorEstado(ctx, estado)
}

// ActualizarEstado actualiza el estado y fecha de actualización de un código ID
func (s *codigoIDService) ActualizarEstado(ctx context.Context, codigo string, nuevoEstado string, fecha time.Time) error {
	if strings.TrimSpace(codigo) == "" {
		return appErrors.ParametroDeBusquedaVacio
	}
	if strings.TrimSpace(nuevoEstado) == "" {
		return appErrors.ParametroDeCargaVacio
	}

	//validamos la fecha con un margen de +1 dia desde el momento que se carga
	limite := time.Now().Add(24 * time.Hour)
	if fecha.After(limite) {
		return appErrors.FechaNoValida
	}
	return s.repo.ActualizarEstado(ctx, codigo, nuevoEstado, fecha)
}

// Listar devuelve un conjunto paginado de códigos ID
func (s *codigoIDService) Listar(ctx context.Context, limit, offset int) ([]models.CodigoID, error) {
	if limit <= 0 || offset < 0 {
		return nil, appErrors.ParametrosDeListaInvalidos
	}
	return s.repo.Listar(ctx, limit, offset)
}
