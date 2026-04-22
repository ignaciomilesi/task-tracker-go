package dbmanager

import (
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("avance")
			defer db.Close()

			if test.setup != nil {
				test.setup(db)
			}

			repo := NewAvanceRepository(db)

			_, err := repo.Cargar(test.input)

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
		setup         func(*avanceRepository) int
		errorEsperado error
	}{
		{
			name: "Todo OK",
			setup: func(r *avanceRepository) int {
				id, err := r.Cargar(&models.Avance{
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
				test.id = test.setup(repo)
			}

			result, err := repo.ObtenerDetalle(test.id)

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
		funcionCarga  func(*avanceRepository) int
		errorEsperado error
	}{
		{
			name: "Todo OK",
			funcionCarga: func(r *avanceRepository) int {
				id, err := r.Cargar(&models.Avance{
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
				id = test.funcionCarga(repo)
			}

			err := repo.Eliminar(id)

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
		funcionCarga  func(*avanceRepository)
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
			funcionCarga: func(r *avanceRepository) {
				_, err := r.Cargar(&models.Avance{
					PendienteID: 1,
					Descripcion: "a1",
					Fecha:       time.Now(),
				})
				if err != nil {
					t.Errorf("Error al crear el avance. Detalle:\n%v", err)
				}
				_, err = r.Cargar(&models.Avance{
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
				test.funcionCarga(repo)
			}

			lista, err := repo.FiltrarPorPendiente(test.pendienteID)

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
