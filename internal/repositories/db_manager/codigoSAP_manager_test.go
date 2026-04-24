package dbmanager

import (
	"context"
	"errors"
	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"

	"testing"
)

func TestCargarCodigoSAP(t *testing.T) {

	descripPrueba := "description de prueba"

	tests := []struct {
		name          string
		codigo        models.CodigoSAP
		funcionCarga  func(context.Context, *codigoSAPRepository, *models.CodigoSAP)
		errorEsperado error
	}{
		{
			name: "Todo Ok",
			codigo: models.CodigoSAP{
				Codigo:      randomString(5),
				Descripcion: &descripPrueba,
			},
			errorEsperado: nil,
		},
		{
			name: "Código en blanco",
			codigo: models.CodigoSAP{
				Codigo:      "",
				Descripcion: &descripPrueba,
			},
			errorEsperado: appErrors.CodigoSAPVacio,
		},
		{
			name: "Código duplicado",
			codigo: models.CodigoSAP{
				Codigo:      randomString(5),
				Descripcion: &descripPrueba,
			},
			funcionCarga: func(ctx context.Context, cs *codigoSAPRepository, s *models.CodigoSAP) {
				cs.Cargar(ctx, s)
			},
			errorEsperado: appErrors.CodigoSAPDuplicado,
		},
		{
			name: "Description nil",
			codigo: models.CodigoSAP{
				Codigo:      randomString(5),
				Descripcion: nil,
			},
			errorEsperado: nil,
		},
	}
	ctx := t.Context()
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("codigo_SAP")
			defer db.Close()

			codigoSAP := NewCodigoSAPRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(ctx, codigoSAP, &test.codigo)
			}

			err := codigoSAP.Cargar(ctx, &test.codigo)

			if !errors.Is(err, test.errorEsperado) {

				t.Errorf("Error no esperado.\nSe esperaba: \n --- %v \nse obtuvo: \n --- %v",
					test.errorEsperado, err)
			}
		})
	}
}

func TestObtenerDetalleCodigoSAP(t *testing.T) {

	tests := []struct {
		name          string
		codigo        string
		funcionCarga  func(context.Context, *codigoSAPRepository, *models.CodigoSAP)
		errorEsperado error
	}{
		{
			name:   "Todo Ok",
			codigo: randomString(5),
			funcionCarga: func(ctx context.Context, cs *codigoSAPRepository, s *models.CodigoSAP) {
				cs.Cargar(ctx, s)
			},
			errorEsperado: nil,
		},
		{
			name:          "codigo no existe",
			codigo:        randomString(8),
			errorEsperado: appErrors.CodigoSAPNoEncontrado,
		},
	}
	ctx := t.Context()
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("codigo_SAP")
			defer db.Close()

			codigoSAPManager := NewCodigoSAPRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(ctx, codigoSAPManager, &models.CodigoSAP{
					Codigo: test.codigo,
				})
			}
			_, err := codigoSAPManager.ObtenerDetalle(ctx, test.codigo)

			if !errors.Is(err, test.errorEsperado) {

				t.Errorf("Error no esperado.\nSe esperaba: \n --- %v \nse obtuvo: \n --- %v",
					test.errorEsperado, err)
			}
		})
	}
}

// probar la búsqueda con filtro y sin filtro (parámetro en blanco, trae lista completa)
func TestBuscarPorDescripcionCodigoSAP(t *testing.T) {

	tests := []struct {
		name          string
		filtro        string
		funcionCarga  func(context.Context, *codigoSAPRepository)
		largoEsperado int
		errorEsperado error
	}{
		{
			name:          "Sin resultados",
			filtro:        "inexistente",
			largoEsperado: 0,
			errorEsperado: nil,
		},
		{
			name:   "Con resultados",
			filtro: "filtro",
			funcionCarga: func(ctx context.Context, r *codigoSAPRepository) {
				descripcionTest := "filtro"
				if err := r.Cargar(ctx,
					&models.CodigoSAP{
						Codigo:      randomString(4),
						Descripcion: &descripcionTest,
					}); err != nil {
					t.Errorf("Error al crear colaborador. Detalle: %v", err)
				}
				if err := r.Cargar(ctx,
					&models.CodigoSAP{
						Codigo:      randomString(4),
						Descripcion: &descripcionTest,
					}); err != nil {
					t.Errorf("Error al crear colaborador. Detalle: %v", err)
				}
			},
			largoEsperado: 2,
			errorEsperado: nil,
		},
	}

	ctx := t.Context()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("codigo_SAP")
			defer db.Close()

			repo := NewCodigoSAPRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(ctx, repo)
			}

			lista, err := repo.BuscarPorDescripcion(ctx, test.filtro)

			if !errors.Is(err, test.errorEsperado) {
				t.Errorf("Error inesperado.\nEsperado: %v\nObtenido: %v",
					test.errorEsperado, err)
			}

			if len(lista) != test.largoEsperado {
				t.Errorf("Cantidad incorrecta.\nEsperado: %d\nObtenido: %d",
					test.largoEsperado, len(lista))
			}
		})
	}

}

