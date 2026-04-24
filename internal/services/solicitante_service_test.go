package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"
)

type solicitanteManagerDbMock struct {
	CrearFunc              func(context.Context, *models.Solicitante) (int, error)
	ObtenerIDPorNombreFunc func(context.Context, string) (int, error)
	BuscarFunc             func(context.Context, string) ([]models.Solicitante, error)
	ListarFunc             func(context.Context, int, int) ([]models.Solicitante, error)
}

func (m *solicitanteManagerDbMock) Crear(ctx context.Context, s *models.Solicitante) (int, error) {
	if m.CrearFunc == nil {
		return 0, fmt.Errorf("CrearFunc no implementado")
	}
	return m.CrearFunc(ctx, s)
}

func (m *solicitanteManagerDbMock) ObtenerIDPorNombre(ctx context.Context, nombre string) (int, error) {
	if m.ObtenerIDPorNombreFunc == nil {
		return 0, fmt.Errorf("ObtenerIDPorNombreFunc no implementado")
	}
	return m.ObtenerIDPorNombreFunc(ctx, nombre)
}

func (m *solicitanteManagerDbMock) Buscar(ctx context.Context, parametro string) ([]models.Solicitante, error) {
	if m.BuscarFunc == nil {
		return nil, fmt.Errorf("BuscarFunc no implementado")
	}
	return m.BuscarFunc(ctx, parametro)
}

func (m *solicitanteManagerDbMock) Listar(ctx context.Context, limit, offset int) ([]models.Solicitante, error) {
	if m.ListarFunc == nil {
		return nil, fmt.Errorf("ListarFunc no implementado")
	}
	return m.ListarFunc(ctx, limit, offset)
}

func TestCrearSolicitante(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *solicitanteManagerDbMock
		nombre    string
		wantErr   error
	}{
		{name: "Todo Ok",
			mockSetup: func() *solicitanteManagerDbMock {
				return &solicitanteManagerDbMock{
					CrearFunc: func(ctx context.Context, s *models.Solicitante) (int, error) {
						return 1, nil
					},
				}
			},
			nombre:  "Juan",
			wantErr: nil,
		},
		{name: "Nombre en blanco",
			mockSetup: func() *solicitanteManagerDbMock {
				return &solicitanteManagerDbMock{
					CrearFunc: func(ctx context.Context, s *models.Solicitante) (int, error) {
						return 0, nil
					},
				}
			},
			nombre:  "   ",
			wantErr: appErrors.ParametroDeCargaVacio,
		},
		{name: "Solicitante duplicado",
			mockSetup: func() *solicitanteManagerDbMock {
				return &solicitanteManagerDbMock{
					CrearFunc: func(ctx context.Context, s *models.Solicitante) (int, error) {
						return 0, appErrors.SolicitanteDuplicado
					},
				}
			},
			nombre:  "Ana",
			wantErr: appErrors.SolicitanteDuplicado,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewSolicitanteService(tt.mockSetup())
			ctx := t.Context()

			_, err := svc.Crear(ctx, tt.nombre)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}

func TestObtenerIDPorNombre(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *solicitanteManagerDbMock
		nombre    string
		wantErr   error
	}{
		{name: "Todo Ok",
			mockSetup: func() *solicitanteManagerDbMock {
				return &solicitanteManagerDbMock{
					ObtenerIDPorNombreFunc: func(ctx context.Context, nombre string) (int, error) {
						return 2, nil
					},
				}
			},
			nombre:  "Marcos",
			wantErr: nil,
		},
		{name: "Nombre vacio",
			mockSetup: func() *solicitanteManagerDbMock {
				return &solicitanteManagerDbMock{
					ObtenerIDPorNombreFunc: func(ctx context.Context, nombre string) (int, error) {
						return 0, nil
					},
				}
			},
			nombre:  "",
			wantErr: appErrors.ParametroDeBusquedaVacio,
		},
		{name: "No encontrado",
			mockSetup: func() *solicitanteManagerDbMock {
				return &solicitanteManagerDbMock{
					ObtenerIDPorNombreFunc: func(ctx context.Context, nombre string) (int, error) {
						return 0, appErrors.SolicitanteNoEncontrado
					},
				}
			},
			nombre:  "NoExiste",
			wantErr: appErrors.SolicitanteNoEncontrado,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewSolicitanteService(tt.mockSetup())
			ctx := t.Context()

			_, err := svc.ObtenerIDPorNombre(ctx, tt.nombre)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}

func TestBuscarSolicitante(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *solicitanteManagerDbMock
		parametro string
		wantErr   error
	}{
		{name: "Todo Ok",
			mockSetup: func() *solicitanteManagerDbMock {
				return &solicitanteManagerDbMock{
					BuscarFunc: func(ctx context.Context, parametro string) ([]models.Solicitante, error) {
						return []models.Solicitante{{ID: 1, Nombre: "Ana"}}, nil
					},
				}
			},
			parametro: "An",
			wantErr:   nil,
		},
		{name: "Parametro vacio",
			mockSetup: func() *solicitanteManagerDbMock {
				return &solicitanteManagerDbMock{
					BuscarFunc: func(ctx context.Context, parametro string) ([]models.Solicitante, error) {
						return nil, nil
					},
				}
			},
			parametro: "  ",
			wantErr:   appErrors.ParametroDeBusquedaVacio,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewSolicitanteService(tt.mockSetup())
			ctx := t.Context()

			_, err := svc.Buscar(ctx, tt.parametro)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}

func TestListarSolicitantes(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *solicitanteManagerDbMock
		limit     int
		offset    int
		wantErr   error
	}{
		{name: "Todo Ok",
			mockSetup: func() *solicitanteManagerDbMock {
				return &solicitanteManagerDbMock{
					ListarFunc: func(ctx context.Context, limit, offset int) ([]models.Solicitante, error) {
						return []models.Solicitante{{ID: 1, Nombre: "Ana"}}, nil
					},
				}
			},
			limit:   10,
			offset:  0,
			wantErr: nil,
		},
		{name: "Limit invalido",
			mockSetup: func() *solicitanteManagerDbMock {
				return &solicitanteManagerDbMock{
					ListarFunc: func(ctx context.Context, limit, offset int) ([]models.Solicitante, error) {
						return nil, nil
					},
				}
			},
			limit:   0,
			offset:  0,
			wantErr: appErrors.ParametrosDeListaInvalidos,
		},
		{name: "Offset invalido",
			mockSetup: func() *solicitanteManagerDbMock {
				return &solicitanteManagerDbMock{
					ListarFunc: func(ctx context.Context, limit, offset int) ([]models.Solicitante, error) {
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
			svc := NewSolicitanteService(tt.mockSetup())
			ctx := t.Context()

			_, err := svc.Listar(ctx, tt.limit, tt.offset)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}
