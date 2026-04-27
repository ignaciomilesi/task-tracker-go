package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"
)

type adjuntoManagerDbMock struct {
	CargarFunc                func(context.Context, *models.Adjunto) (int, error)
	ObtenerDetalleFunc        func(context.Context, int) (*models.Adjunto, error)
	EliminarFunc              func(context.Context, int) error
	FiltrarPorPendienteFunc   func(context.Context, int) ([]models.Adjunto, error)
	ActualizarPathFunc        func(context.Context, int, string) error
	ActualizarDescripcionFunc func(context.Context, int, string) error
}

func (m *adjuntoManagerDbMock) Cargar(ctx context.Context, a *models.Adjunto) (int, error) {
	if m.CargarFunc == nil {
		return 0, fmt.Errorf("CargarFunc no implementado")
	}
	return m.CargarFunc(ctx, a)
}

func (m *adjuntoManagerDbMock) ObtenerDetalle(ctx context.Context, id int) (*models.Adjunto, error) {
	if m.ObtenerDetalleFunc == nil {
		return nil, fmt.Errorf("ObtenerDetalleFunc no implementado")
	}
	return m.ObtenerDetalleFunc(ctx, id)
}

func (m *adjuntoManagerDbMock) Eliminar(ctx context.Context, id int) error {
	if m.EliminarFunc == nil {
		return fmt.Errorf("EliminarFunc no implementado")
	}
	return m.EliminarFunc(ctx, id)
}

func (m *adjuntoManagerDbMock) FiltrarPorPendiente(ctx context.Context, pendienteID int) ([]models.Adjunto, error) {
	if m.FiltrarPorPendienteFunc == nil {
		return nil, fmt.Errorf("FiltrarPorPendienteFunc no implementado")
	}
	return m.FiltrarPorPendienteFunc(ctx, pendienteID)
}

func (m *adjuntoManagerDbMock) ActualizarPath(ctx context.Context, id int, path string) error {
	if m.ActualizarPathFunc == nil {
		return fmt.Errorf("ActualizarPathFunc no implementado")
	}
	return m.ActualizarPathFunc(ctx, id, path)
}

func (m *adjuntoManagerDbMock) ActualizarDescripcion(ctx context.Context, id int, descripcion string) error {
	if m.ActualizarDescripcionFunc == nil {
		return fmt.Errorf("ActualizarDescripcionFunc no implementado")
	}
	return m.ActualizarDescripcionFunc(ctx, id, descripcion)
}

