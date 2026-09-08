package service

import (
	"context"
	"fmt"
	"strings"

	"task-tracker-go/internal/pendiente/dto"
	"task-tracker-go/internal/pendiente/model"
)

type repoAdjuntoInterface interface {
	CargarAdjunto(ctx context.Context, nuevoAdjunto *model.Adjunto) (int, error)
	EliminarAdjunto(ctx context.Context, id int) error
	ObtenerListaAdjuntos(ctx context.Context, pendienteID int) ([]dto.AdjuntoDetalleRespuesta, error)
}

type ServiceAdjunto struct {
	repo repoAdjuntoInterface
}

func NewServiceAdjunto(repo repoAdjuntoInterface) *ServiceAdjunto {
	return &ServiceAdjunto{repo: repo}
}

func (s *ServiceAdjunto) CargarAdjunto(ctx context.Context, nuevoAdjunto *dto.AdjuntoRequest) error {

	if nuevoAdjunto == nil {
		return fmt.Errorf("service: %w", model.ErrModelDeCargaNil)
	}

	if nuevoAdjunto.PendienteID <= 0 {
		return fmt.Errorf("%w, parametro: pendiente_id", model.ErrParametrosNoValidos)
	}
	if strings.TrimSpace(nuevoAdjunto.Descripcion) == "" {
		return fmt.Errorf("%w, parametro: descripcion", model.ErrParametrosNoValidos)
	}
	if strings.TrimSpace(nuevoAdjunto.ArchivoPath) == "" {
		return fmt.Errorf("%w, parametro: archivo_path", model.ErrParametrosNoValidos)
	}

	_, err := s.repo.CargarAdjunto(ctx, &model.Adjunto{
		PendienteID: nuevoAdjunto.PendienteID,
		Descripcion: nuevoAdjunto.Descripcion,
		ArchivoPath: nuevoAdjunto.ArchivoPath,
	})
	return err
}

func (s *ServiceAdjunto) EliminarAdjunto(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("service: %w, ID de adjunto inválido", model.ErrParametrosNoValidos)
	}
	return s.repo.EliminarAdjunto(ctx, id)
}

func (s *ServiceAdjunto) ObtenerListaAdjuntos(ctx context.Context, pendienteID int) ([]dto.AdjuntoDetalleRespuesta, error) {
	if pendienteID <= 0 {
		return nil, fmt.Errorf("service: %w, pendienteID inválido", model.ErrParametrosNoValidos)
	}
	return s.repo.ObtenerListaAdjuntos(ctx, pendienteID)
}
