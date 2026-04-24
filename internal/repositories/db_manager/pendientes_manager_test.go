package dbmanager

import (
	"context"
	"database/sql"
	"errors"
	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"
	"testing"
	"time"
)

func TestPendientesCrear(t *testing.T) {

	tests := []struct {
		name                        string
		input                       *models.Pendientes
		funcionGeneradorSolicitante func(context.Context, *sql.DB, *models.Pendientes)
		errorEsperado               error
	}{
		{
			name: "Todo OK",
			input: &models.Pendientes{
				Titulo:        "titulo de prueba",
				Descripcion:   "decripcion de prueba",
				SolicitanteID: 1,
				FechaPedido:   time.Now(),
			},
			funcionGeneradorSolicitante: func(ctx context.Context, d *sql.DB, p *models.Pendientes) {
				cleanDB(d, "solicitante")
				repo := NewsolicitanteRepository(d)
				id, err := repo.Crear(ctx, &models.Solicitante{
					Nombre: randomString(4),
				})
				if err != nil {
					t.Errorf("Error al crear al solicitante. Detalle:\n%v", err)
				}
				p.SolicitanteID = id
			},
			errorEsperado: nil,
		},
		{
			name: "Solicitante inexistente",
			input: &models.Pendientes{
				Titulo:        "test",
				Descripcion:   "test",
				SolicitanteID: 9999,
				FechaPedido:   time.Now(),
			},
			errorEsperado: appErrors.SolicitanteNoEncontrado,
		},
	}

	ctx := t.Context()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("pendientes")
			defer db.Close()

			repo := NewPendientesRepository(db)

			if test.funcionGeneradorSolicitante != nil {
				test.funcionGeneradorSolicitante(ctx, db, test.input)
			}

			_, err := repo.Crear(ctx, test.input)

			if !errors.Is(err, test.errorEsperado) {
				t.Errorf("Error inesperado.\nEsperado: %v\nObtenido: %v",
					test.errorEsperado, err)
			}
		})
	}
}

func TestPendientesAsignar(t *testing.T) {

	tests := []struct {
		name          string
		asignadoID    int
		setup         func(context.Context, *sql.DB, *int) int
		errorEsperado error
	}{
		{
			name: "Todo OK",
			setup: func(ctx context.Context, d *sql.DB, c *int) int {

				// generamos el colaborador
				cleanDB(d, "colaborador")
				repoCol := NewColaboradorRepository(d)
				id, err := repoCol.Crear(ctx, &models.Colaborador{
					Nombre: randomString(4),
				})
				if err != nil {
					t.Errorf("Error al crear al solicitante. Detalle:\n%v", err)
				}

				*c = id // guardamos el id del colaborador generado

				// generamos el solicitante
				cleanDB(d, "solicitante")
				repoSol := NewsolicitanteRepository(d)
				id, err = repoSol.Crear(ctx, &models.Solicitante{
					Nombre: randomString(4),
				})
				if err != nil {
					t.Errorf("Error al crear al solicitante. Detalle:\n%v", err)
				}

				// generamos un pendiente
				repoPend := NewPendientesRepository(d)
				id, err = repoPend.Crear(ctx, &models.Pendientes{
					Titulo:        "titulo de prueba",
					Descripcion:   "decripcion de prueba",
					SolicitanteID: id,
					FechaPedido:   time.Now(),
				})
				if err != nil {
					t.Errorf("Error al crear el pendiente. Detalle:\n%v", err)
				}
				return id
			},
			errorEsperado: nil,
		},
		{
			name:          "Pendiente no existe",
			errorEsperado: appErrors.PendienteNoEncontrado,
		},
		{
			name:       "Colaborador inexistente",
			asignadoID: 9999,
			setup: func(ctx context.Context, d *sql.DB, c *int) int {

				// No generamos el colaborador

				// generamos el solicitante
				cleanDB(d, "solicitante")
				repoSol := NewsolicitanteRepository(d)
				id, err := repoSol.Crear(ctx, &models.Solicitante{
					Nombre: randomString(4),
				})
				if err != nil {
					t.Errorf("Error al crear al solicitante. Detalle:\n%v", err)
				}

				// generamos un pendiente
				repoPend := NewPendientesRepository(d)
				id, err = repoPend.Crear(ctx, &models.Pendientes{
					Titulo:        "titulo de prueba",
					Descripcion:   "decripcion de prueba",
					SolicitanteID: id,
					FechaPedido:   time.Now(),
				})
				if err != nil {
					t.Errorf("Error al crear el pendiente. Detalle:\n%v", err)
				}
				return id
			},
			errorEsperado: appErrors.ColaboradorNoEncontrado,
		},
	}

	ctx := t.Context()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("pendientes")
			defer db.Close()

			repo := NewPendientesRepository(db)

			id := 999
			if test.setup != nil {
				id = test.setup(ctx, db, &test.asignadoID)
			}
			err := repo.Asignar(ctx, id, test.asignadoID, time.Now())

			if !errors.Is(err, test.errorEsperado) {
				t.Errorf("Error inesperado.\nEsperado: %v\nObtenido: %v",
					test.errorEsperado, err)
			}
		})
	}
}

