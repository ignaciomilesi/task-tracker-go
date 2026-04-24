package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"
)

type codigoSAPManagerDbMock struct {
	CargarFunc               func(context.Context, *models.CodigoSAP) error
	ObtenerDetalleFunc       func(context.Context, string) (*models.CodigoSAP, error)
	BuscarPorDescripcionFunc func(context.Context, string) ([]models.CodigoSAP, error)
	ModificarDescripcionFunc func(context.Context, string, string) error
	ListarFunc               func(context.Context, int, int) ([]models.CodigoSAP, error)
}

func (m *codigoSAPManagerDbMock) Cargar(ctx context.Context, c *models.CodigoSAP) error {
	if m.CargarFunc == nil {
		return fmt.Errorf("CargarFunc no implementado")
	}
	return m.CargarFunc(ctx, c)
}

func (m *codigoSAPManagerDbMock) ObtenerDetalle(ctx context.Context, codigo string) (*models.CodigoSAP, error) {
	if m.ObtenerDetalleFunc == nil {
		return nil, fmt.Errorf("ObtenerDetalleFunc no implementado")
	}
	return m.ObtenerDetalleFunc(ctx, codigo)
}

func (m *codigoSAPManagerDbMock) BuscarPorDescripcion(ctx context.Context, parametro string) ([]models.CodigoSAP, error) {
	if m.BuscarPorDescripcionFunc == nil {
		return nil, fmt.Errorf("BuscarPorDescripcionFunc no implementado")
	}
	return m.BuscarPorDescripcionFunc(ctx, parametro)
}

func (m *codigoSAPManagerDbMock) ModificarDescripcion(ctx context.Context, codigo string, nuevaDescripcion string) error {
	if m.ModificarDescripcionFunc == nil {
		return fmt.Errorf("ModificarDescripcionFunc no implementado")
	}
	return m.ModificarDescripcionFunc(ctx, codigo, nuevaDescripcion)
}

func (m *codigoSAPManagerDbMock) Listar(ctx context.Context, limit, offset int) ([]models.CodigoSAP, error) {
	if m.ListarFunc == nil {
		return nil, fmt.Errorf("ListarFunc no implementado")
	}
	return m.ListarFunc(ctx, limit, offset)
}

