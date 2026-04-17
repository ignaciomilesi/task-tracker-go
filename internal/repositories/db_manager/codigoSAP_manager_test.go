package dbmanager

import (
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
		funcionCarga  func(*codigoSAPRepository, *models.CodigoSAP)
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
			funcionCarga: func(cs *codigoSAPRepository, s *models.CodigoSAP) {
				cs.Cargar(s)
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

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			db := NewGestotorDb("../../../database/app_test.db")
			defer db.Close()
			codigoSAP := NewCodigoSAPRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(codigoSAP, &test.codigo)
			}

			err := codigoSAP.Cargar(&test.codigo)

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
		funcionCarga  func(*codigoSAPRepository, *models.CodigoSAP)
		errorEsperado error
	}{
		{
			name:   "Todo Ok",
			codigo: randomString(5),
			funcionCarga: func(cs *codigoSAPRepository, s *models.CodigoSAP) {
				cs.Cargar(s)
			},
			errorEsperado: nil,
		},
		{
			name:          "codigo no existe",
			codigo:        randomString(8),
			errorEsperado: appErrors.CodigoSAPNoEncontrado,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			db := NewGestotorDb("../../../database/app_test.db")
			defer db.Close()
			codigoSAPManager := NewCodigoSAPRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(codigoSAPManager, &models.CodigoSAP{
					Codigo: test.codigo,
				})
			}
			_, err := codigoSAPManager.ObtenerDetalle(test.codigo)

			if !errors.Is(err, test.errorEsperado) {

				t.Errorf("Error no esperado.\nSe esperaba: \n --- %v \nse obtuvo: \n --- %v",
					test.errorEsperado, err)
			}
		})
	}
}

// probar la búsqueda con filtro y sin filtro (parámetro en blanco, trae lista completa)
func TestBuscarPorDescripcionCodigoSAP(t *testing.T) {

	db := NewGestotorDb("../../../database/app_test.db")
	defer db.Close()
	codigoSAPManager := NewCodigoSAPRepository(db)

	// creamos 3 códigos para realizar la prueba

	descripcionTest := "filtro"
	if err := codigoSAPManager.Cargar(&models.CodigoSAP{
		Codigo:      randomString(4),
		Descripcion: &descripcionTest,
	}); err != nil {
		t.Errorf("Error al crear codigoSAP. Detalle: %v", err)
	}

	if err := codigoSAPManager.Cargar(&models.CodigoSAP{
		Codigo:      randomString(4),
		Descripcion: &descripcionTest,
	}); err != nil {
		t.Errorf("Error al crear codigoSAP. Detalle: %v", err)
	}

	if err := codigoSAPManager.Cargar(&models.CodigoSAP{
		Codigo:      randomString(4),
		Descripcion: &descripcionTest,
	}); err != nil {
		t.Errorf("Error al crear codigoSAP. Detalle: %v", err)
	}

	// Busco filtrando
	listaFiltrada, err := codigoSAPManager.BuscarPorDescripcion("Filtro")
	if err != nil {
		t.Errorf("Error no esperado.\nSe esperaba: \n --- nil \nse obtuvo: \n --- %v", err)
	}

	if len(listaFiltrada) != 3 {
		t.Errorf("Error en la búsqueda, se esperaban 3 resultados pero se obtuvo %d. Detalle: \n %v", len(listaFiltrada), listaFiltrada)
	}

	// Busco toda la lista
	listaSinFiltrar, err := codigoSAPManager.BuscarPorDescripcion("")
	if err != nil {
		t.Errorf("Error no esperado.\nSe esperaba: \n --- nil \nse obtuvo: \n --- %v", err)
	}

	if len(listaFiltrada) >= len(listaSinFiltrar) {
		t.Errorf("Error en la búsqueda, se esperaban más resultados sin filtro pero se obtuvo %d", len(listaSinFiltrar))
	}
}

func TestModificarDescripcionCodigoSAP(t *testing.T) {

	descripcionModificada := "description modificación de prueba"

	tests := []struct {
		name             string
		codigoModificado *models.CodigoSAP
		funcionCarga     func(*codigoSAPRepository, string)
		errorEsperado    error
	}{
		{
			name: "Todo Ok",
			codigoModificado: &models.CodigoSAP{
				Codigo:      randomString(5),
				Descripcion: &descripcionModificada,
			},
			funcionCarga: func(cs *codigoSAPRepository, codigo string) {
				descripcionInicial := "description inicial de prueba"
				cs.Cargar(&models.CodigoSAP{
					Codigo:      codigo,
					Descripcion: &descripcionInicial,
				})
			},
			errorEsperado: nil,
		},
		{
			name:             "Código en blanco",
			codigoModificado: &models.CodigoSAP{},
			errorEsperado:    appErrors.CodigoSAPVacio,
		},
		{
			name: "Código No encontrado",
			codigoModificado: &models.CodigoSAP{
				Codigo:      randomString(5),
				Descripcion: &descripcionModificada,
			},
			errorEsperado: appErrors.CodigoSAPNoEncontrado,
		},
		{
			name: "Description inicial nil",
			codigoModificado: &models.CodigoSAP{
				Codigo:      randomString(5),
				Descripcion: &descripcionModificada,
			},
			funcionCarga: func(cs *codigoSAPRepository, codigo string) {
				cs.Cargar(&models.CodigoSAP{
					Codigo: codigo,
				})
			},
			errorEsperado: nil,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			db := NewGestotorDb("../../../database/app_test.db")
			defer db.Close()
			codigoSAP := NewCodigoSAPRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(codigoSAP, test.codigoModificado.Codigo)
			}

			err := codigoSAP.ModificarCodigoSAP(test.codigoModificado)

			if !errors.Is(err, test.errorEsperado) {

				t.Errorf("Error no esperado.\nSe esperaba: \n --- %v \nse obtuvo: \n --- %v",
					test.errorEsperado, err)
			}
		})
	}
}
