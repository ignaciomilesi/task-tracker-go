package dbmanager

import (
	"errors"
	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"
	"testing"
	"time"
)

func TestCodigoIDCargar(t *testing.T) {

	tests := []struct {
		name          string
		input         *models.CodigoID
		funcionCarga  func(*codigoIDRepository, *models.CodigoID)
		errorEsperado error
	}{
		{
			name: "Todo OK",
			input: &models.CodigoID{
				Codigo:      randomString(6),
				Estado:      "pendiente",
				FechaPedido: time.Now(),
			},
			errorEsperado: nil,
		},
		{
			name: "Codigo vacío",
			input: &models.CodigoID{
				Codigo: "   ",
			},
			errorEsperado: appErrors.CodigoIDVacio,
		},
		{
			name: "Codigo duplicado",
			input: &models.CodigoID{
				Codigo:      randomString(6),
				Estado:      "pendiente",
				FechaPedido: time.Now(),
			},
			funcionCarga: func(r *codigoIDRepository, c *models.CodigoID) {
				r.Cargar(c)
			},
			errorEsperado: appErrors.CodigoIDDuplicado,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("codigo_ID")
			defer db.Close()

			repo := NewCodigoIDRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(repo, test.input)
			}

			err := repo.Cargar(test.input)

			if !errors.Is(err, test.errorEsperado) {
				t.Errorf("Error inesperado.\nEsperado: %v\nObtenido: %v",
					test.errorEsperado, err)
			}
		})
	}
}

func TestCodigoIDObtenerDetalle(t *testing.T) {

	tests := []struct {
		name          string
		codigo        string
		funcionCarga  func(*codigoIDRepository, string)
		errorEsperado error
	}{
		{
			name:   "Todo OK",
			codigo: randomString(6),
			funcionCarga: func(r *codigoIDRepository, codigo string) {
				r.Cargar(&models.CodigoID{
					Codigo:      codigo,
					Estado:      "pendiente",
					FechaPedido: time.Now(),
				})
			},
			errorEsperado: nil,
		},
		{
			name:          "Codigo vacío",
			codigo:        " ",
			errorEsperado: appErrors.CodigoIDVacio,
		},
		{
			name:          "No encontrado",
			codigo:        "no_existe",
			errorEsperado: appErrors.CodigoIDNoEncontrado,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("codigo_ID")
			defer db.Close()

			repo := NewCodigoIDRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(repo, test.codigo)
			}

			result, err := repo.ObtenerDetalle(test.codigo)

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

func TestCodigoIDFiltrarPorEstado(t *testing.T) {

	tests := []struct {
		name          string
		estado        string
		funcionCarga  func(*codigoIDRepository)
		largoEsperado int
		errorEsperado error
	}{
		{
			name:          "Sin resultados",
			estado:        "inexistente",
			largoEsperado: 0,
			errorEsperado: nil,
		},
		{
			name:   "Con resultados",
			estado: "pendiente",
			funcionCarga: func(r *codigoIDRepository) {
				r.Cargar(&models.CodigoID{
					Codigo:      randomString(6),
					Estado:      "pendiente",
					FechaPedido: time.Now(),
				})
				r.Cargar(&models.CodigoID{
					Codigo:      randomString(6),
					Estado:      "pendiente",
					FechaPedido: time.Now(),
				})
			},
			largoEsperado: 2,
			errorEsperado: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("codigo_ID")
			defer db.Close()

			repo := NewCodigoIDRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(repo)
			}

			lista, err := repo.FiltrarPorEstado(test.estado)

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

func TestCodigoIDActualizarEstado(t *testing.T) {

	tests := []struct {
		name          string
		codigo        string
		funcionCarga  func(*codigoIDRepository, string)
		errorEsperado error
	}{
		{
			name:   "Todo OK",
			codigo: randomString(6),
			funcionCarga: func(r *codigoIDRepository, codigo string) {
				r.Cargar(&models.CodigoID{
					Codigo:      codigo,
					Estado:      "pendiente",
					FechaPedido: time.Now(),
				})
			},
			errorEsperado: nil,
		},
		{
			name:          "Codigo vacío",
			codigo:        " ",
			errorEsperado: appErrors.CodigoIDVacio,
		},
		{
			name:          "No encontrado",
			codigo:        "no_existe",
			errorEsperado: appErrors.CodigoIDNoEncontrado,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("codigo_ID")
			defer db.Close()

			repo := NewCodigoIDRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(repo, test.codigo)
			}

			err := repo.ActualizarEstado(test.codigo, "aprobado", time.Now())

			if !errors.Is(err, test.errorEsperado) {
				t.Errorf("Error inesperado.\nEsperado: %v\nObtenido: %v",
					test.errorEsperado, err)
			}
		})
	}
}

func TestCodigoIDListar(t *testing.T) {

	tests := []struct {
		name          string
		funcionCarga  func(*codigoIDRepository)
		largoEsperado int
		errorEsperado error
	}{
		{
			name:          "Lista vacía",
			largoEsperado: 0,
			errorEsperado: nil,
		},
		{
			name: "Con datos",
			funcionCarga: func(r *codigoIDRepository) {
				r.Cargar(&models.CodigoID{
					Codigo:      randomString(6),
					Estado:      "pendiente",
					FechaPedido: time.Now(),
				})
			},
			largoEsperado: 1,
			errorEsperado: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("codigo_ID")
			defer db.Close()

			repo := NewCodigoIDRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(repo)
			}

			lista, err := repo.Listar(10, 0)

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
