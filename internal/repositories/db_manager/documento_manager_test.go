package dbmanager

import (
	"context"
	"errors"
	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"
	"testing"
)

func TestDocumentoCargar(t *testing.T) {

	tests := []struct {
		name          string
		input         *models.Documento
		funcionCarga  func(context.Context, *documentoRepository, *models.Documento)
		errorEsperado error
	}{
		{
			name: "Todo OK",
			input: &models.Documento{
				Codigo: randomString(6),
				Titulo: "Doc test",
				Tipo:   "pdf",
			},
			errorEsperado: nil,
		},
		{
			name: "Codigo duplicado",
			input: &models.Documento{
				Codigo: randomString(6),
				Titulo: "Doc test",
				Tipo:   "pdf",
			},
			funcionCarga: func(ctx context.Context, r *documentoRepository, d *models.Documento) {
				r.Cargar(ctx, d)
			},
			errorEsperado: appErrors.DocumentoDuplicado,
		},
	}

	ctx := t.Context()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("documento")
			defer db.Close()

			repo := NewDocumentoRepository(db)

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

func TestDocumentoObtenerDetalle(t *testing.T) {

	tests := []struct {
		name          string
		codigo        string
		funcionCarga  func(context.Context, *documentoRepository, string)
		errorEsperado error
	}{
		{
			name:   "Todo OK",
			codigo: randomString(6),
			funcionCarga: func(ctx context.Context, r *documentoRepository, codigo string) {
				r.Cargar(ctx, &models.Documento{
					Codigo: codigo,
					Titulo: "Doc",
					Tipo:   "pdf",
				})
			},
			errorEsperado: nil,
		},
		{
			name:          "Codigo vacío",
			codigo:        " ",
			errorEsperado: appErrors.ParametroDeBusquedaVacio,
		},
		{
			name:          "No encontrado",
			codigo:        "no_existe",
			errorEsperado: appErrors.DocumentoNoEncontrado,
		},
	}

	ctx := t.Context()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("documento")
			defer db.Close()

			repo := NewDocumentoRepository(db)

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

func TestDocumentoFiltrarPorTipo(t *testing.T) {

	tests := []struct {
		name          string
		tipo          string
		funcionCarga  func(context.Context, *documentoRepository)
		largoEsperado int
		errorEsperado error
	}{
		{
			name:          "Sin resultados",
			tipo:          "inexistente",
			largoEsperado: 0,
			errorEsperado: nil,
		},
		{
			name: "Con resultados",
			tipo: "pdf",
			funcionCarga: func(ctx context.Context, r *documentoRepository) {
				r.Cargar(ctx, &models.Documento{
					Codigo: randomString(6),
					Titulo: "Doc1",
					Tipo:   "pdf",
				})
				r.Cargar(ctx, &models.Documento{
					Codigo: randomString(6),
					Titulo: "Doc2",
					Tipo:   "pdf",
				})
			},
			largoEsperado: 2,
			errorEsperado: nil,
		},
		{
			name:          "Tipo vacío",
			tipo:          " ",
			errorEsperado: appErrors.ParametroDeBusquedaVacio,
		},
	}

	ctx := t.Context()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("documento")
			defer db.Close()

			repo := NewDocumentoRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(ctx, repo)
			}

			lista, err := repo.FiltrarPorTipo(ctx, test.tipo)

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

func TestDocumentoFiltrarPorTitulo(t *testing.T) {

	tests := []struct {
		name          string
		titulo        string
		funcionCarga  func(context.Context, *documentoRepository)
		largoEsperado int
		errorEsperado error
	}{
		{
			name:          "Sin resultados",
			titulo:        "zzz",
			largoEsperado: 0,
			errorEsperado: nil,
		},
		{
			name:   "Con resultados",
			titulo: "Doc",
			funcionCarga: func(ctx context.Context, r *documentoRepository) {
				r.Cargar(ctx, &models.Documento{
					Codigo: randomString(6),
					Titulo: "Documento 1",
				})
				r.Cargar(ctx, &models.Documento{
					Codigo: randomString(6),
					Titulo: "Documento 2",
				})
			},
			largoEsperado: 2,
			errorEsperado: nil,
		},
		{
			name:          "Titulo vacío",
			titulo:        " ",
			errorEsperado: appErrors.ParametroDeBusquedaVacio,
		},
	}

	ctx := t.Context()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("documento")
			defer db.Close()

			repo := NewDocumentoRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(ctx, repo)
			}

			lista, err := repo.FiltrarPorTitulo(ctx, test.titulo)

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

func TestDocumentoActualizarPath(t *testing.T) {

	tests := []struct {
		name          string
		codigo        string
		funcionCarga  func(context.Context, *documentoRepository, string)
		errorEsperado error
	}{
		{
			name:   "Todo OK",
			codigo: randomString(6),
			funcionCarga: func(ctx context.Context, r *documentoRepository, codigo string) {
				r.Cargar(ctx, &models.Documento{
					Codigo: codigo,
				})
			},
			errorEsperado: nil,
		},
		{
			name:          "Codigo vacío",
			codigo:        " ",
			errorEsperado: appErrors.ParametroDeBusquedaVacio,
		},
		{
			name:          "No encontrado",
			codigo:        "no_existe",
			errorEsperado: appErrors.DocumentoNoEncontrado,
		},
	}

	ctx := t.Context()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("documento")
			defer db.Close()

			repo := NewDocumentoRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(ctx, repo, test.codigo)
			}

			err := repo.ActualizarPath(ctx, test.codigo, "/nuevo/path.pdf")

			if !errors.Is(err, test.errorEsperado) {
				t.Errorf("Error inesperado.\nEsperado: %v\nObtenido: %v",
					test.errorEsperado, err)
			}
		})
	}
}
