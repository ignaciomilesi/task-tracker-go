package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"
)

type documentoManagerDbMock struct {
	CargarFunc               func(context.Context, *models.Documento) error
	ObtenerDetalleFunc       func(context.Context, string) (*models.Documento, error)
	FiltrarPorTipoFunc       func(context.Context, string) ([]models.Documento, error)
	FiltrarPorTituloFunc     func(context.Context, string) ([]models.Documento, error)
	ActualizarPathFunc       func(context.Context, string, string) error
	ActualizarBackupPathFunc func(context.Context, string, string) error
}

func (m *documentoManagerDbMock) Cargar(ctx context.Context, d *models.Documento) error {
	if m.CargarFunc == nil {
		return fmt.Errorf("CargarFunc no implementado")
	}
	return m.CargarFunc(ctx, d)
}

func (m *documentoManagerDbMock) ObtenerDetalle(ctx context.Context, codigo string) (*models.Documento, error) {
	if m.ObtenerDetalleFunc == nil {
		return nil, fmt.Errorf("ObtenerDetalleFunc no implementado")
	}
	return m.ObtenerDetalleFunc(ctx, codigo)
}

func (m *documentoManagerDbMock) FiltrarPorTipo(ctx context.Context, tipo string) ([]models.Documento, error) {
	if m.FiltrarPorTipoFunc == nil {
		return nil, fmt.Errorf("FiltrarPorTipoFunc no implementado")
	}
	return m.FiltrarPorTipoFunc(ctx, tipo)
}

func (m *documentoManagerDbMock) FiltrarPorTitulo(ctx context.Context, titulo string) ([]models.Documento, error) {
	if m.FiltrarPorTituloFunc == nil {
		return nil, fmt.Errorf("FiltrarPorTituloFunc no implementado")
	}
	return m.FiltrarPorTituloFunc(ctx, titulo)
}

func (m *documentoManagerDbMock) ActualizarPath(ctx context.Context, codigo string, nuevoPath string) error {
	if m.ActualizarPathFunc == nil {
		return fmt.Errorf("ActualizarPathFunc no implementado")
	}
	return m.ActualizarPathFunc(ctx, codigo, nuevoPath)
}

func (m *documentoManagerDbMock) ActualizarBackupPath(ctx context.Context, codigo string, nuevoPath string) error {
	if m.ActualizarBackupPathFunc == nil {
		return fmt.Errorf("ActualizarBackupPathFunc no implementado")
	}
	return m.ActualizarBackupPathFunc(ctx, codigo, nuevoPath)
}