func TestModificarDescripcionCodigoSAP(t *testing.T) {

	tests := []struct {
		name                  string
		codigo                string
		descripcionModificada string
		funcionCarga          func(context.Context, *codigoSAPRepository, string)
		errorEsperado         error
	}{
		{
			name:                  "Todo Ok",
			codigo:                randomString(5),
			descripcionModificada: "description modificación de prueba",
			funcionCarga: func(ctx context.Context, cs *codigoSAPRepository, codigo string) {
				descripcionInicial := "description inicial de prueba"
				cs.Cargar(ctx, &models.CodigoSAP{
					Codigo:      codigo,
					Descripcion: &descripcionInicial,
				})
			},
			errorEsperado: nil,
		},
		{
			name:                  "Código en blanco",
			codigo:                "",
			descripcionModificada: "description modificación de prueba",
			errorEsperado:         appErrors.CodigoSAPVacio,
		},
		{
			name:                  "Código No encontrado",
			codigo:                randomString(5),
			descripcionModificada: "description modificación de prueba",
			errorEsperado:         appErrors.CodigoSAPNoEncontrado,
		},
		{
			name:   "Description inicial nil",
			codigo: randomString(5),
			funcionCarga: func(ctx context.Context, cs *codigoSAPRepository, codigo string) {
				cs.Cargar(ctx, &models.CodigoSAP{
					Codigo:      codigo,
					Descripcion: nil,
				})
			},
			descripcionModificada: "description modificación de prueba-nil",
			errorEsperado:         nil,
		},
	}

	ctx := t.Context()
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("codigo_SAP")
			defer db.Close()

			codigoSAP := NewCodigoSAPRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(ctx, codigoSAP, test.codigo)
			}

			err := codigoSAP.ModificarDescripcion(ctx, test.codigo, test.descripcionModificada)

			if !errors.Is(err, test.errorEsperado) {

				t.Errorf("Error no esperado.\nSe esperaba: \n --- %v \nse obtuvo: \n --- %v",
					test.errorEsperado, err)
			}
		})
	}
}

func TestListarCodigoSAP(t *testing.T) {

	tests := []struct {
		name          string
		limit         int
		offset        int
		funcionCarga  func(context.Context, *codigoSAPRepository)
		largoEsperado int
		errorEsperado error
	}{
		{
			name:          "Lista con datos",
			limit:         10,
			offset:        0,
			largoEsperado: 2,
			funcionCarga: func(ctx context.Context, sr *codigoSAPRepository) {
				sr.Cargar(ctx, &models.CodigoSAP{Codigo: randomString(6)})
				sr.Cargar(ctx, &models.CodigoSAP{Codigo: randomString(6)})
			},
			errorEsperado: nil,
		},
		{
			name:   "Respeta limit",
			limit:  1,
			offset: 0,
			funcionCarga: func(ctx context.Context, sr *codigoSAPRepository) {
				sr.Cargar(ctx, &models.CodigoSAP{Codigo: randomString(6)})
				sr.Cargar(ctx, &models.CodigoSAP{Codigo: randomString(6)})
			},
			largoEsperado: 1,
			errorEsperado: nil,
		},
	}

	ctx := t.Context()
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("codigo_SAP")
			defer db.Close()

			codigoSAPManager := NewCodigoSAPRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(ctx, codigoSAPManager)
			}

			lista, err := codigoSAPManager.Listar(ctx, test.limit, test.offset)

			if !errors.Is(err, test.errorEsperado) {
				t.Errorf("Error no esperado.\nSe esperaba:\n --- %v\nSe obtuvo:\n --- %v",
					test.errorEsperado, err)
			}

			if len(lista) != test.largoEsperado {
				t.Errorf("Cantidad incorrecta.\nSe esperaba:\n --- %v\nSe obtuvo:\n --- %v",
					test.largoEsperado, len(lista))
			}
		})
	}
}
