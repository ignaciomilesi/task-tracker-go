package services

import (
	"context"
	"strings"
	"time"

	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"
)

// Interface para el acceso a datos de avances. Se usa como contrato para el servicio
// Los métodos reciben context.Context para permitir trazabilidad/cancelación desde capas superiores.
type avanceManagerDbInterface interface {
	// Parámetros:
	//      - Contexto
	//      - Avance a crear (models.Avance)
	// Salida:
	//      - Id creado (int)
	// Errores que puede devuelve:
	//      - appErrors.ParametroDeCargaVacio
	//      - appErrors.PendienteNoEncontrado
	Cargar(context.Context, *models.Avance) (int, error)

	// Parámetros:
	//      - Contexto
	//      - Id del avance (int)
	// Salida:
	//      - *models.Avance
	// Errores que puede devuelve:
	//      - appErrors.AvanceNoEncontrado
	ObtenerDetalle(context.Context, int) (*models.Avance, error)

	// Parámetros:
	//      - Contexto
	//      - Id del avance (int)
	// Salida:
	//      - error
	// Errores que puede devuelve:
	//      - appErrors.AvanceNoEncontrado
	Eliminar(context.Context, int) error

	// Parámetros:
	//      - Contexto
	//      - Id del pendiente (int)
	// Salida:
	//      - []models.Avance
	// Errores que puede devuelve:
	//      - appErrors.ParametroDeBusquedaVacio
	FiltrarPorPendiente(context.Context, int) ([]models.Avance, error)
}

// Service
type avanceService struct {
	repo avanceManagerDbInterface
}

// NewAvanceService crea el servicio con la implementación del repo que cumpla la interfaz
func NewAvanceService(repo avanceManagerDbInterface) *avanceService {
	return &avanceService{repo: repo}
}

// Cargar crea un nuevo avance validando parámetros mínimos
func (s *avanceService) Cargar(ctx context.Context, pendienteID int, descripcion string, fecha time.Time, mailPath string) (int, error) {
	if pendienteID <= 0 {
		return 0, appErrors.ParametroDeCargaVacio
	}
	if strings.TrimSpace(descripcion) == "" {
		return 0, appErrors.ParametroDeCargaVacio
	}
	if fecha.After(time.Now()) {
		return 0, appErrors.FechaNoValida
	}

	a := &models.Avance{PendienteID: pendienteID, Descripcion: descripcion, Fecha: fecha}
	if strings.TrimSpace(mailPath) != "" {
		a.MailPath = &mailPath
	}

	return s.repo.Cargar(ctx, a)
}

// ObtenerDetalle devuelve el detalle correspondiente al avance provisto
func (s *avanceService) ObtenerDetalle(ctx context.Context, id int) (*models.Avance, error) {
	if id <= 0 {
		return nil, appErrors.ParametroDeBusquedaVacio
	}
	return s.repo.ObtenerDetalle(ctx, id)
}

// Eliminar elimina un avance por id
func (s *avanceService) Eliminar(ctx context.Context, id int) error {
	if id <= 0 {
		return appErrors.ParametroDeBusquedaVacio
	}
	return s.repo.Eliminar(ctx, id)
}

// FiltrarPorPendiente devuelve avances asociados a un pendiente
func (s *avanceService) FiltrarPorPendiente(ctx context.Context, pendienteID int) ([]models.Avance, error) {
	if pendienteID <= 0 {
		return nil, appErrors.ParametroDeBusquedaVacio
	}
	return s.repo.FiltrarPorPendiente(ctx, pendienteID)
}
