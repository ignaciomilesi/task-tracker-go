package services

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"
)

type avanceManagerDbMock struct {
	CargarFunc              func(context.Context, *models.Avance) (int, error)
	ObtenerDetalleFunc      func(context.Context, int) (*models.Avance, error)
	EliminarFunc            func(context.Context, int) error
	FiltrarPorPendienteFunc func(context.Context, int) ([]models.Avance, error)
}

func (m *avanceManagerDbMock) Cargar(ctx context.Context, a *models.Avance) (int, error) {
	if m.CargarFunc == nil {
		return 0, fmt.Errorf("CargarFunc no implementado")
	}
	return m.CargarFunc(ctx, a)
}

func (m *avanceManagerDbMock) ObtenerDetalle(ctx context.Context, id int) (*models.Avance, error) {
	if m.ObtenerDetalleFunc == nil {
		return nil, fmt.Errorf("ObtenerDetalleFunc no implementado")
	}
	return m.ObtenerDetalleFunc(ctx, id)
}

func (m *avanceManagerDbMock) Eliminar(ctx context.Context, id int) error {
	if m.EliminarFunc == nil {
		return fmt.Errorf("EliminarFunc no implementado")
	}
	return m.EliminarFunc(ctx, id)
}

func (m *avanceManagerDbMock) FiltrarPorPendiente(ctx context.Context, pendienteID int) ([]models.Avance, error) {
	if m.FiltrarPorPendienteFunc == nil {
		return nil, fmt.Errorf("FiltrarPorPendienteFunc no implementado")
	}
	return m.FiltrarPorPendienteFunc(ctx, pendienteID)
}

func TestCargarAvance(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name        string
		mockSetup   func() *avanceManagerDbMock
		pendienteID int
		descripcion string
		fecha       time.Time
		mailPath    string
		wantErr     error
	}{
		{name: "Todo Ok",
			mockSetup: func() *avanceManagerDbMock {
				return &avanceManagerDbMock{CargarFunc: func(ctx context.Context, a *models.Avance) (int, error) { return 1, nil }}
			},
			pendienteID: 1,
			descripcion: "d",
			fecha:       now.Add(-time.Hour),
			mailPath:    "/m",
			wantErr:     nil,
		},
		{name: "Pendiente invalido",
			mockSetup: func() *avanceManagerDbMock {
				return &avanceManagerDbMock{CargarFunc: func(ctx context.Context, a *models.Avance) (int, error) { return 0, appErrors.PendienteNoEncontrado }}
			},
			pendienteID: 0,
			descripcion: "d",
			fecha:       now,
			mailPath:    "",
			wantErr:     appErrors.ParametroDeCargaVacio,
		},
		{name: "Descripcion vacia",
			mockSetup: func() *avanceManagerDbMock {
				return &avanceManagerDbMock{CargarFunc: func(ctx context.Context, a *models.Avance) (int, error) { return 0, nil }}
			},
			pendienteID: 1,
			descripcion: " ",
			fecha:       now,
			mailPath:    "",
			wantErr:     appErrors.ParametroDeCargaVacio,
		},
		{name: "Fecha futura",
			mockSetup: func() *avanceManagerDbMock {
				return &avanceManagerDbMock{CargarFunc: func(ctx context.Context, a *models.Avance) (int, error) { return 0, nil }}
			},
			pendienteID: 1,
			descripcion: "d",
			fecha:       now.Add(48 * time.Hour),
			mailPath:    "",
			wantErr:     appErrors.FechaNoValida,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAvanceService(tt.mockSetup())
			ctx := t.Context()

			_, err := svc.Cargar(ctx, tt.pendienteID, tt.descripcion, tt.fecha, tt.mailPath)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}

func TestObtenerEliminarFiltrarAvance(t *testing.T) {
	tests := []struct {
		name        string
		mockSetup   func() *avanceManagerDbMock
		metodo      string // obtener | eliminar | filtrar
		id          int
		pendienteID int
		wantErr     error
	}{
		{name: "Obtener OK",
			mockSetup: func() *avanceManagerDbMock {
				return &avanceManagerDbMock{ObtenerDetalleFunc: func(ctx context.Context, id int) (*models.Avance, error) { return &models.Avance{ID: id}, nil }}
			},
			metodo:  "obtener",
			id:      1,
			wantErr: nil,
		},
		{name: "Obtener id invalido",
			mockSetup: func() *avanceManagerDbMock {
				return &avanceManagerDbMock{ObtenerDetalleFunc: func(ctx context.Context, id int) (*models.Avance, error) { return nil, nil }}
			},
			metodo:  "obtener",
			id:      0,
			wantErr: appErrors.ParametroDeBusquedaVacio,
		},
		{name: "Eliminar OK",
			mockSetup: func() *avanceManagerDbMock {
				return &avanceManagerDbMock{EliminarFunc: func(ctx context.Context, id int) error { return nil }}
			},
			metodo:  "eliminar",
			id:      1,
			wantErr: nil,
		},
		{name: "Eliminar id invalido",
			mockSetup: func() *avanceManagerDbMock {
				return &avanceManagerDbMock{EliminarFunc: func(ctx context.Context, id int) error { return nil }}
			},
			metodo:  "eliminar",
			id:      0,
			wantErr: appErrors.ParametroDeBusquedaVacio,
		},
		{name: "Filtrar OK",
			mockSetup: func() *avanceManagerDbMock {
				return &avanceManagerDbMock{FiltrarPorPendienteFunc: func(ctx context.Context, pendienteID int) ([]models.Avance, error) {
					return []models.Avance{{ID: 1}}, nil
				}}
			},
			metodo:      "filtrar",
			pendienteID: 1,
			wantErr:     nil,
		},
		{name: "Filtrar id invalido",
			mockSetup: func() *avanceManagerDbMock {
				return &avanceManagerDbMock{FiltrarPorPendienteFunc: func(ctx context.Context, pendienteID int) ([]models.Avance, error) { return nil, nil }}
			},
			metodo:      "filtrar",
			pendienteID: 0,
			wantErr:     appErrors.ParametroDeBusquedaVacio,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAvanceService(tt.mockSetup())
			ctx := t.Context()

			var err error
			switch tt.metodo {
			case "obtener":
				_, err = svc.ObtenerDetalle(ctx, tt.id)
			case "eliminar":
				err = svc.Eliminar(ctx, tt.id)
			case "filtrar":
				_, err = svc.FiltrarPorPendiente(ctx, tt.pendienteID)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}
