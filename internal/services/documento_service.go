package services

import (
	"context"
	"strings"

	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"
)

// Interface para el acceso a datos de documentos. Se usa como contrato para el servicio
// Los métodos reciben context.Context para permitir trazabilidad/cancelación desde capas superiores.
type documentoManagerDbInterface interface {
	// Parámetros:
	//      - Contexto
	//      - Documento a crear (models.Documento)
	// Salida:
	//      - error
	// Errores que puede devuelve:
	//      - appErrors.ParametroDeCargaVacio
	//      - appErrors.CodigoDocumentoVacio
	//      - appErrors.DocumentoDuplicado
	Cargar(context.Context, *models.Documento) error

	// Parámetros:
	//      - Contexto
	//      - Código a buscar (string)
	// Salida:
	//      - *models.Documento
	// Errores que puede devuelve:
	//      - appErrors.ParametroDeBusquedaVacio
	//      - appErrors.DocumentoNoEncontrado
	ObtenerDetalle(context.Context, string) (*models.Documento, error)

	// Parámetros:
	//      - Contexto
	//      - Tipo a filtrar (string)
	// Salida:
	//      - []models.Documento
	// Errores que puede devuelve:
	//      - appErrors.ParametroDeBusquedaVacio
	FiltrarPorTipo(context.Context, string) ([]models.Documento, error)

	// Parámetros:
	//      - Contexto
	//      - Título a buscar (string)
	// Salida:
	//      - []models.Documento
	// Errores que puede devuelve:
	//      - appErrors.ParametroDeBusquedaVacio
	FiltrarPorTitulo(context.Context, string) ([]models.Documento, error)

	// Parámetros:
	//      - Contexto
	//      - Código a actualizar (string)
	//      - Nuevo path (string)
	// Salida:
	//      - error
	// Errores que puede devuelve:
	//      - appErrors.ParametroDeBusquedaVacio
	//      - appErrors.ParametroDeCargaVacio
	//      - appErrors.DocumentoNoEncontrado
	ActualizarPath(context.Context, string, string) error

	// Parámetros:
	//      - Contexto
	//      - Código a actualizar (string)
	//      - Nuevo backup path (string)
	// Salida:
	//      - error
	// Errores que puede devuelve:
	//      - appErrors.ParametroDeBusquedaVacio
	//      - appErrors.ParametroDeCargaVacio
	//      - appErrors.DocumentoNoEncontrado
	ActualizarBackupPath(context.Context, string, string) error
}

// Service
type documentoService struct {
	repo documentoManagerDbInterface
}

// NewDocumentoService crea el servicio con la implementación del repo que cumpla la interfaz
func NewDocumentoService(repo documentoManagerDbInterface) *documentoService {
	return &documentoService{repo: repo}
}

// Cargar crea un nuevo documento validando parámetros mínimos
func (s *documentoService) Cargar(ctx context.Context, codigo, emision, titulo, tipo string, ubicacionPath, backupPath string) error {
	if strings.TrimSpace(emision) == "" ||
		strings.TrimSpace(titulo) == "" ||
		strings.TrimSpace(tipo) == "" {
		return appErrors.ParametroDeCargaVacio
	}
	if strings.TrimSpace(codigo) == "" {
		return appErrors.CodigoDocumentoVacio
	}

	doc := &models.Documento{Codigo: codigo, Emision: emision, Titulo: titulo, Tipo: tipo}

	if strings.TrimSpace(ubicacionPath) != "" {
		doc.UbicacionPath = &ubicacionPath
	}
	if strings.TrimSpace(backupPath) != "" {
		doc.BackupPath = &backupPath
	}

	return s.repo.Cargar(ctx, doc)
}

// ObtenerDetalle devuelve el detalle correspondiente al código provisto
func (s *documentoService) ObtenerDetalle(ctx context.Context, codigo string) (*models.Documento, error) {
	if strings.TrimSpace(codigo) == "" {
		return nil, appErrors.ParametroDeBusquedaVacio
	}
	return s.repo.ObtenerDetalle(ctx, codigo)
}

// FiltrarPorTipo filtra documentos por tipo
func (s *documentoService) FiltrarPorTipo(ctx context.Context, tipo string) ([]models.Documento, error) {
	if strings.TrimSpace(tipo) == "" {
		return nil, appErrors.ParametroDeBusquedaVacio
	}
	return s.repo.FiltrarPorTipo(ctx, tipo)
}

// FiltrarPorTitulo busca documentos por título (LIKE)
func (s *documentoService) FiltrarPorTitulo(ctx context.Context, titulo string) ([]models.Documento, error) {
	if strings.TrimSpace(titulo) == "" {
		return nil, appErrors.ParametroDeBusquedaVacio
	}
	return s.repo.FiltrarPorTitulo(ctx, titulo)
}

// ActualizarPath actualiza el path de un documento
func (s *documentoService) ActualizarPath(ctx context.Context, codigo string, nuevoPath string) error {
	if strings.TrimSpace(codigo) == "" {
		return appErrors.ParametroDeBusquedaVacio
	}
	if strings.TrimSpace(nuevoPath) == "" {
		return appErrors.ParametroDeCargaVacio
	}
	return s.repo.ActualizarPath(ctx, codigo, nuevoPath)
}

// ActualizarBackupPath actualiza el backup path de un documento
func (s *documentoService) ActualizarBackupPath(ctx context.Context, codigo string, nuevoPath string) error {
	if strings.TrimSpace(codigo) == "" {
		return appErrors.ParametroDeBusquedaVacio
	}
	if strings.TrimSpace(nuevoPath) == "" {
		return appErrors.ParametroDeCargaVacio
	}
	return s.repo.ActualizarBackupPath(ctx, codigo, nuevoPath)
}