func TestCargarAdjunto(t *testing.T) {
	tests := []struct {
		name        string
		mockSetup   func() *adjuntoManagerDbMock
		pendienteID int
		descripcion string
		archivoPath string
		wantErr     error
	}{
		{name: "Todo Ok",
			mockSetup: func() *adjuntoManagerDbMock {
				return &adjuntoManagerDbMock{CargarFunc: func(ctx context.Context, a *models.Adjunto) (int, error) { return 1, nil }}
			},
			pendienteID: 1,
			descripcion: "desc",
			archivoPath: "/f",
			wantErr:     nil,
		},
		{name: "Pendiente invalido",
			mockSetup: func() *adjuntoManagerDbMock {
				return &adjuntoManagerDbMock{CargarFunc: func(ctx context.Context, a *models.Adjunto) (int, error) { return 0, appErrors.PendienteNoEncontrado }}
			},
			pendienteID: 0,
			descripcion: "d",
			archivoPath: "/f",
			wantErr:     appErrors.ParametroDeCargaVacio,
		},
		{name: "ArchivoPath vacio",
			mockSetup: func() *adjuntoManagerDbMock {
				return &adjuntoManagerDbMock{CargarFunc: func(ctx context.Context, a *models.Adjunto) (int, error) { return 0, nil }}
			},
			pendienteID: 1,
			descripcion: "d",
			archivoPath: " ",
			wantErr:     appErrors.ParametroDeCargaVacio,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAdjuntoService(tt.mockSetup())
			ctx := t.Context()

			_, err := svc.Cargar(ctx, tt.pendienteID, tt.descripcion, tt.archivoPath)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}

func TestObtenerEliminarFiltrarAdjunto(t *testing.T) {
	//now := time.Now()
	tests := []struct {
		name        string
		mockSetup   func() *adjuntoManagerDbMock
		metodo      string // obtener | eliminar | filtrar
		id          int
		pendienteID int
		wantErr     error
	}{
		{name: "Obtener OK",
			mockSetup: func() *adjuntoManagerDbMock {
				return &adjuntoManagerDbMock{ObtenerDetalleFunc: func(ctx context.Context, id int) (*models.Adjunto, error) { return &models.Adjunto{ID: id}, nil }}
			},
			metodo:  "obtener",
			id:      1,
			wantErr: nil,
		},
		{name: "Obtener id invalido",
			mockSetup: func() *adjuntoManagerDbMock {
				return &adjuntoManagerDbMock{ObtenerDetalleFunc: func(ctx context.Context, id int) (*models.Adjunto, error) { return nil, nil }}
			},
			metodo:  "obtener",
			id:      0,
			wantErr: appErrors.ParametroDeBusquedaVacio,
		},
		{name: "Eliminar OK",
			mockSetup: func() *adjuntoManagerDbMock {
				return &adjuntoManagerDbMock{EliminarFunc: func(ctx context.Context, id int) error { return nil }}
			},
			metodo:  "eliminar",
			id:      1,
			wantErr: nil,
		},
		{name: "Eliminar id invalido",
			mockSetup: func() *adjuntoManagerDbMock {
				return &adjuntoManagerDbMock{EliminarFunc: func(ctx context.Context, id int) error { return nil }}
			},
			metodo:  "eliminar",
			id:      0,
			wantErr: appErrors.ParametroDeBusquedaVacio,
		},
		{name: "Filtrar OK",
			mockSetup: func() *adjuntoManagerDbMock {
				return &adjuntoManagerDbMock{FiltrarPorPendienteFunc: func(ctx context.Context, pendienteID int) ([]models.Adjunto, error) {
					return []models.Adjunto{{ID: 1}}, nil
				}}
			},
			metodo:      "filtrar",
			pendienteID: 1,
			wantErr:     nil,
		},
		{name: "Filtrar id invalido",
			mockSetup: func() *adjuntoManagerDbMock {
				return &adjuntoManagerDbMock{FiltrarPorPendienteFunc: func(ctx context.Context, pendienteID int) ([]models.Adjunto, error) { return nil, nil }}
			},
			metodo:      "filtrar",
			pendienteID: 0,
			wantErr:     appErrors.ParametroDeBusquedaVacio,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAdjuntoService(tt.mockSetup())
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

func TestActualizarAdjunto(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *adjuntoManagerDbMock
		metodo    string // path | descripcion
		id        int
		valor     string
		wantErr   error
	}{
		{name: "Actualizar path OK",
			mockSetup: func() *adjuntoManagerDbMock {
				return &adjuntoManagerDbMock{ActualizarPathFunc: func(ctx context.Context, id int, path string) error { return nil }}
			},
			metodo:  "path",
			id:      1,
			valor:   "/p",
			wantErr: nil,
		},
		{name: "Actualizar descripcion OK",
			mockSetup: func() *adjuntoManagerDbMock {
				return &adjuntoManagerDbMock{ActualizarDescripcionFunc: func(ctx context.Context, id int, descripcion string) error { return nil }}
			},
			metodo:  "descripcion",
			id:      1,
			valor:   "d",
			wantErr: nil,
		},
		{name: "Id invalido",
			mockSetup: func() *adjuntoManagerDbMock {
				return &adjuntoManagerDbMock{ActualizarPathFunc: func(ctx context.Context, id int, path string) error { return nil }}
			},
			metodo:  "path",
			id:      0,
			valor:   "/p",
			wantErr: appErrors.ParametroDeBusquedaVacio,
		},
		{name: "Valor vacio",
			mockSetup: func() *adjuntoManagerDbMock {
				return &adjuntoManagerDbMock{ActualizarPathFunc: func(ctx context.Context, id int, path string) error { return nil }}
			},
			metodo:  "path",
			id:      1,
			valor:   " ",
			wantErr: appErrors.ParametroDeCargaVacio,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAdjuntoService(tt.mockSetup())
			ctx := t.Context()

			var err error
			if tt.metodo == "path" {
				err = svc.ActualizarPath(ctx, tt.id, tt.valor)
			} else {
				err = svc.ActualizarDescripcion(ctx, tt.id, tt.valor)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}