func TestPendientesTerminar(t *testing.T) {

	tests := []struct {
		name          string
		setup         func(context.Context, *pendientesRepository) int
		errorEsperado error
	}{
		{
			name: "Todo OK",
			setup: func(ctx context.Context, d *pendientesRepository) int {
				id, err := d.Crear(ctx, &models.Pendientes{
					Titulo:        "titulo de prueba",
					Descripcion:   "decripcion de prueba",
					SolicitanteID: 1,
					FechaPedido:   time.Now(),
				})
				if err != nil {
					t.Errorf("Error al crear el pendiente. Detalle:\n%v", err)
				}
				return id
			},
			errorEsperado: nil,
		},
		{
			name: "No existe",
			setup: func(ctx context.Context, r *pendientesRepository) int {
				return 9999
			},
			errorEsperado: appErrors.PendienteNoEncontrado,
		},
	}

	ctx := t.Context()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("pendientes")
			defer db.Close()

			// desactivo la verificación de foreign_keys, este test no lo requiere
			if _, err := db.Exec("PRAGMA foreign_keys = OFF"); err != nil {
				t.Errorf("Error al desactivar el PRAGMA. Detalle: %v", err)
			}

			repo := NewPendientesRepository(db)

			id := test.setup(ctx, repo)

			err := repo.Terminar(ctx, id, nil, time.Now())

			if !errors.Is(err, test.errorEsperado) {
				t.Errorf("Error inesperado.\nEsperado: %v\nObtenido: %v",
					test.errorEsperado, err)
			}
		})
	}
}

func TestPendientesModificarDescripcion(t *testing.T) {

	tests := []struct {
		name          string
		setup         func(context.Context, *pendientesRepository) int
		errorEsperado error
	}{
		{
			name: "Todo OK",
			setup: func(ctx context.Context, d *pendientesRepository) int {
				id, err := d.Crear(ctx, &models.Pendientes{
					Titulo:        "titulo de prueba",
					Descripcion:   "decripcion de prueba",
					SolicitanteID: 1,
					FechaPedido:   time.Now(),
				})
				if err != nil {
					t.Errorf("Error al crear el pendiente. Detalle:\n%v", err)
				}
				return id
			},
			errorEsperado: nil,
		},
		{
			name: "No existe",
			setup: func(ctx context.Context, r *pendientesRepository) int {
				return 9999
			},
			errorEsperado: appErrors.PendienteNoEncontrado,
		},
	}

	ctx := t.Context()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("pendientes")
			defer db.Close()

			// desactivo la verificación de foreign_keys, este test no lo requiere
			if _, err := db.Exec("PRAGMA foreign_keys = OFF"); err != nil {
				t.Errorf("Error al desactivar el PRAGMA. Detalle: %v", err)
			}

			repo := NewPendientesRepository(db)

			id := test.setup(ctx, repo)

			err := repo.ModificarDescripcion(ctx, id, "nueva descripcion")
			if !errors.Is(err, test.errorEsperado) {
				t.Errorf("Error inesperado.\nEsperado: %v\nObtenido: %v",
					test.errorEsperado, err)
			}
		})
	}
}

