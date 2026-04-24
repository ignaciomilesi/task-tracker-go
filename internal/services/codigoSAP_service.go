package services

import (
	"context"
	"strings"

	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"
)

// Interface para el acceso a datos de códigos SAP. Se usa como contrato para el servicio
// Los métodos reciben context.Context para permitir trazabilidad/cancelación desde capas superiores.
type codigoSAPManagerDbInterface interface {
	// Parámetros:
	//      - Contexto
	//      - Código a crear (models.CodigoSAP)
	// Salida:
	//      - error
	// Errores que puede devuelve:
	//      - appErrors.ParametroDeCargaVacio
	//      - appErrors.CodigoSAPVacio
	//      - appErrors.CodigoSAPDuplicado
	Cargar(context.Context, *models.CodigoSAP) error

	// Parámetros:
	//      - Contexto
	//      - Código a buscar (string)
	// Salida:
	//      - *models.CodigoSAP
	// Errores que puede devuelve:
	//      - appErrors.CodigoSAPNoEncontrado
	ObtenerDetalle(context.Context, string) (*models.CodigoSAP, error)

	// Parámetros:
	//      - Contexto
	//      - Palabra a buscar en la descripción (string)
	// Salida:
	//      - []models.CodigoSAP
	// Errores que puede devuelve:
	//      - appErrors.ParametroDeBusquedaVacio
	BuscarPorDescripcion(context.Context, string) ([]models.CodigoSAP, error)

	// Parámetros:
	//      - Contexto
	//      - Código a modificar (string)
	//      - Nueva descripción (string)
	// Salida:
	//      - error
	// Errores que puede devuelve:
	//      - appErrors.CodigoSAPVacio
	//      - appErrors.CodigoSAPNoEncontrado
	ModificarDescripcion(context.Context, string, string) error

	// Parámetros:
	//      - Contexto
	//      - limit
	//      - offset
	// Salida:
	//      - []models.CodigoSAP
	// Errores que puede devuelve:
	//      - appErrors.ParametrosDeListaInvalidos
	Listar(context.Context, int, int) ([]models.CodigoSAP, error)
}

// Service
type codigoSAPService struct {
	repo codigoSAPManagerDbInterface
}

// NewCodigoSAPService crea el servicio con la implementación del repo que cumpla la interfaz
func NewCodigoSAPService(repo codigoSAPManagerDbInterface) *codigoSAPService {
	return &codigoSAPService{repo: repo}
}

// Cargar crea un nuevo código SAP validando parámetros mínimos
func (s *codigoSAPService) Cargar(ctx context.Context, codigo string, descripcion string) error {
	if strings.TrimSpace(codigo) == "" {
		return appErrors.ParametroDeCargaVacio
	}

	nuevo := &models.CodigoSAP{Codigo: codigo}

	if strings.TrimSpace(descripcion) != "" {
		nuevo.Descripcion = &descripcion
	}
	return s.repo.Cargar(ctx, nuevo)
}

// ObtenerDetalle devuelve el detalle correspondiente al código provisto
func (s *codigoSAPService) ObtenerDetalle(ctx context.Context, codigo string) (*models.CodigoSAP, error) {
	if strings.TrimSpace(codigo) == "" {
		return nil, appErrors.ParametroDeBusquedaVacio
	}
	return s.repo.ObtenerDetalle(ctx, codigo)
}

// BuscarPorDescripcion busca códigos SAP por una subcadena en la descripción
func (s *codigoSAPService) BuscarPorDescripcion(ctx context.Context, parametro string) ([]models.CodigoSAP, error) {
	if strings.TrimSpace(parametro) == "" {
		return nil, appErrors.ParametroDeBusquedaVacio
	}
	return s.repo.BuscarPorDescripcion(ctx, parametro)
}

// ModificarDescripcion actualiza la descripción de un código SAP
func (s *codigoSAPService) ModificarDescripcion(ctx context.Context, codigo string, nuevaDescripcion string) error {
	if strings.TrimSpace(codigo) == "" {
		return appErrors.ParametroDeBusquedaVacio
	}
	if strings.TrimSpace(nuevaDescripcion) == "" {
		return appErrors.ParametroDeCargaVacio
	}
	return s.repo.ModificarDescripcion(ctx, codigo, nuevaDescripcion)
}

// Listar devuelve un conjunto paginado de códigos SAP
func (s *codigoSAPService) Listar(ctx context.Context, limit, offset int) ([]models.CodigoSAP, error) {
	if limit <= 0 || offset < 0 {
		return nil, appErrors.ParametrosDeListaInvalidos
	}
	return s.repo.Listar(ctx, limit, offset)
}
