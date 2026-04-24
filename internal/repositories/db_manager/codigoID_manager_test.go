package dbmanager

import (
	"context"
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
		funcionCarga  func(context.Context, *codigoIDRepository, *models.CodigoID)
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
			funcionCarga: func(ctx context.Context, r *codigoIDRepository, c *models.CodigoID) {
				r.Cargar(ctx, c)
			},
			errorEsperado: appErrors.CodigoIDDuplicado,
		},
	}
	ctx := t.Context()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("codigo_ID")
			defer db.Close()

			repo := NewCodigoIDRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(ctx, repo, test.input)
			}

			err := repo.Cargar(ctx, test.input)

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
		funcionCarga  func(context.Context, *codigoIDRepository, string)
		errorEsperado error
	}{
		{
			name:   "Todo OK",
			codigo: randomString(6),
			funcionCarga: func(ctx context.Context, r *codigoIDRepository, codigo string) {
				r.Cargar(ctx, &models.CodigoID{
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
	ctx := t.Context()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("codigo_ID")
			defer db.Close()

			repo := NewCodigoIDRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(ctx, repo, test.codigo)
			}

			result, err := repo.ObtenerDetalle(ctx, test.codigo)

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
		funcionCarga  func(context.Context, *codigoIDRepository)
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
			funcionCarga: func(ctx context.Context, r *codigoIDRepository) {
				r.Cargar(ctx, &models.CodigoID{
					Codigo:      randomString(6),
					Estado:      "pendiente",
					FechaPedido: time.Now(),
				})
				r.Cargar(ctx, &models.CodigoID{
					Codigo:      randomString(6),
					Estado:      "pendiente",
					FechaPedido: time.Now(),
				})
			},
			largoEsperado: 2,
			errorEsperado: nil,
		},
	}
	ctx := t.Context()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("codigo_ID")
			defer db.Close()

			repo := NewCodigoIDRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(ctx, repo)
			}

			lista, err := repo.FiltrarPorEstado(ctx, test.estado)

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
		funcionCarga  func(context.Context, *codigoIDRepository, string)
		errorEsperado error
	}{
		{
			name:   "Todo OK",
			codigo: randomString(6),
			funcionCarga: func(ctx context.Context, r *codigoIDRepository, codigo string) {
				r.Cargar(ctx, &models.CodigoID{
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
	ctx := t.Context()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("codigo_ID")
			defer db.Close()

			repo := NewCodigoIDRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(ctx, repo, test.codigo)
			}

			err := repo.ActualizarEstado(ctx, test.codigo, "aprobado", time.Now())

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
		funcionCarga  func(context.Context, *codigoIDRepository)
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
			funcionCarga: func(ctx context.Context, r *codigoIDRepository) {
				r.Cargar(ctx, &models.CodigoID{
					Codigo:      randomString(6),
					Estado:      "pendiente",
					FechaPedido: time.Now(),
				})
			},
			largoEsperado: 1,
			errorEsperado: nil,
		},
	}
	ctx := t.Context()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("codigo_ID")
			defer db.Close()

			repo := NewCodigoIDRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(ctx, repo)
			}

			lista, err := repo.Listar(ctx, 10, 0)

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