func TestPendientesBuscarPorTituloDescripcion(t *testing.T) {

	tests := []struct {
		name          string
		texto         string
		setup         func(context.Context, *pendientesRepository)
		errorEsperado error
	}{
		{
			name:          "Texto vacío",
			texto:         " ",
			setup:         func(ctx context.Context, r *pendientesRepository) {},
			errorEsperado: appErrors.ParametroDeBusquedaVacio,
		},
		{
			name:  "Con resultados (titulo)",
			texto: "abc",
			setup: func(ctx context.Context, d *pendientesRepository) {
				_, err := d.Crear(ctx, &models.Pendientes{
					Titulo:        "titulo de prueba abc",
					Descripcion:   "decripción de prueba",
					SolicitanteID: 1,
					FechaPedido:   time.Now(),
				})
				if err != nil {
					t.Errorf("Error al crear el pendiente. Detalle:\n%v", err)
				}

			},
			errorEsperado: nil,
		},
		{
			name:  "Con resultados (descripción)",
			texto: "abc",
			setup: func(ctx context.Context, d *pendientesRepository) {
				_, err := d.Crear(ctx, &models.Pendientes{
					Titulo:        "titulo de prueba",
					Descripcion:   "decripción de prueba abc",
					SolicitanteID: 1,
					FechaPedido:   time.Now(),
				})
				if err != nil {
					t.Errorf("Error al crear el pendiente. Detalle:\n%v", err)
				}

			},
			errorEsperado: nil,
		},
	}

	ctx := t.Context()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("pendientes")
			defer db.Close()

			// desactivo la verificación de foreign_keys, este test no lo requiere
			if _, err := db.Exec("PRAGMA foreign_keys = OFF"); err != nil {
				t.Errorf("Error al desactivar el PRAGMA. Detalle: %v", err)
			}

			repo := NewPendientesRepository(db)

			test.setup(ctx, repo)

			lista, err := repo.BuscarPorTituloDescripcion(ctx, test.texto, false)

			if !errors.Is(err, test.errorEsperado) {
				t.Errorf("Error inesperado.\nEsperado: %v\nObtenido: %v",
					test.errorEsperado, err)
			}

			if err == nil && len(lista) == 0 {
				t.Errorf("Error inesperado. No se obtuvo resultado")
			}
		})
	}
}

func TestPendientesListar(t *testing.T) {

	tests := []struct {
		name              string
		incluirFinalizado bool
		setup             func(context.Context, *pendientesRepository)
		largoEsperado     int
		errorEsperado     error
	}{
		{
			name:          "Lista vacía",
			setup:         func(ctx context.Context, r *pendientesRepository) {},
			largoEsperado: 0,
			errorEsperado: nil,
		},
		{
			name:              "Lista Con resultados (Finalizado = false)",
			incluirFinalizado: false,
			setup: func(ctx context.Context, d *pendientesRepository) {
				_, err := d.Crear(ctx, &models.Pendientes{
					Titulo:        "titulo de prueba abc",
					Descripcion:   "decripción de prueba",
					SolicitanteID: 1,
					FechaPedido:   time.Now(),
				})
				if err != nil {
					t.Errorf("Error al crear el pendiente. Detalle:\n%v", err)
				}

			},
			largoEsperado: 1,
			errorEsperado: nil,
		},
		{
			name:              "Lista Con resultados (Finalizado = true)",
			incluirFinalizado: true,
			setup: func(ctx context.Context, d *pendientesRepository) {
				id, err := d.Crear(ctx, &models.Pendientes{
					Titulo:        "titulo de prueba abc",
					Descripcion:   "decripción de prueba",
					SolicitanteID: 1,
					FechaPedido:   time.Now(),
				})
				if err != nil {
					t.Errorf("Error al crear el pendiente. Detalle:\n%v", err)
				}
				err = d.Terminar(ctx, id, nil, time.Now())

			},
			largoEsperado: 1,
			errorEsperado: nil,
		},
	}

	ctx := t.Context()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("pendientes")
			defer db.Close()

			// desactivo la verificación de foreign_keys, este test no lo requiere
			if _, err := db.Exec("PRAGMA foreign_keys = OFF"); err != nil {
				t.Errorf("Error al desactivar el PRAGMA. Detalle: %v", err)
			}

			repo := NewPendientesRepository(db)

			test.setup(ctx, repo)

			lista, err := repo.Listar(ctx, test.incluirFinalizado, 10, 0)

			if !errors.Is(err, test.errorEsperado) {
				t.Errorf("Error inesperado.\nEsperado: %v\nObtenido: %v",
					test.errorEsperado, err)
			}

			if err == nil && len(lista) < test.largoEsperado {
				t.Errorf("Error inesperado. No se obtuvo el largo buscado. Se esperaba: %d pero se obtuvo %d", test.largoEsperado, len(lista))
			}
		})
	}
}

