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

type pendientesManagerDbMock struct {
	CrearFunc                        func(context.Context, *models.Pendientes) (int, error)
	AsignarFunc                      func(context.Context, int, int, time.Time) error
	TerminarFunc                     func(context.Context, int, *string, time.Time) error
	ModificarDescripcionFunc         func(context.Context, int, string) error
	ModificarIdentificacionTablaFunc func(context.Context, int, string) error
	BuscarPorTituloDescripcionFunc   func(context.Context, string, bool) ([]models.Pendientes, error)
	ListarPorAsignadoFunc            func(context.Context, int, bool) ([]models.Pendientes, error)
	ListarFunc                       func(context.Context, bool, int, int) ([]models.Pendientes, error)
}

func (m *pendientesManagerDbMock) Crear(ctx context.Context, p *models.Pendientes) (int, error) {
	if m.CrearFunc == nil {
		return 0, fmt.Errorf("CrearFunc no implementado")
	}
	return m.CrearFunc(ctx, p)
}

func (m *pendientesManagerDbMock) Asignar(ctx context.Context, id int, asignadoID int, fecha time.Time) error {
	if m.AsignarFunc == nil {
		return fmt.Errorf("AsignarFunc no implementado")
	}
	return m.AsignarFunc(ctx, id, asignadoID, fecha)
}

func (m *pendientesManagerDbMock) Terminar(ctx context.Context, id int, cierre *string, fecha time.Time) error {
	if m.TerminarFunc == nil {
		return fmt.Errorf("TerminarFunc no implementado")
	}
	return m.TerminarFunc(ctx, id, cierre, fecha)
}

func (m *pendientesManagerDbMock) ModificarDescripcion(ctx context.Context, id int, descripcion string) error {
	if m.ModificarDescripcionFunc == nil {
		return fmt.Errorf("ModificarDescripcionFunc no implementado")
	}
	return m.ModificarDescripcionFunc(ctx, id, descripcion)
}

func (m *pendientesManagerDbMock) ModificarIdentificacionTablaPendiente(ctx context.Context, id int, identificacion string) error {
	if m.ModificarIdentificacionTablaFunc == nil {
		return fmt.Errorf("ModificarIdentificacionTablaFunc no implementado")
	}
	return m.ModificarIdentificacionTablaFunc(ctx, id, identificacion)
}

func (m *pendientesManagerDbMock) BuscarPorTituloDescripcion(ctx context.Context, texto string, finalizado bool) ([]models.Pendientes, error) {
	if m.BuscarPorTituloDescripcionFunc == nil {
		return nil, fmt.Errorf("BuscarPorTituloDescripcionFunc no implementado")
	}
	return m.BuscarPorTituloDescripcionFunc(ctx, texto, finalizado)
}

func (m *pendientesManagerDbMock) ListarPorAsignado(ctx context.Context, asignadoID int, finalizado bool) ([]models.Pendientes, error) {
	if m.ListarPorAsignadoFunc == nil {
		return nil, fmt.Errorf("ListarPorAsignadoFunc no implementado")
	}
	return m.ListarPorAsignadoFunc(ctx, asignadoID, finalizado)
}

func (m *pendientesManagerDbMock) Listar(ctx context.Context, finalizado bool, limit, offset int) ([]models.Pendientes, error) {
	if m.ListarFunc == nil {
		return nil, fmt.Errorf("ListarFunc no implementado")
	}
	return m.ListarFunc(ctx, finalizado, limit, offset)
}

