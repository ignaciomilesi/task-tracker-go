package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"task-tracker-go/internal/pendiente/dto"
	"task-tracker-go/internal/pendiente/model"
)

type repoAvanceInterface interface {
	CargarAvance(ctx context.Context, nuevoavance *model.Avance) (int, error)
	EliminarAvance(ctx context.Context, id int) error
	ObtenerListaAvances(ctx context.Context, pendienteID int) ([]dto.AvanceDetalleRespuesta, error)
}

type ServiceAvance struct {
	repo repoAvanceInterface
}

func NewServiceAvance(repo repoAvanceInterface) *ServiceAvance {
	return &ServiceAvance{repo: repo}
}

func (s *ServiceAvance) CargarAvance(ctx context.Context, nuevoAvance *dto.AvanceRequest) error {

	if nuevoAvance == nil {
		return fmt.Errorf("service: %w", model.ErrModelDeCargaNil)
	}

	if nuevoAvance.PendienteID <= 0 {
		return fmt.Errorf("service: %w, parametro: pendiente_id", model.ErrParametrosNoValidos)
	}
	if strings.TrimSpace(nuevoAvance.Descripcion) == "" {
		return fmt.Errorf("service: %w, parametro: descripcion", model.ErrParametrosNoValidos)
	}
	if nuevoAvance.Fecha.After(time.Now()) {
		return fmt.Errorf("service: %w, parametro: fecha", model.ErrParametrosNoValidos)
	}
	if nuevoAvance.MailPath != nil && strings.TrimSpace(*nuevoAvance.MailPath) == "" {
		nuevoAvance.MailPath = nil
	}

	_, err := s.repo.CargarAvance(ctx, &model.Avance{
		PendienteID: nuevoAvance.PendienteID,
		Descripcion: nuevoAvance.Descripcion,
		Fecha:       nuevoAvance.Fecha,
		MailPath:    nuevoAvance.MailPath,
	})
	return err
}

func (s *ServiceAvance) EliminarAvance(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("service: %w, ID del avance inválido", model.ErrParametrosNoValidos)
	}
	return s.repo.EliminarAvance(ctx, id)
}

func (s *ServiceAvance) ObtenerListaAvances(ctx context.Context, pendienteID int) ([]dto.AvanceDetalleRespuesta, error) {
	if pendienteID <= 0 {
		return nil, fmt.Errorf("service: %w, pendienteID inválido", model.ErrParametrosNoValidos)
	}
	return s.repo.ObtenerListaAvances(ctx, pendienteID)
}