func TestPendientesModificarIdentificacionTablaPendiente(t *testing.T) {

	tests := []struct {
		name          string
		setup         func(context.Context, *pendientesRepository) int
		errorEsperado error
	}{
		{
			name: "Todo OK",
			setup: func(ctx context.Context, d *pendientesRepository) int {
				id, err := d.Crear(ctx, &models.Pendientes{
					Titulo:        "titulo de prueba",
					Descripcion:   "decripcion de prueba",
					SolicitanteID: 1,
					FechaPedido:   time.Now(),
				})
				if err != nil {
					t.Errorf("Error al crear el pendiente. Detalle:\n%v", err)
				}
				return id
			},
			errorEsperado: nil,
		},
		{
			name: "No existe",
			setup: func(ctx context.Context, r *pendientesRepository) int {
				return 9999
			},
			errorEsperado: appErrors.PendienteNoEncontrado,
		},
	}

	ctx := t.Context()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("pendientes")
			defer db.Close()

			// desactivo la verificación de foreign_keys, este test no lo requiere
			if _, err := db.Exec("PRAGMA foreign_keys = OFF"); err != nil {
				t.Errorf("Error al desactivar el PRAGMA. Detalle: %v", err)
			}

			repo := NewPendientesRepository(db)

			id := test.setup(ctx, repo)

			err := repo.ModificarIdentificacionTablaPendiente(ctx, id, "nueva descripcion")
			if !errors.Is(err, test.errorEsperado) {
				t.Errorf("Error inesperado.\nEsperado: %v\nObtenido: %v",
					test.errorEsperado, err)
			}
		})
	}
}

func TestPendientesListarPorAsignado(t *testing.T) {

	tests := []struct {
		name          string
		asignadoID    int
		finalizado    bool
		setup         func(context.Context, *pendientesRepository)
		largoEsperado int
		errorEsperado error
	}{
		{
			name:          "Sin resultados",
			asignadoID:    1,
			finalizado:    false,
			largoEsperado: 0,
			errorEsperado: nil,
		},
		{
			name:       "Con resultados",
			asignadoID: 2,
			finalizado: false,
			setup: func(ctx context.Context, r *pendientesRepository) {

				id, err := r.Crear(ctx, &models.Pendientes{
					Titulo:        "titulo de prueba",
					Descripcion:   "decripcion de prueba",
					SolicitanteID: 1,
					FechaPedido:   time.Now(),
				})
				if err != nil {
					t.Errorf("Error al crear el pendiente. Detalle:\n%v", err)
				}
				err = r.Asignar(ctx, id, 2, time.Now())
				if err != nil {
					t.Errorf("Error al asignar el pendiente. Detalle:\n%v", err)
				}
				id2, err := r.Crear(ctx, &models.Pendientes{
					Titulo:        "titulo de prueba",
					Descripcion:   "decripcion de prueba",
					SolicitanteID: 1,
					FechaPedido:   time.Now(),
				})
				if err != nil {
					t.Errorf("Error al crear el pendiente. Detalle:\n%v", err)
				}
				err = r.Asignar(ctx, id2, 2, time.Now())
				if err != nil {
					t.Errorf("Error al asignar el pendiente. Detalle:\n%v", err)
				}
			},
			largoEsperado: 2,
			errorEsperado: nil,
		},
		{
			name:       "Filtra por finalizado",
			asignadoID: 3,
			finalizado: true,
			setup: func(ctx context.Context, r *pendientesRepository) {

				id, err := r.Crear(ctx, &models.Pendientes{
					Titulo:        "titulo de prueba",
					Descripcion:   "decripcion de prueba",
					SolicitanteID: 1,
					FechaPedido:   time.Now(),
				})
				if err != nil {
					t.Errorf("Error al crear el pendiente. Detalle:\n%v", err)
				}
				err = r.Asignar(ctx, id, 3, time.Now())
				if err != nil {
					t.Errorf("Error al asignar el pendiente. Detalle:\n%v", err)
				}
				err = r.Terminar(ctx, id, nil, time.Now())
				if err != nil {
					t.Errorf("Error al terminar el pendiente. Detalle:\n%v", err)
				}
			},
			largoEsperado: 1,
			errorEsperado: nil,
		},
	}

	ctx := t.Context()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("pendientes")
			defer db.Close()

			// desactivo la verificación de foreign_keys, este test no lo requiere
			if _, err := db.Exec("PRAGMA foreign_keys = OFF"); err != nil {
				t.Errorf("Error al desactivar el PRAGMA. Detalle: %v", err)
			}

			repo := NewPendientesRepository(db)

			if test.setup != nil {
				test.setup(ctx, repo)
			}

			lista, err := repo.ListarPorAsignado(ctx, test.asignadoID, test.finalizado)

			if !errors.Is(err, test.errorEsperado) {
				t.Errorf("Error inesperado.\nEsperado: %v\nObtenido: %v",
					test.errorEsperado, err)
			}

			if len(lista) < test.largoEsperado {
				t.Errorf("Cantidad incorrecta.\nEsperado: %d\nObtenido: %d",
					test.largoEsperado, len(lista))
			}
		})
	}
}