func TestCrearPendiente(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name          string
		mockSetup     func() *pendientesManagerDbMock
		titulo        string
		descripcion   string
		solicitanteID int
		fechaPedido   time.Time
		wantErr       error
	}{
		{name: "Todo Ok",
			mockSetup: func() *pendientesManagerDbMock {
				return &pendientesManagerDbMock{CrearFunc: func(ctx context.Context, p *models.Pendientes) (int, error) { return 1, nil }}
			},
			titulo:        "T",
			descripcion:   "D",
			solicitanteID: 1,
			fechaPedido:   now.Add(-time.Hour),
			wantErr:       nil,
		},
		{name: "Titulo vacio",
			mockSetup: func() *pendientesManagerDbMock {
				return &pendientesManagerDbMock{CrearFunc: func(ctx context.Context, p *models.Pendientes) (int, error) { return 0, nil }}
			},
			titulo:        " ",
			descripcion:   "D",
			solicitanteID: 1,
			fechaPedido:   now,
			wantErr:       appErrors.ParametroDeCargaVacio,
		},
		{name: "Solicitante invalido",
			mockSetup: func() *pendientesManagerDbMock {
				return &pendientesManagerDbMock{CrearFunc: func(ctx context.Context, p *models.Pendientes) (int, error) {
					return 0, appErrors.SolicitanteNoEncontrado
				}}
			},
			titulo:        "T",
			descripcion:   "D",
			solicitanteID: 0,
			fechaPedido:   now,
			wantErr:       appErrors.ParametroDeCargaVacio,
		},
		{name: "Fecha futura",
			mockSetup: func() *pendientesManagerDbMock {
				return &pendientesManagerDbMock{CrearFunc: func(ctx context.Context, p *models.Pendientes) (int, error) { return 0, nil }}
			},
			titulo:        "T",
			descripcion:   "D",
			solicitanteID: 1,
			fechaPedido:   now.Add(48 * time.Hour),
			wantErr:       appErrors.FechaNoValida,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewPendientesService(tt.mockSetup())
			ctx := t.Context()

			_, err := svc.Crear(ctx, tt.titulo, tt.descripcion, tt.solicitanteID, tt.fechaPedido)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}

func TestAsignarPendiente(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name       string
		mockSetup  func() *pendientesManagerDbMock
		id         int
		asignadoID int
		fecha      time.Time
		wantErr    error
	}{
		{name: "Todo Ok",
			mockSetup: func() *pendientesManagerDbMock {
				return &pendientesManagerDbMock{AsignarFunc: func(ctx context.Context, id int, asignadoID int, fecha time.Time) error { return nil }}
			},
			id:         1,
			asignadoID: 2,
			fecha:      now,
			wantErr:    nil,
		},
		{name: "Id invalido",
			mockSetup: func() *pendientesManagerDbMock {
				return &pendientesManagerDbMock{AsignarFunc: func(ctx context.Context, id int, asignadoID int, fecha time.Time) error { return nil }}
			},
			id:         0,
			asignadoID: 2,
			fecha:      now,
			wantErr:    appErrors.ParametroDeBusquedaVacio,
		},
		{name: "Fecha futura",
			mockSetup: func() *pendientesManagerDbMock {
				return &pendientesManagerDbMock{AsignarFunc: func(ctx context.Context, id int, asignadoID int, fecha time.Time) error { return nil }}
			},
			id:         1,
			asignadoID: 2,
			fecha:      now.Add(48 * time.Hour),
			wantErr:    appErrors.FechaNoValida,
		},
		{name: "Colaborador no encontrado (repo)",
			mockSetup: func() *pendientesManagerDbMock {
				return &pendientesManagerDbMock{AsignarFunc: func(ctx context.Context, id int, asignadoID int, fecha time.Time) error {
					return appErrors.ColaboradorNoEncontrado
				}}
			},
			id:         1,
			asignadoID: 99,
			fecha:      now,
			wantErr:    appErrors.ColaboradorNoEncontrado,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewPendientesService(tt.mockSetup())
			ctx := t.Context()

			err := svc.Asignar(ctx, tt.id, tt.asignadoID, tt.fecha)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}

func TestTerminarPendiente(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name      string
		mockSetup func() *pendientesManagerDbMock
		id        int
		cierre    string
		fecha     time.Time
		wantErr   error
	}{
		{name: "Todo Ok sin cierre",
			mockSetup: func() *pendientesManagerDbMock {
				return &pendientesManagerDbMock{TerminarFunc: func(ctx context.Context, id int, cierre *string, fecha time.Time) error { return nil }}
			},
			id:      1,
			cierre:  "",
			fecha:   now,
			wantErr: nil,
		},
		{name: "Todo Ok con cierre",
			mockSetup: func() *pendientesManagerDbMock {
				return &pendientesManagerDbMock{TerminarFunc: func(ctx context.Context, id int, cierre *string, fecha time.Time) error { return nil }}
			},
			id:      1,
			cierre:  "Cierre",
			fecha:   now,
			wantErr: nil,
		},
		{name: "Id invalido",
			mockSetup: func() *pendientesManagerDbMock {
				return &pendientesManagerDbMock{TerminarFunc: func(ctx context.Context, id int, cierre *string, fecha time.Time) error { return nil }}
			},
			id:      0,
			cierre:  "",
			fecha:   now,
			wantErr: appErrors.ParametroDeBusquedaVacio,
		},
		{name: "Fecha futura",
			mockSetup: func() *pendientesManagerDbMock {
				return &pendientesManagerDbMock{TerminarFunc: func(ctx context.Context, id int, cierre *string, fecha time.Time) error { return nil }}
			},
			id:      1,
			cierre:  "",
			fecha:   now.Add(48 * time.Hour),
			wantErr: appErrors.FechaNoValida,
		},
		{name: "Pendiente no encontrado (repo)",
			mockSetup: func() *pendientesManagerDbMock {
				return &pendientesManagerDbMock{TerminarFunc: func(ctx context.Context, id int, cierre *string, fecha time.Time) error {
					return appErrors.PendienteNoEncontrado
				}}
			},
			id:      999,
			cierre:  "",
			fecha:   now,
			wantErr: appErrors.PendienteNoEncontrado,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewPendientesService(tt.mockSetup())
			ctx := t.Context()

			err := svc.Terminar(ctx, tt.id, tt.cierre, tt.fecha)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}

func TestModificarDescripcionYIdentificacionPendiente(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func() *pendientesManagerDbMock
		metodo    string // "descripcion" o "identificacion"
		id        int
		valor     string
		wantErr   error
	}{
		{name: "Modificar descripcion OK",
			mockSetup: func() *pendientesManagerDbMock {
				return &pendientesManagerDbMock{ModificarDescripcionFunc: func(ctx context.Context, id int, descripcion string) error { return nil }}
			},
			metodo:  "descripcion",
			id:      1,
			valor:   "Nueva",
			wantErr: nil,
		},
		{name: "Modificar identificacion OK",
			mockSetup: func() *pendientesManagerDbMock {
				return &pendientesManagerDbMock{ModificarIdentificacionTablaFunc: func(ctx context.Context, id int, identificacion string) error { return nil }}
			},
			metodo:  "identificacion",
			id:      1,
			valor:   "ID_TAB",
			wantErr: nil,
		},
		{name: "Id invalido",
			mockSetup: func() *pendientesManagerDbMock {
				return &pendientesManagerDbMock{ModificarDescripcionFunc: func(ctx context.Context, id int, descripcion string) error { return nil }}
			},
			metodo:  "descripcion",
			id:      0,
			valor:   "X",
			wantErr: appErrors.ParametroDeBusquedaVacio,
		},
		{name: "Valor vacio",
			mockSetup: func() *pendientesManagerDbMock {
				return &pendientesManagerDbMock{ModificarDescripcionFunc: func(ctx context.Context, id int, descripcion string) error { return nil }}
			},
			metodo:  "descripcion",
			id:      1,
			valor:   " ",
			wantErr: appErrors.ParametroDeCargaVacio,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewPendientesService(tt.mockSetup())
			ctx := t.Context()

			var err error
			if tt.metodo == "descripcion" {
				err = svc.ModificarDescripcion(ctx, tt.id, tt.valor)
			} else {
				err = svc.ModificarIdentificacionTablaPendiente(ctx, tt.id, tt.valor)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}

func TestBuscarYListarPendientes(t *testing.T) {
	tests := []struct {
		name       string
		mockSetup  func() *pendientesManagerDbMock
		metodo     string // "buscar" | "listarPorAsignado" | "listar"
		texto      string
		asignadoID int
		finalizado bool
		limit      int
		offset     int
		wantErr    error
	}{
		{name: "Buscar OK",
			mockSetup: func() *pendientesManagerDbMock {
				return &pendientesManagerDbMock{BuscarPorTituloDescripcionFunc: func(ctx context.Context, texto string, finalizado bool) ([]models.Pendientes, error) {
					return []models.Pendientes{{ID: 1}}, nil
				}}
			},
			metodo:     "buscar",
			texto:      "titulo",
			finalizado: false,
			wantErr:    nil,
		},
		{name: "Buscar parametro vacio",
			mockSetup: func() *pendientesManagerDbMock {
				return &pendientesManagerDbMock{BuscarPorTituloDescripcionFunc: func(ctx context.Context, texto string, finalizado bool) ([]models.Pendientes, error) { return nil, nil }}
			},
			metodo:     "buscar",
			texto:      " ",
			finalizado: false,
			wantErr:    appErrors.ParametroDeBusquedaVacio,
		},
		{name: "ListarPorAsignado OK",
			mockSetup: func() *pendientesManagerDbMock {
				return &pendientesManagerDbMock{ListarPorAsignadoFunc: func(ctx context.Context, asignadoID int, finalizado bool) ([]models.Pendientes, error) {
					return []models.Pendientes{{ID: 1}}, nil
				}}
			},
			metodo:     "listarPorAsignado",
			asignadoID: 2,
			finalizado: true,
			wantErr:    nil,
		},
		{name: "ListarPorAsignado id invalido",
			mockSetup: func() *pendientesManagerDbMock {
				return &pendientesManagerDbMock{ListarPorAsignadoFunc: func(ctx context.Context, asignadoID int, finalizado bool) ([]models.Pendientes, error) {
					return nil, nil
				}}
			},
			metodo:     "listarPorAsignado",
			asignadoID: 0,
			finalizado: true,
			wantErr:    appErrors.ParametroDeBusquedaVacio,
		},
		{name: "Listar OK",
			mockSetup: func() *pendientesManagerDbMock {
				return &pendientesManagerDbMock{ListarFunc: func(ctx context.Context, finalizado bool, limit, offset int) ([]models.Pendientes, error) {
					return []models.Pendientes{{ID: 1}}, nil
				}}
			},
			metodo:     "listar",
			finalizado: false,
			limit:      10,
			offset:     0,
			wantErr:    nil,
		},
		{name: "Listar parametros invalidos",
			mockSetup: func() *pendientesManagerDbMock {
				return &pendientesManagerDbMock{ListarFunc: func(ctx context.Context, finalizado bool, limit, offset int) ([]models.Pendientes, error) {
					return nil, nil
				}}
			},
			metodo:     "listar",
			finalizado: false,
			limit:      0,
			offset:     -1,
			wantErr:    appErrors.ParametrosDeListaInvalidos,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewPendientesService(tt.mockSetup())
			ctx := t.Context()

			var err error
			switch tt.metodo {
			case "buscar":
				_, err = svc.BuscarPorTituloDescripcion(ctx, tt.texto, tt.finalizado)
			case "listarPorAsignado":
				_, err = svc.ListarPorAsignado(ctx, tt.asignadoID, tt.finalizado)
			case "listar":
				_, err = svc.Listar(ctx, tt.finalizado, tt.limit, tt.offset)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Error no esperado.\nSe esperaba: %v\nSe obtuvo: %v", tt.wantErr, err)
			}
		})
	}
}
