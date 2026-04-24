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

type codigoIDManagerDbMock struct {
	CargarFunc           func(context.Context, *models.CodigoID) error
	ObtenerDetalleFunc   func(context.Context, string) (*models.CodigoID, error)
	FiltrarPorEstadoFunc func(context.Context, string) ([]models.CodigoID, error)
	ActualizarEstadoFunc func(context.Context, string, string, time.Time) error
	ListarFunc           func(context.Context, int, int) ([]models.CodigoID, error)
}

func (m *codigoIDManagerDbMock) Cargar(ctx context.Context, c *models.CodigoID) error {
	if m.CargarFunc == nil {
		return fmt.Errorf("CargarFunc no implementado")
	}
	return m.CargarFunc(ctx, c)
}

func (m *codigoIDManagerDbMock) ObtenerDetalle(ctx context.Context, codigo string) (*models.CodigoID, error) {
	if m.ObtenerDetalleFunc == nil {
		return nil, fmt.Errorf("ObtenerDetalleFunc no implementado")
	}
	return m.ObtenerDetalleFunc(ctx, codigo)
}

func (m *codigoIDManagerDbMock) FiltrarPorEstado(ctx context.Context, estado string) ([]models.CodigoID, error) {
	if m.FiltrarPorEstadoFunc == nil {
		return nil, fmt.Errorf("FiltrarPorEstadoFunc no implementado")
	}
	return m.FiltrarPorEstadoFunc(ctx, estado)
}

func (m *codigoIDManagerDbMock) ActualizarEstado(ctx context.Context, codigo string, nuevoEstado string, fecha time.Time) error {
	if m.ActualizarEstadoFunc == nil {
		return fmt.Errorf("ActualizarEstadoFunc no implementado")
	}
	return m.ActualizarEstadoFunc(ctx, codigo, nuevoEstado, fecha)
}

func (m *codigoIDManagerDbMock) Listar(ctx context.Context, limit, offset int) ([]models.CodigoID, error) {
	if m.ListarFunc == nil {
		return nil, fmt.Errorf("ListarFunc no implementado")
	}
	return m.ListarFunc(ctx, limit, offset)
}

