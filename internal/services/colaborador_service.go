package services

import (
	"context"
	"strings"

	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"
)

// Interface para el acceso a datos de colaboradores. Se usa como contrato para el servicio
// Los métodos reciben context.Context para permitir trazabilidad/cancelación desde capas superiores.
type colaboradorManagerDbInterface interface {
	// Parámetros:
	//      - Contexto
	//      - Colaborador a crear
	// Salida:
	//      - Id creado
	// Errores que puede devuelve:
	//      - appErrors.ParametroDeCargaVacio
	//      - appErrors.ColaboradorDuplicado
	Crear(context.Context, *models.Colaborador) (int, error)

	// Parámetros:
	//      - Contexto
	//      - Nombre a buscar
	// Salida:
	//      - Id del colaborador
	// Errores que puede devuelve:
	//      - appErrors.ColaboradorNoEncontrado
	ObtenerIDPorNombre(context.Context, string) (int, error)

	// Parámetros:
	//      - Contexto
	//      - Palabra a buscar
	// Salida:
	//      - Array con los colaboradores encontrados
	// Errores que puede devuelve:
	//      - appErrors.ParametroDeBusquedaVacio
	Buscar(context.Context, string) ([]models.Colaborador, error)

	// Parámetros:
	//      - Contexto
	//      - limit
	//      - offset
	// Salida:
	//      - Array con los colaboradores
	Listar(context.Context, int, int) ([]models.Colaborador, error)
}

// Service
type colaboradorService struct {
	repo colaboradorManagerDbInterface
}

// NewColaboradorService crea el servicio con la implementación del repo que cumpla la interfaz
func NewColaboradorService(repo colaboradorManagerDbInterface) *colaboradorService {
	return &colaboradorService{repo: repo}
}

// Crear crea un nuevo colaborador validando parámetros mínimos
func (s *colaboradorService) Crear(ctx context.Context, nombre string) (int, error) {
	if strings.TrimSpace(nombre) == "" {
		return 0, appErrors.ParametroDeCargaVacio
	}

	nuevo := &models.Colaborador{Nombre: nombre}
	return s.repo.Crear(ctx, nuevo)
}

// ObtenerIDPorNombre devuelve el id correspondiente al nombre provisto
func (s *colaboradorService) ObtenerIDPorNombre(ctx context.Context, nombre string) (int, error) {
	if strings.TrimSpace(nombre) == "" {
		return 0, appErrors.ParametroDeBusquedaVacio
	}
	return s.repo.ObtenerIDPorNombre(ctx, nombre)
}

// Buscar busca colaboradores por una subcadena en el nombre
func (s *colaboradorService) Buscar(ctx context.Context, parametro string) ([]models.Colaborador, error) {
	if strings.TrimSpace(parametro) == "" {
		return nil, appErrors.ParametroDeBusquedaVacio
	}
	return s.repo.Buscar(ctx, parametro)
}

// Listar devuelve un conjunto paginado de colaboradores
func (s *colaboradorService) Listar(ctx context.Context, limit, offset int) ([]models.Colaborador, error) {

	if limit <= 0 || offset < 0 {
		return nil, appErrors.ParametrosDeListaInvalidos
	}
	return s.repo.Listar(ctx, limit, offset)
}
