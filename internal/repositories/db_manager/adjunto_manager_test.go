package dbmanager

import (
	"database/sql"
	"errors"
	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"
	"testing"
)

func TestAdjuntoCargar(t *testing.T) {

	tests := []struct {
		name          string
		input         *models.Adjunto
		setup         func(*sql.DB)
		errorEsperado error
	}{
		{
			name: "Todo OK",
			input: &models.Adjunto{
				PendienteID: 1,
				Descripcion: "adjunto inicial",
			},
			setup: func(db *sql.DB) {
				if _, err := db.Exec("PRAGMA foreign_keys = OFF"); err != nil {
					t.Errorf("Error al desactivar FK: %v", err)
				}
			},
			errorEsperado: nil,
		},
		{
			name:          "Adjunto nil",
			input:         nil,
			errorEsperado: appErrors.ParametroDeCargaVacio,
		},
		{
			name: "Pendiente inexistente (FK)",
			input: &models.Adjunto{
				PendienteID: 9999,
				Descripcion: "adjunto",
			},
			errorEsperado: appErrors.PendienteNoEncontrado,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("adjunto")
			defer db.Close()

			if test.setup != nil {
				test.setup(db)
			}

			repo := NewAdjuntoRepository(db)

			_, err := repo.Cargar(test.input)

			if !errors.Is(err, test.errorEsperado) {
				t.Errorf("Error inesperado.\nEsperado: %v\nObtenido: %v",
					test.errorEsperado, err)
			}
		})
	}
}

func TestAdjuntoObtenerDetalle(t *testing.T) {

	tests := []struct {
		name          string
		id            int
		setup         func(*adjuntoRepository) int
		errorEsperado error
	}{
		{
			name: "Todo OK",
			setup: func(r *adjuntoRepository) int {
				id, err := r.Cargar(&models.Adjunto{
					PendienteID: 1,
					Descripcion: "detalle",
				})
				if err != nil {
					t.Errorf("Error al crear adjunto: %v", err)
				}
				return id
			},
			errorEsperado: nil,
		},
		{
			name:          "No encontrado",
			id:            9999,
			errorEsperado: appErrors.AdjuntoNoEncontrado,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("adjunto")
			defer db.Close()

			if _, err := db.Exec("PRAGMA foreign_keys = OFF"); err != nil {
				t.Errorf("Error al desactivar FK: %v", err)
			}

			repo := NewAdjuntoRepository(db)

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

func TestAdjuntoEliminar(t *testing.T) {

	tests := []struct {
		name          string
		funcionCarga  func(*adjuntoRepository) int
		errorEsperado error
	}{
		{
			name: "Todo OK",
			funcionCarga: func(r *adjuntoRepository) int {
				id, err := r.Cargar(&models.Adjunto{
					PendienteID: 1,
					Descripcion: "para borrar",
				})
				if err != nil {
					t.Errorf("Error al crear adjunto: %v", err)
				}
				return id
			},
			errorEsperado: nil,
		},
		{
			name:          "No encontrado",
			funcionCarga:  nil,
			errorEsperado: appErrors.AdjuntoNoEncontrado,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("adjunto")
			defer db.Close()

			if _, err := db.Exec("PRAGMA foreign_keys = OFF"); err != nil {
				t.Errorf("Error al desactivar FK: %v", err)
			}

			repo := NewAdjuntoRepository(db)

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

func TestAdjuntoFiltrarPorPendiente(t *testing.T) {

	tests := []struct {
		name          string
		pendienteID   int
		funcionCarga  func(*adjuntoRepository)
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
			funcionCarga: func(r *adjuntoRepository) {
				_, err := r.Cargar(&models.Adjunto{
					PendienteID: 1,
					Descripcion: "a1",
				})
				if err != nil {
					t.Errorf("Error al crear adjunto: %v", err)
				}
				_, err = r.Cargar(&models.Adjunto{
					PendienteID: 1,
					Descripcion: "a2",
				})
				if err != nil {
					t.Errorf("Error al crear adjunto: %v", err)
				}
			},
			largoEsperado: 2,
			errorEsperado: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("adjunto")
			defer db.Close()

			if _, err := db.Exec("PRAGMA foreign_keys = OFF"); err != nil {
				t.Errorf("Error al desactivar FK: %v", err)
			}

			repo := NewAdjuntoRepository(db)

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