func TestCargarCodigoID(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name        string
		mockSetup   func() *codigoIDManagerDbMock
		codigo      string
		descripcion string
		estado      string
		fechaPedido time.Time
		wantErr     error
	}{
		{name: "Todo Ok",
			mockSetup: func() *codigoIDManagerDbMock {
				return &codigoIDManagerDbMock{
					CargarFunc: func(ctx context.Context, c *models.CodigoID) error { return nil },
				}
			},
			codigo:      "ID1",
			descripcion: "desc",
			estado:      "PEND",
			fechaPedido: now.Add(-time.Hour),
			wantErr:     nil,
		},
		{name: "Codigo vacio",
			mockSetup: func() *codigoIDManagerDbMock {
				return &codigoIDManagerDbMock{CargarFunc: func(ctx context.Context, c *models.CodigoID) error { return nil }}
			},
			codigo:      " ",
			descripcion: "d",
			estado:      "PEND",
			fechaPedido: now,
			wantErr:     appErrors.ParametroDeCargaVacio,
		},
		{name: "Estado vacio",
			mockSetup: func() *codigoIDManagerDbMock {
				return &codigoIDManagerDbMock{CargarFunc: func(ctx context.Context, c *models.CodigoID) error { return nil }}
			},
			codigo:      "ID2",
			descripcion: "d",
			estado:      " ",
			fechaPedido: now,
			wantErr:     appErrors.ParametroDeCargaVacio,
		},
		{name: "Fecha futura invalida",
			mockSetup: func() *codigoIDManagerDbMock {
				return &codigoIDManagerDbMock{CargarFunc: func(ctx context.Context, c *models.CodigoID) error { return nil }}
			},
			codigo:      "ID3",
			descripcion: "d",
			estado:      "PEND",
			fechaPedido: now.Add(48 * time.Hour),
			wantErr:     appErrors.FechaNoValida,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewCodigoIDService(tt.mockSetup())
			ctx := t.Context()

			err := svc.Cargar(ctx, tt.codigo, tt.descripcion, tt.estado, tt.fechaPedido)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}

func TestObtenerDetalleCodigoID(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *codigoIDManagerDbMock
		codigo    string
		wantErr   error
	}{
		{name: "Todo Ok",
			mockSetup: func() *codigoIDManagerDbMock {
				return &codigoIDManagerDbMock{ObtenerDetalleFunc: func(ctx context.Context, codigo string) (*models.CodigoID, error) {
					return &models.CodigoID{Codigo: codigo}, nil
				}}
			},
			codigo:  "ID1",
			wantErr: nil,
		},
		{name: "Codigo vacio",
			mockSetup: func() *codigoIDManagerDbMock {
				return &codigoIDManagerDbMock{ObtenerDetalleFunc: func(ctx context.Context, codigo string) (*models.CodigoID, error) { return nil, nil }}
			},
			codigo:  "",
			wantErr: appErrors.ParametroDeBusquedaVacio,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewCodigoIDService(tt.mockSetup())
			ctx := t.Context()

			_, err := svc.ObtenerDetalle(ctx, tt.codigo)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}

func TestFiltrarPorEstadoCodigoID(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *codigoIDManagerDbMock
		estado    string
		wantErr   error
	}{
		{name: "Todo Ok",
			mockSetup: func() *codigoIDManagerDbMock {
				return &codigoIDManagerDbMock{FiltrarPorEstadoFunc: func(ctx context.Context, estado string) ([]models.CodigoID, error) {
					return []models.CodigoID{{Codigo: "ID1"}}, nil
				}}
			},
			estado:  "PEND",
			wantErr: nil,
		},
		{name: "Estado vacio",
			mockSetup: func() *codigoIDManagerDbMock {
				return &codigoIDManagerDbMock{FiltrarPorEstadoFunc: func(ctx context.Context, estado string) ([]models.CodigoID, error) { return nil, nil }}
			},
			estado:  " ",
			wantErr: appErrors.ParametroDeBusquedaVacio,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewCodigoIDService(tt.mockSetup())
			ctx := t.Context()

			_, err := svc.FiltrarPorEstado(ctx, tt.estado)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}

func TestActualizarEstadoCodigoID(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name        string
		mockSetup   func() *codigoIDManagerDbMock
		codigo      string
		nuevoEstado string
		fecha       time.Time
		wantErr     error
	}{
		{name: "Todo Ok",
			mockSetup: func() *codigoIDManagerDbMock {
				return &codigoIDManagerDbMock{ActualizarEstadoFunc: func(ctx context.Context, codigo string, nuevo string, fecha time.Time) error { return nil }}
			},
			codigo:      "ID1",
			nuevoEstado: "OK",
			fecha:       now,
			wantErr:     nil,
		},
		{name: "Codigo vacio",
			mockSetup: func() *codigoIDManagerDbMock {
				return &codigoIDManagerDbMock{ActualizarEstadoFunc: func(ctx context.Context, codigo string, nuevo string, fecha time.Time) error { return nil }}
			},
			codigo:      " ",
			nuevoEstado: "OK",
			fecha:       now,
			wantErr:     appErrors.ParametroDeBusquedaVacio,
		},
		{name: "Nuevo estado vacio",
			mockSetup: func() *codigoIDManagerDbMock {
				return &codigoIDManagerDbMock{ActualizarEstadoFunc: func(ctx context.Context, codigo string, nuevo string, fecha time.Time) error { return nil }}
			},
			codigo:      "ID1",
			nuevoEstado: " ",
			fecha:       now,
			wantErr:     appErrors.ParametroDeCargaVacio,
		},
		{name: "Fecha fuera de rango",
			mockSetup: func() *codigoIDManagerDbMock {
				return &codigoIDManagerDbMock{ActualizarEstadoFunc: func(ctx context.Context, codigo string, nuevo string, fecha time.Time) error { return nil }}
			},
			codigo:      "ID1",
			nuevoEstado: "OK",
			fecha:       now.Add(72 * time.Hour),
			wantErr:     appErrors.FechaNoValida,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewCodigoIDService(tt.mockSetup())
			ctx := t.Context()

			err := svc.ActualizarEstado(ctx, tt.codigo, tt.nuevoEstado, tt.fecha)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}

func TestListarCodigoID(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *codigoIDManagerDbMock
		limit     int
		offset    int
		wantErr   error
	}{
		{name: "Todo Ok",
			mockSetup: func() *codigoIDManagerDbMock {
				return &codigoIDManagerDbMock{ListarFunc: func(ctx context.Context, limit, offset int) ([]models.CodigoID, error) {
					return []models.CodigoID{{Codigo: "ID1"}}, nil
				}}
			},
			limit:   10,
			offset:  0,
			wantErr: nil,
		},
		{name: "Limit invalido",
			mockSetup: func() *codigoIDManagerDbMock {
				return &codigoIDManagerDbMock{ListarFunc: func(ctx context.Context, limit, offset int) ([]models.CodigoID, error) { return nil, nil }}
			},
			limit:   0,
			offset:  0,
			wantErr: appErrors.ParametrosDeListaInvalidos,
		},
		{name: "Offset invalido",
			mockSetup: func() *codigoIDManagerDbMock {
				return &codigoIDManagerDbMock{ListarFunc: func(ctx context.Context, limit, offset int) ([]models.CodigoID, error) { return nil, nil }}
			},
			limit:   10,
			offset:  -1,
			wantErr: appErrors.ParametrosDeListaInvalidos,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewCodigoIDService(tt.mockSetup())
			ctx := t.Context()

			_, err := svc.Listar(ctx, tt.limit, tt.offset)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}