func TestCargarDocumento(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *documentoManagerDbMock
		codigo    string
		emision   string
		titulo    string
		tipo      string
		ubicacion string
		backup    string
		wantErr   error
	}{
		{name: "Todo Ok",
			mockSetup: func() *documentoManagerDbMock {
				return &documentoManagerDbMock{CargarFunc: func(ctx context.Context, d *models.Documento) error { return nil }}
			},
			codigo:    "DOC1",
			emision:   "2020",
			titulo:    "T",
			tipo:      "manual",
			ubicacion: "/path",
			backup:    "/b",
			wantErr:   nil,
		},
		{name: "Falta campo obligatorio",
			mockSetup: func() *documentoManagerDbMock {
				return &documentoManagerDbMock{CargarFunc: func(ctx context.Context, d *models.Documento) error { return nil }}
			},
			codigo:    "DOC2",
			emision:   "",
			titulo:    "T",
			tipo:      "manual",
			ubicacion: "",
			backup:    "",
			wantErr:   appErrors.ParametroDeCargaVacio,
		},
		{name: "Codigo vacio",
			mockSetup: func() *documentoManagerDbMock {
				return &documentoManagerDbMock{CargarFunc: func(ctx context.Context, d *models.Documento) error { return nil }}
			},
			codigo:    " ",
			emision:   "2020",
			titulo:    "T",
			tipo:      "manual",
			ubicacion: "",
			backup:    "",
			wantErr:   appErrors.CodigoDocumentoVacio,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewDocumentoService(tt.mockSetup())
			ctx := t.Context()

			err := svc.Cargar(ctx, tt.codigo, tt.emision, tt.titulo, tt.tipo, tt.ubicacion, tt.backup)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}

func TestObtenerDetalleDocumento(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *documentoManagerDbMock
		codigo    string
		wantErr   error
	}{
		{name: "Todo Ok",
			mockSetup: func() *documentoManagerDbMock {
				return &documentoManagerDbMock{ObtenerDetalleFunc: func(ctx context.Context, codigo string) (*models.Documento, error) {
					return &models.Documento{Codigo: codigo}, nil
				}}
			},
			codigo:  "DOC1",
			wantErr: nil,
		},
		{name: "Codigo vacio",
			mockSetup: func() *documentoManagerDbMock {
				return &documentoManagerDbMock{ObtenerDetalleFunc: func(ctx context.Context, codigo string) (*models.Documento, error) { return nil, nil }}
			},
			codigo:  "",
			wantErr: appErrors.ParametroDeBusquedaVacio,
		},
		{name: "No encontrado",
			mockSetup: func() *documentoManagerDbMock {
				return &documentoManagerDbMock{ObtenerDetalleFunc: func(ctx context.Context, codigo string) (*models.Documento, error) {
					return nil, appErrors.DocumentoNoEncontrado
				}}
			},
			codigo:  "NOEX",
			wantErr: appErrors.DocumentoNoEncontrado,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewDocumentoService(tt.mockSetup())
			ctx := t.Context()

			_, err := svc.ObtenerDetalle(ctx, tt.codigo)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}

func TestFiltrarPorTipoYTituloDocumento(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *documentoManagerDbMock
		campo     string
		metodo    string // "tipo" o "titulo"
		wantErr   error
	}{
		{name: "Filtrar por tipo OK",
			mockSetup: func() *documentoManagerDbMock {
				return &documentoManagerDbMock{FiltrarPorTipoFunc: func(ctx context.Context, tipo string) ([]models.Documento, error) {
					return []models.Documento{{Codigo: "DOC1"}}, nil
				}}
			},
			campo:   "manual",
			metodo:  "tipo",
			wantErr: nil,
		},
		{name: "Filtrar por titulo OK",
			mockSetup: func() *documentoManagerDbMock {
				return &documentoManagerDbMock{FiltrarPorTituloFunc: func(ctx context.Context, titulo string) ([]models.Documento, error) {
					return []models.Documento{{Codigo: "DOC1"}}, nil
				}}
			},
			campo:   "Titulo",
			metodo:  "titulo",
			wantErr: nil,
		},
		{name: "Parametro vacio tipo",
			mockSetup: func() *documentoManagerDbMock {
				return &documentoManagerDbMock{FiltrarPorTipoFunc: func(ctx context.Context, tipo string) ([]models.Documento, error) { return nil, nil }}
			},
			campo:   " ",
			metodo:  "tipo",
			wantErr: appErrors.ParametroDeBusquedaVacio,
		},
		{name: "Parametro vacio titulo",
			mockSetup: func() *documentoManagerDbMock {
				return &documentoManagerDbMock{FiltrarPorTituloFunc: func(ctx context.Context, titulo string) ([]models.Documento, error) { return nil, nil }}
			},
			campo:   "",
			metodo:  "titulo",
			wantErr: appErrors.ParametroDeBusquedaVacio,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewDocumentoService(tt.mockSetup())
			ctx := t.Context()

			var err error
			if tt.metodo == "tipo" {
				_, err = svc.FiltrarPorTipo(ctx, tt.campo)
			} else {
				_, err = svc.FiltrarPorTitulo(ctx, tt.campo)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}

func TestActualizarPathsDocumento(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *documentoManagerDbMock
		codigo    string
		nuevoPath string
		metodo    string // "path" o "backup"
		wantErr   error
	}{
		{name: "Actualizar path OK",
			mockSetup: func() *documentoManagerDbMock {
				return &documentoManagerDbMock{ActualizarPathFunc: func(ctx context.Context, codigo string, nuevoPath string) error { return nil }}
			},
			codigo:    "DOC1",
			nuevoPath: "/new",
			metodo:    "path",
			wantErr:   nil,
		},
		{name: "Actualizar backup OK",
			mockSetup: func() *documentoManagerDbMock {
				return &documentoManagerDbMock{ActualizarBackupPathFunc: func(ctx context.Context, codigo string, nuevoPath string) error { return nil }}
			},
			codigo:    "DOC1",
			nuevoPath: "/newb",
			metodo:    "backup",
			wantErr:   nil,
		},
		{name: "Codigo vacio",
			mockSetup: func() *documentoManagerDbMock {
				return &documentoManagerDbMock{ActualizarPathFunc: func(ctx context.Context, codigo string, nuevoPath string) error { return nil }}
			},
			codigo:    "",
			nuevoPath: "/p",
			metodo:    "path",
			wantErr:   appErrors.ParametroDeBusquedaVacio,
		},
		{name: "NuevoPath vacio",
			mockSetup: func() *documentoManagerDbMock {
				return &documentoManagerDbMock{ActualizarPathFunc: func(ctx context.Context, codigo string, nuevoPath string) error { return nil }}
			},
			codigo:    "DOC1",
			nuevoPath: " ",
			metodo:    "path",
			wantErr:   appErrors.ParametroDeCargaVacio,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewDocumentoService(tt.mockSetup())
			ctx := t.Context()

			var err error
			if tt.metodo == "path" {
				err = svc.ActualizarPath(ctx, tt.codigo, tt.nuevoPath)
			} else {
				err = svc.ActualizarBackupPath(ctx, tt.codigo, tt.nuevoPath)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}
