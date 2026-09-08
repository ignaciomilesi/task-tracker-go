package documento

import (
	"context"
	"fmt"
	"strings"
)

type repoInterface interface {
	Cargar(context.Context, *Documento) error
	ObtenerDetalle(context.Context, string) (*Documento, error)
	FiltrarPorTipo(context.Context, string) ([]Documento, error)
	FiltrarPorTitulo(context.Context, string) ([]Documento, error)
	ActualizarPath(context.Context, string, string) error
	ActualizarBackupPath(context.Context, string, string) error
	PendientesRelacionados(context.Context, string) ([]pendiente, error)
}

type service struct {
	repo repoInterface
}

func newService(repo repoInterface) *service {
	return &service{repo: repo}
}

func (s *service) Cargar(ctx context.Context, nuevaSolicitudDocumento *nuevoDocumentoRequest) error {
	if nuevaSolicitudDocumento == nil {
		return fmt.Errorf("service: %w", errModelDeCargaNil)
	}

	nuevoDocumento, err := NewDocumento(
		nuevaSolicitudDocumento.Codigo,
		nuevaSolicitudDocumento.Emision,
		nuevaSolicitudDocumento.Titulo,
		nuevaSolicitudDocumento.Tipo,
		nuevaSolicitudDocumento.UbicacionPath,
		nuevaSolicitudDocumento.BackupPath,
	)

	if err != nil {
		return fmt.Errorf("service: %w", errModelDeCargaNil)
	}

	return s.repo.Cargar(ctx, nuevoDocumento)
}

func (s *service) ObtenerDetalle(ctx context.Context, codigo string) (*Documento, error) {
	if strings.TrimSpace(codigo) == "" {
		return nil, fmt.Errorf("service: %w, codigo vacío", errParametrosNoValidos)
	}
	return s.repo.ObtenerDetalle(ctx, codigo)
}

func (s *service) FiltrarPorTipo(ctx context.Context, tipo string) ([]Documento, error) {
	if strings.TrimSpace(tipo) == "" {
		return nil, fmt.Errorf("service: %w, tipo vacío", errParametrosNoValidos)
	}
	return s.repo.FiltrarPorTipo(ctx, tipo)
}

func (s *service) FiltrarPorTitulo(ctx context.Context, titulo string) ([]Documento, error) {
	if strings.TrimSpace(titulo) == "" {
		return nil, fmt.Errorf("service: %w, titulo vacío", errParametrosNoValidos)
	}
	return s.repo.FiltrarPorTitulo(ctx, titulo)
}

func (s *service) ActualizarPath(ctx context.Context, codigo, nuevoPath, tipo string) error {
	if strings.TrimSpace(codigo) == "" {
		return fmt.Errorf("service: %w, codigo vacío", errParametrosNoValidos)
	}
	if strings.TrimSpace(nuevoPath) == "" {
		return fmt.Errorf("service: %w, nuevoPath vacío", errParametrosNoValidos)
	}

	if strings.TrimSpace(tipo) == "backup" {
		return s.repo.ActualizarBackupPath(ctx, codigo, nuevoPath)
	}
	if strings.TrimSpace(tipo) == "merge" {
		return s.repo.ActualizarPath(ctx, codigo, nuevoPath)
	}
	return fmt.Errorf("service: %w, parametro 'path' debe ser 'backup' o 'merge'", errParametrosNoValidos)
}

func (s *service) PendientesRelacionados(ctx context.Context, documentoID string) ([]pendiente, error) {
	if strings.TrimSpace(documentoID) == "" {
		return nil, fmt.Errorf("service: %w, parametro de busqueda vacío", errParametrosNoValidos)
	}
	return s.repo.PendientesRelacionados(ctx, documentoID)
}
