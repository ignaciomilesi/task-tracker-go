package services

import (
	"context"

	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"
)

// Interface para el acceso a datos de solicitantes. Se usa como contrato para el servicio
// Los métodos reciben context.Context para permitir trazabilidad/cancelación desde capas superiores.
type solicitanteManagerDbInterface interface {
	// Parámetros:
	//      - Contexto
	//      - Solicitante a crear
	// Salida:
	//      - Id creado
	// Errores que puede devuelve:
	//      - appErrors.ParametroDeCargaVacio
	//      - appErrors.SolicitanteDuplicado
	Crear(context.Context, *models.Solicitante) (int, error)

	// Parámetros:
	//      - Contexto
	//      - Nombre a buscar
	// Salida:
	//      - Id del solicitante
	// Errores que puede devuelve:
	//      - appErrors.SolicitanteNoEncontrado
	ObtenerIDPorNombre(context.Context, string) (int, error)

	// Parámetros:
	//      - Contexto
	//      - Palabra a buscar
	// Salida:
	//      - Array con los solicitantes encontrados
	// Errores que puede devuelve:
	//      - appErrors.ParametroDeBusquedaVacio
	Buscar(context.Context, string) ([]models.Solicitante, error)

	// Parámetros:
	//      - Contexto
	//      - limit
	//      - offset
	// Salida:
	//      - Array con los solicitantes
	Listar(context.Context, int, int) ([]models.Solicitante, error)
}

// Service
type solicitanteService struct {
	repo solicitanteManagerDbInterface
}

// NewSolicitanteService crea el servicio con la implementación del repo que cumpla la interfaz
func NewSolicitanteService(repo solicitanteManagerDbInterface) *solicitanteService {
	return &solicitanteService{repo: repo}
}

// Crear crea un nuevo solicitante validando parámetros mínimos
func (s *solicitanteService) Crear(ctx context.Context, nombre string) (int, error) {
	if nombre == "" {
		return 0, appErrors.ParametroDeCargaVacio
	}

	nuevo := &models.Solicitante{Nombre: nombre}
	return s.repo.Crear(ctx, nuevo)
}

// ObtenerIDPorNombre devuelve el id correspondiente al nombre provisto
func (s *solicitanteService) ObtenerIDPorNombre(ctx context.Context, nombre string) (int, error) {
	if nombre == "" {
		return 0, appErrors.ParametroDeBusquedaVacio
	}
	return s.repo.ObtenerIDPorNombre(ctx, nombre)
}

// Buscar busca solicitantes por una subcadena en el nombre
func (s *solicitanteService) Buscar(ctx context.Context, parametro string) ([]models.Solicitante, error) {
	return s.repo.Buscar(ctx, parametro)
}

// Listar devuelve un conjunto paginado de solicitantes
func (s *solicitanteService) Listar(ctx context.Context, limit, offset int) ([]models.Solicitante, error) {
	// si limit es 0 o negativo, ajustamos un valor por defecto razonable
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.Listar(ctx, limit, offset)
}