func TestCargarCodigoSAP(t *testing.T) {
	tests := []struct {
		name        string
		mockSetup   func() *codigoSAPManagerDbMock
		codigo      string
		descripcion string
		wantErr     error
	}{
		{name: "Todo Ok",
			mockSetup: func() *codigoSAPManagerDbMock {
				return &codigoSAPManagerDbMock{
					CargarFunc: func(ctx context.Context, c *models.CodigoSAP) error {
						return nil
					},
				}
			},
			codigo:      "ABC123",
			descripcion: "Desc",
			wantErr:     nil,
		},
		{name: "Codigo vacio",
			mockSetup: func() *codigoSAPManagerDbMock {
				return &codigoSAPManagerDbMock{
					CargarFunc: func(ctx context.Context, c *models.CodigoSAP) error { return nil },
				}
			},
			codigo:      "  ",
			descripcion: "Desc",
			wantErr:     appErrors.ParametroDeCargaVacio,
		},
		{name: "Codigo duplicado (repo)",
			mockSetup: func() *codigoSAPManagerDbMock {
				return &codigoSAPManagerDbMock{
					CargarFunc: func(ctx context.Context, c *models.CodigoSAP) error { return appErrors.CodigoSAPDuplicado },
				}
			},
			codigo:      "DUP",
			descripcion: "",
			wantErr:     appErrors.CodigoSAPDuplicado,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewCodigoSAPService(tt.mockSetup())
			ctx := t.Context()

			err := svc.Cargar(ctx, tt.codigo, tt.descripcion)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}

func TestObtenerDetalleCodigoSAP(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *codigoSAPManagerDbMock
		codigo    string
		wantErr   error
	}{
		{name: "Todo Ok",
			mockSetup: func() *codigoSAPManagerDbMock {
				return &codigoSAPManagerDbMock{
					ObtenerDetalleFunc: func(ctx context.Context, codigo string) (*models.CodigoSAP, error) {
						d := "desc"
						return &models.CodigoSAP{Codigo: codigo, Descripcion: &d}, nil
					},
				}
			},
			codigo:  "X1",
			wantErr: nil,
		},
		{name: "Codigo vacio",
			mockSetup: func() *codigoSAPManagerDbMock {
				return &codigoSAPManagerDbMock{
					ObtenerDetalleFunc: func(ctx context.Context, codigo string) (*models.CodigoSAP, error) { return nil, nil },
				}
			},
			codigo:  "",
			wantErr: appErrors.ParametroDeBusquedaVacio,
		},
		{name: "No encontrado",
			mockSetup: func() *codigoSAPManagerDbMock {
				return &codigoSAPManagerDbMock{
					ObtenerDetalleFunc: func(ctx context.Context, codigo string) (*models.CodigoSAP, error) {
						return nil, appErrors.CodigoSAPNoEncontrado
					},
				}
			},
			codigo:  "NOEX",
			wantErr: appErrors.CodigoSAPNoEncontrado,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewCodigoSAPService(tt.mockSetup())
			ctx := t.Context()

			_, err := svc.ObtenerDetalle(ctx, tt.codigo)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}

func TestBuscarPorDescripcionCodigoSAP(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *codigoSAPManagerDbMock
		parametro string
		wantErr   error
	}{
		{name: "Todo Ok",
			mockSetup: func() *codigoSAPManagerDbMock {
				return &codigoSAPManagerDbMock{
					BuscarPorDescripcionFunc: func(ctx context.Context, parametro string) ([]models.CodigoSAP, error) {
						return []models.CodigoSAP{{Codigo: "C1"}}, nil
					},
				}
			},
			parametro: "desc",
			wantErr:   nil,
		},
		{name: "Parametro vacio",
			mockSetup: func() *codigoSAPManagerDbMock {
				return &codigoSAPManagerDbMock{
					BuscarPorDescripcionFunc: func(ctx context.Context, parametro string) ([]models.CodigoSAP, error) { return nil, nil },
				}
			},
			parametro: " ",
			wantErr:   appErrors.ParametroDeBusquedaVacio,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewCodigoSAPService(tt.mockSetup())
			ctx := t.Context()

			_, err := svc.BuscarPorDescripcion(ctx, tt.parametro)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}

func TestModificarDescripcionCodigoSAP(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *codigoSAPManagerDbMock
		codigo    string
		nuevaDesc string
		wantErr   error
	}{
		{name: "Todo Ok",
			mockSetup: func() *codigoSAPManagerDbMock {
				return &codigoSAPManagerDbMock{
					ModificarDescripcionFunc: func(ctx context.Context, codigo string, nueva string) error { return nil },
				}
			},
			codigo:    "C1",
			nuevaDesc: "Nueva",
			wantErr:   nil,
		},
		{name: "Codigo vacio",
			mockSetup: func() *codigoSAPManagerDbMock {
				return &codigoSAPManagerDbMock{
					ModificarDescripcionFunc: func(ctx context.Context, codigo string, nueva string) error { return nil },
				}
			},
			codigo:    "",
			nuevaDesc: "Nueva",
			wantErr:   appErrors.ParametroDeBusquedaVacio,
		},
		{name: "Descripcion vacia",
			mockSetup: func() *codigoSAPManagerDbMock {
				return &codigoSAPManagerDbMock{
					ModificarDescripcionFunc: func(ctx context.Context, codigo string, nueva string) error { return nil },
				}
			},
			codigo:    "C1",
			nuevaDesc: "  ",
			wantErr:   appErrors.ParametroDeCargaVacio,
		},
		{name: "No encontrado (repo)",
			mockSetup: func() *codigoSAPManagerDbMock {
				return &codigoSAPManagerDbMock{
					ModificarDescripcionFunc: func(ctx context.Context, codigo string, nueva string) error { return appErrors.CodigoSAPNoEncontrado },
				}
			},
			codigo:    "NOEX",
			nuevaDesc: "Nueva",
			wantErr:   appErrors.CodigoSAPNoEncontrado,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewCodigoSAPService(tt.mockSetup())
			ctx := t.Context()

			err := svc.ModificarDescripcion(ctx, tt.codigo, tt.nuevaDesc)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}

func TestListarCodigoSAP(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *codigoSAPManagerDbMock
		limit     int
		offset    int
		wantErr   error
	}{
		{name: "Todo Ok",
			mockSetup: func() *codigoSAPManagerDbMock {
				return &codigoSAPManagerDbMock{
					ListarFunc: func(ctx context.Context, limit, offset int) ([]models.CodigoSAP, error) {
						return []models.CodigoSAP{{Codigo: "C1"}}, nil
					},
				}
			},
			limit:   10,
			offset:  0,
			wantErr: nil,
		},
		{name: "Limit invalido",
			mockSetup: func() *codigoSAPManagerDbMock {
				return &codigoSAPManagerDbMock{
					ListarFunc: func(ctx context.Context, limit, offset int) ([]models.CodigoSAP, error) { return nil, nil },
				}
			},
			limit:   0,
			offset:  0,
			wantErr: appErrors.ParametrosDeListaInvalidos,
		},
		{name: "Offset invalido",
			mockSetup: func() *codigoSAPManagerDbMock {
				return &codigoSAPManagerDbMock{
					ListarFunc: func(ctx context.Context, limit, offset int) ([]models.CodigoSAP, error) { return nil, nil },
				}
			},
			limit:   10,
			offset:  -1,
			wantErr: appErrors.ParametrosDeListaInvalidos,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewCodigoSAPService(tt.mockSetup())
			ctx := t.Context()

			_, err := svc.Listar(ctx, tt.limit, tt.offset)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}
