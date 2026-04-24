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

func TestAvanceCargar(t *testing.T) {

	tests := []struct {
		name          string
		input         *models.Avance
		setup         func(*sql.DB)
		errorEsperado error
	}{
		{
			name: "Todo OK",
			input: &models.Avance{
				PendienteID: 1,
				Descripcion: "avance inicial",
				Fecha:       time.Now(),
			},
			setup: func(db *sql.DB) {
				// desactivo la verificación de foreign_keys, este test no lo requiere
				if _, err := db.Exec("PRAGMA foreign_keys = OFF"); err != nil {
					t.Errorf("Error al desactivar el PRAGMA. Detalle: %v", err)
				}
			},
			errorEsperado: nil,
		},
		{
			name:  "Avance nil",
			input: nil,
			setup: func(db *sql.DB) {
				// desactivo la verificación de foreign_keys, este test no lo requiere
				if _, err := db.Exec("PRAGMA foreign_keys = OFF"); err != nil {
					t.Errorf("Error al desactivar el PRAGMA. Detalle: %v", err)
				}
			},
			errorEsperado: appErrors.ParametroDeCargaVacio,
		},
		{
			name: "Pendiente inexistente (FK)",
			input: &models.Avance{
				PendienteID: 9999,
				Descripcion: "avance",
				Fecha:       time.Now(),
			},
			errorEsperado: appErrors.PendienteNoEncontrado,
		},
	}
	ctx := t.Context()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("avance")
			defer db.Close()

			if test.setup != nil {
				test.setup(db)
			}

			repo := NewAvanceRepository(db)

			_, err := repo.Cargar(ctx, test.input)

			if !errors.Is(err, test.errorEsperado) {
				t.Errorf("Error inesperado.\nEsperado: %v\nObtenido: %v",
					test.errorEsperado, err)
			}
		})
	}
}

func TestAvanceObtenerDetalle(t *testing.T) {

	tests := []struct {
		name          string
		id            int
		setup         func(context.Context, *avanceRepository) int
		errorEsperado error
	}{
		{
			name: "Todo OK",
			setup: func(ctx context.Context, r *avanceRepository) int {
				id, err := r.Cargar(ctx, &models.Avance{
					PendienteID: 1,
					Descripcion: "detalle",
					Fecha:       time.Now(),
				})
				if err != nil {
					t.Errorf("Error al crear el avance. Detalle:\n%v", err)
				}
				return id
			},
			errorEsperado: nil,
		},
		{
			name:          "No encontrado",
			id:            9999,
			errorEsperado: appErrors.AvanceNoEncontrado,
		},
	}
	ctx := t.Context()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("avance")
			defer db.Close()

			// desactivo la verificación de foreign_keys, este test no lo requiere
			if _, err := db.Exec("PRAGMA foreign_keys = OFF"); err != nil {
				t.Errorf("Error al desactivar el PRAGMA. Detalle: %v", err)
			}
			repo := NewAvanceRepository(db)

			if test.setup != nil {
				test.id = test.setup(ctx, repo)
			}

			result, err := repo.ObtenerDetalle(ctx, test.id)

			if !errors.Is(err, test.errorEsperado) {
				t.Errorf("Error inesperado.\nEsperado: %v\nObtenido: %v",
					test.errorEsperado, err)
			}

			if err == nil && result == nil {
				t.Errorf("Se esperaba resultado y se obtuvo nil")
			}
		})
	}
}

func TestAvanceEliminar(t *testing.T) {

	tests := []struct {
		name          string
		funcionCarga  func(context.Context, *avanceRepository) int
		errorEsperado error
	}{
		{
			name: "Todo OK",
			funcionCarga: func(ctx context.Context, r *avanceRepository) int {
				id, err := r.Cargar(ctx, &models.Avance{
					PendienteID: 1,
					Descripcion: "para borrar",
					Fecha:       time.Now(),
				})
				if err != nil {
					t.Errorf("Error al crear el avance. Detalle:\n%v", err)
				}
				return id
			},
			errorEsperado: nil,
		},
		{
			name:          "No encontrado",
			funcionCarga:  nil,
			errorEsperado: appErrors.AvanceNoEncontrado,
		},
	}
	ctx := t.Context()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("avance")
			defer db.Close()

			// desactivo la verificación de foreign_keys, este test no lo requiere
			if _, err := db.Exec("PRAGMA foreign_keys = OFF"); err != nil {
				t.Errorf("Error al desactivar el PRAGMA. Detalle: %v", err)
			}

			repo := NewAvanceRepository(db)

			id := 9999
			if test.funcionCarga != nil {
				id = test.funcionCarga(ctx, repo)
			}

			err := repo.Eliminar(ctx, id)

			if !errors.Is(err, test.errorEsperado) {
				t.Errorf("Error inesperado.\nEsperado: %v\nObtenido: %v",
					test.errorEsperado, err)
			}
		})
	}
}

func TestAvanceFiltrarPorPendiente(t *testing.T) {

	tests := []struct {
		name          string
		pendienteID   int
		funcionCarga  func(context.Context, *avanceRepository)
		largoEsperado int
		errorEsperado error
	}{
		{
			name:          "Sin resultados",
			pendienteID:   1,
			largoEsperado: 0,
			errorEsperado: nil,
		},
		{
			name:        "Con resultados",
			pendienteID: 1,
			funcionCarga: func(ctx context.Context, r *avanceRepository) {
				_, err := r.Cargar(ctx, &models.Avance{
					PendienteID: 1,
					Descripcion: "a1",
					Fecha:       time.Now(),
				})
				if err != nil {
					t.Errorf("Error al crear el avance. Detalle:\n%v", err)
				}
				_, err = r.Cargar(ctx, &models.Avance{
					PendienteID: 1,
					Descripcion: "a2",
					Fecha:       time.Now(),
				})
				if err != nil {
					t.Errorf("Error al crear el avance. Detalle:\n%v", err)
				}
			},
			largoEsperado: 2,
			errorEsperado: nil,
		},
	}
	ctx := t.Context()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("avance")
			defer db.Close()

			// desactivo la verificación de foreign_keys, este test no lo requiere
			if _, err := db.Exec("PRAGMA foreign_keys = OFF"); err != nil {
				t.Errorf("Error al desactivar el PRAGMA. Detalle: %v", err)
			}

			repo := NewAvanceRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(ctx, repo)
			}

			lista, err := repo.FiltrarPorPendiente(ctx, test.pendienteID)

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
