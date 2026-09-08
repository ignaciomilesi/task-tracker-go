package codigoSAP

import (
	"context"
	"fmt"
	"strings"
)

type repoInterface interface {
	Cargar(context.Context, *CodigoSAP) error
	ObtenerDetalle(context.Context, string) (*CodigoSAP, error)
	BuscarPorDescripcion(context.Context, string) ([]CodigoSAP, error)
	ModificarDescripcion(context.Context, string, string) error
	Listar(context.Context, int, int) ([]CodigoSAP, error)
	PendientesRelacionados(context.Context, string) ([]pendiente, error)
}

type service struct {
	repo repoInterface
}

func newService(repo repoInterface) *service {
	return &service{repo: repo}
}

// funciona como DTO
type nuevoCodigoSAPRequest struct {
	Codigo      string
	Descripcion *string
}

func (s *service) Cargar(ctx context.Context, nuevoSolicitudCodigo *nuevoCodigoSAPRequest) error {
	if nuevoSolicitudCodigo == nil {
		return fmt.Errorf("service: %w", errModelDeCargaNil)
	}
	nuevoCodigo, err := NewCodigoID(
		nuevoSolicitudCodigo.Codigo,
		nuevoSolicitudCodigo.Descripcion,
	)
	if err != nil {
		return fmt.Errorf("service: %w", errModelDeCargaNil)
	}

	return s.repo.Cargar(ctx, nuevoCodigo)
}

func (s *service) ObtenerDetalle(ctx context.Context, codigo string) (*CodigoSAP, error) {
	err := validarCodigo(codigo)
	if err != nil {
		return nil, fmt.Errorf("service: %w", err)
	}
	return s.repo.ObtenerDetalle(ctx, codigo)
}

func (s *service) BuscarPorDescripcion(ctx context.Context, parametro string) ([]CodigoSAP, error) {
	if strings.TrimSpace(parametro) == "" {
		return nil, fmt.Errorf("service: %w, parametro de busqueda vacío", errParametrosNoValidos)
	}
	return s.repo.BuscarPorDescripcion(ctx, parametro)
}

func (s *service) ModificarDescripcion(ctx context.Context, codigo string, nuevaDescripcion string) error {
	if strings.TrimSpace(codigo) == "" {
		return fmt.Errorf("service: %w, código vacío", errParametrosNoValidos)
	}
	if strings.TrimSpace(nuevaDescripcion) == "" {
		return fmt.Errorf("service: %w, descripción vacía", errParametrosNoValidos)
	}
	return s.repo.ModificarDescripcion(ctx, codigo, nuevaDescripcion)
}

func (s *service) Listar(ctx context.Context, limit, offset int) ([]CodigoSAP, error) {
	if limit <= 0 || offset < 0 {
		return nil, fmt.Errorf("service: %w, parámetros de lista inválidos", errParametrosNoValidos)
	}
	return s.repo.Listar(ctx, limit, offset)
}

func (s *service) PendientesRelacionados(ctx context.Context, documentoID string) ([]pendiente, error) {
	if strings.TrimSpace(documentoID) == "" {
		return nil, fmt.Errorf("service: %w, parametro de busqueda vacío", errParametrosNoValidos)
	}
	return s.repo.PendientesRelacionados(ctx, documentoID)
}
