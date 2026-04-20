package dbmanager

import (
	"errors"
	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"
	"testing"
)

func TestDocumentoCargar(t *testing.T) {

	tests := []struct {
		name          string
		input         *models.Documento
		funcionCarga  func(*documentoRepository, *models.Documento)
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
			funcionCarga: func(r *documentoRepository, d *models.Documento) {
				r.Cargar(d)
			},
			errorEsperado: appErrors.DocumentoDuplicado,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("documento")
			defer db.Close()

			repo := NewDocumentoRepository(db)

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

func TestDocumentoObtenerDetalle(t *testing.T) {

	tests := []struct {
		name          string
		codigo        string
		funcionCarga  func(*documentoRepository, string)
		errorEsperado error
	}{
		{
			name:   "Todo OK",
			codigo: randomString(6),
			funcionCarga: func(r *documentoRepository, codigo string) {
				r.Cargar(&models.Documento{
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("documento")
			defer db.Close()

			repo := NewDocumentoRepository(db)

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

func TestDocumentoFiltrarPorTipo(t *testing.T) {

	tests := []struct {
		name          string
		tipo          string
		funcionCarga  func(*documentoRepository)
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
			funcionCarga: func(r *documentoRepository) {
				r.Cargar(&models.Documento{
					Codigo: randomString(6),
					Titulo: "Doc1",
					Tipo:   "pdf",
				})
				r.Cargar(&models.Documento{
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("documento")
			defer db.Close()

			repo := NewDocumentoRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(repo)
			}

			lista, err := repo.FiltrarPorTipo(test.tipo)

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
		funcionCarga  func(*documentoRepository)
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
			funcionCarga: func(r *documentoRepository) {
				r.Cargar(&models.Documento{
					Codigo: randomString(6),
					Titulo: "Documento 1",
				})
				r.Cargar(&models.Documento{
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("documento")
			defer db.Close()

			repo := NewDocumentoRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(repo)
			}

			lista, err := repo.FiltrarPorTitulo(test.titulo)

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
		funcionCarga  func(*documentoRepository, string)
		errorEsperado error
	}{
		{
			name:   "Todo OK",
			codigo: randomString(6),
			funcionCarga: func(r *documentoRepository, codigo string) {
				r.Cargar(&models.Documento{
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("documento")
			defer db.Close()

			repo := NewDocumentoRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(repo, test.codigo)
			}

			path := "/nuevo/path.pdf"
			err := repo.ActualizarPath(test.codigo, &path)

			if !errors.Is(err, test.errorEsperado) {
				t.Errorf("Error inesperado.\nEsperado: %v\nObtenido: %v",
					test.errorEsperado, err)
			}
		})
	}
}
