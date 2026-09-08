package codigoID

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type repoInterface interface {
	Cargar(context.Context, *CodigoID) error
	ObtenerDetalle(context.Context, string) (*CodigoID, error)
	FiltrarPorEstado(context.Context, string) ([]CodigoID, error)
	ActualizarEstado(context.Context, string, string, time.Time) error
	Listar(context.Context, int, int) ([]CodigoID, error)
}

type service struct {
	repo repoInterface
}

func newService(repo repoInterface) *service {
	return &service{repo: repo}
}

// funciona como DTO
type nuevoCodigoIDRequest struct {
	Codigo      string
	Descripcion *string
	FechaPedido time.Time
}

func (s *service) Cargar(ctx context.Context, nuevaSolicitudCodigo *nuevoCodigoIDRequest) error {
	if nuevaSolicitudCodigo == nil {
		return fmt.Errorf("service: %w", errModelDeCargaNil)
	}

	nuevoCodigo, err := NewCodigoID(
		nuevaSolicitudCodigo.Codigo,
		nuevaSolicitudCodigo.Descripcion,
		"solicitado",
		nuevaSolicitudCodigo.FechaPedido,
		nil,
	)

	if err != nil {
		return fmt.Errorf("service: %w", errModelDeCargaNil)
	}

	return s.repo.Cargar(ctx, nuevoCodigo)
}

func (s *service) ObtenerDetalle(ctx context.Context, codigo string) (*CodigoID, error) {
	err := validarCodigo(codigo)
	if err != nil {
		return nil, fmt.Errorf("service: %w", err)
	}
	return s.repo.ObtenerDetalle(ctx, codigo)
}

func (s *service) ActualizarEstado(ctx context.Context, codigo string, nuevoEstado string, fecha time.Time) error {
	err := validarCodigo(codigo)
	if err != nil {
		return fmt.Errorf("service: %w", err)
	}
	if strings.TrimSpace(nuevoEstado) == "" {
		return fmt.Errorf("service: %w, estado vacío", errParametrosNoValidos)
	}
	if fecha.After(time.Now()) {
		return fmt.Errorf("service: %w, fecha no válida", errParametrosNoValidos)
	}
	return s.repo.ActualizarEstado(ctx, codigo, nuevoEstado, fecha)
}

func (s *service) FiltrarPorEstado(ctx context.Context, estado string) ([]CodigoID, error) {
	if strings.TrimSpace(estado) == "" {
		return nil, fmt.Errorf("service: %w, estado vacío", errParametrosNoValidos)
	}
	return s.repo.FiltrarPorEstado(ctx, estado)
}

func (s *service) Listar(ctx context.Context, limit, offset int) ([]CodigoID, error) {
	if limit <= 0 || offset < 0 {
		return nil, fmt.Errorf("service: %w, parámetros de lista inválidos", errParametrosNoValidos)
	}
	return s.repo.Listar(ctx, limit, offset)
}
