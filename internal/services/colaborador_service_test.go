package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"
)

type colaboradorManagerDbMock struct {
	CrearFunc              func(context.Context, *models.Colaborador) (int, error)
	ObtenerIDPorNombreFunc func(context.Context, string) (int, error)
	BuscarFunc             func(context.Context, string) ([]models.Colaborador, error)
	ListarFunc             func(context.Context, int, int) ([]models.Colaborador, error)
}

func (m *colaboradorManagerDbMock) Crear(ctx context.Context, c *models.Colaborador) (int, error) {
	if m.CrearFunc == nil {
		return 0, fmt.Errorf("CrearFunc no implementado")
	}
	return m.CrearFunc(ctx, c)
}

func (m *colaboradorManagerDbMock) ObtenerIDPorNombre(ctx context.Context, nombre string) (int, error) {
	if m.ObtenerIDPorNombreFunc == nil {
		return 0, fmt.Errorf("ObtenerIDPorNombreFunc no implementado")
	}
	return m.ObtenerIDPorNombreFunc(ctx, nombre)
}

func (m *colaboradorManagerDbMock) Buscar(ctx context.Context, parametro string) ([]models.Colaborador, error) {
	if m.BuscarFunc == nil {
		return nil, fmt.Errorf("BuscarFunc no implementado")
	}
	return m.BuscarFunc(ctx, parametro)
}

func (m *colaboradorManagerDbMock) Listar(ctx context.Context, limit, offset int) ([]models.Colaborador, error) {
	if m.ListarFunc == nil {
		return nil, fmt.Errorf("ListarFunc no implementado")
	}
	return m.ListarFunc(ctx, limit, offset)
}

func TestCrearColaborador(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *colaboradorManagerDbMock
		nombre    string
		wantErr   error
	}{
		{name: "Todo Ok",
			mockSetup: func() *colaboradorManagerDbMock {
				return &colaboradorManagerDbMock{
					CrearFunc: func(ctx context.Context, c *models.Colaborador) (int, error) {
						return 1, nil
					},
				}
			},
			nombre:  "Pedro",
			wantErr: nil,
		},
		{name: "Nombre en blanco",
			mockSetup: func() *colaboradorManagerDbMock {
				return &colaboradorManagerDbMock{
					CrearFunc: func(ctx context.Context, c *models.Colaborador) (int, error) {
						return 0, nil
					},
				}
			},
			nombre:  "",
			wantErr: appErrors.ParametroDeCargaVacio,
		},
		{name: "Colaborador duplicado",
			mockSetup: func() *colaboradorManagerDbMock {
				return &colaboradorManagerDbMock{
					CrearFunc: func(ctx context.Context, c *models.Colaborador) (int, error) {
						return 0, appErrors.ColaboradorDuplicado
					},
				}
			},
			nombre:  "Lucia",
			wantErr: appErrors.ColaboradorDuplicado,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewColaboradorService(tt.mockSetup())
			ctx := t.Context()

			_, err := svc.Crear(ctx, tt.nombre)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}

func TestObtenerIDPorNombreColaborador(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *colaboradorManagerDbMock
		nombre    string
		wantErr   error
	}{
		{name: "Todo Ok",
			mockSetup: func() *colaboradorManagerDbMock {
				return &colaboradorManagerDbMock{
					ObtenerIDPorNombreFunc: func(ctx context.Context, nombre string) (int, error) {
						return 3, nil
					},
				}
			},
			nombre:  "Marta",
			wantErr: nil,
		},
		{name: "Nombre vacio",
			mockSetup: func() *colaboradorManagerDbMock {
				return &colaboradorManagerDbMock{
					ObtenerIDPorNombreFunc: func(ctx context.Context, nombre string) (int, error) {
						return 0, nil
					},
				}
			},
			nombre:  "   ",
			wantErr: appErrors.ParametroDeBusquedaVacio,
		},
		{name: "No encontrado",
			mockSetup: func() *colaboradorManagerDbMock {
				return &colaboradorManagerDbMock{
					ObtenerIDPorNombreFunc: func(ctx context.Context, nombre string) (int, error) {
						return 0, appErrors.ColaboradorNoEncontrado
					},
				}
			},
			nombre:  "NoExiste",
			wantErr: appErrors.ColaboradorNoEncontrado,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewColaboradorService(tt.mockSetup())
			ctx := t.Context()

			_, err := svc.ObtenerIDPorNombre(ctx, tt.nombre)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}

func TestBuscarColaborador(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *colaboradorManagerDbMock
		parametro string
		wantErr   error
	}{
		{name: "Todo Ok",
			mockSetup: func() *colaboradorManagerDbMock {
				return &colaboradorManagerDbMock{
					BuscarFunc: func(ctx context.Context, parametro string) ([]models.Colaborador, error) {
						return []models.Colaborador{{ID: 1, Nombre: "Ana"}}, nil
					},
				}
			},
			parametro: "An",
			wantErr:   nil,
		},
		{name: "Parametro vacio",
			mockSetup: func() *colaboradorManagerDbMock {
				return &colaboradorManagerDbMock{
					BuscarFunc: func(ctx context.Context, parametro string) ([]models.Colaborador, error) {
						return nil, nil
					},
				}
			},
			parametro: "",
			wantErr:   appErrors.ParametroDeBusquedaVacio,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewColaboradorService(tt.mockSetup())
			ctx := t.Context()

			_, err := svc.Buscar(ctx, tt.parametro)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}

func TestListarColaboradores(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *colaboradorManagerDbMock
		limit     int
		offset    int
		wantErr   error
	}{
		{name: "Todo Ok",
			mockSetup: func() *colaboradorManagerDbMock {
				return &colaboradorManagerDbMock{
					ListarFunc: func(ctx context.Context, limit, offset int) ([]models.Colaborador, error) {
						return []models.Colaborador{{ID: 1, Nombre: "Ana"}}, nil
					},
				}
			},
			limit:   10,
			offset:  0,
			wantErr: nil,
		},
		{name: "Limit invalido",
			mockSetup: func() *colaboradorManagerDbMock {
				return &colaboradorManagerDbMock{
					ListarFunc: func(ctx context.Context, limit, offset int) ([]models.Colaborador, error) {
						return nil, nil
					},
				}
			},
			limit:   0,
			offset:  0,
			wantErr: appErrors.ParametrosDeListaInvalidos,
		},
		{name: "Offset invalido",
			mockSetup: func() *colaboradorManagerDbMock {
				return &colaboradorManagerDbMock{
					ListarFunc: func(ctx context.Context, limit, offset int) ([]models.Colaborador, error) {
						return nil, nil
					},
				}
			},
			limit:   10,
			offset:  -1,
			wantErr: appErrors.ParametrosDeListaInvalidos,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewColaboradorService(tt.mockSetup())
			ctx := t.Context()

			_, err := svc.Listar(ctx, tt.limit, tt.offset)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}
