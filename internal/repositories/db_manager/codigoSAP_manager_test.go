package dbmanager

import (
	"errors"
	"task-tracker-go/internal/appErrors"

	"testing"
)

func TestCargarCodigoSAP(t *testing.T) {

	descripPrueba := "description de prueba"
	codigoTest := randomString(5)

	tests := []struct {
		name          string
		codigo        string
		descripcion   *string
		errorEsperado error
	}{
		{
			name:          "Todo Ok",
			codigo:        codigoTest,
			descripcion:   &descripPrueba,
			errorEsperado: nil,
		},
		{
			name:          "Código en blanco",
			codigo:        "",
			descripcion:   &descripPrueba,
			errorEsperado: appErrors.CodigoSAPVacio,
		},
		{
			name:          "Código duplicado",
			codigo:        codigoTest,
			descripcion:   &descripPrueba,
			errorEsperado: appErrors.CodigoSAPDuplicado,
		},
		{
			name:          "Description en blanco",
			codigo:        randomString(5),
			descripcion:   nil,
			errorEsperado: nil,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			db := NewGestotorDb("../../../database/app_test.db")
			defer db.Close()
			codigoSAP := NewCodigoSAPRepository(db)

			err := codigoSAP.Cargar(test.codigo, test.descripcion)

			if !errors.Is(err, test.errorEsperado) {

				t.Errorf("Error no esperado.\nSe esperaba: \n --- %v \nse obtuvo: \n --- %v",
					test.errorEsperado, err)
			}
		})
	}
}

func TestObtenerDetalleCodigoSAP(t *testing.T) {
	codigoTest := randomString(5)
	tests := []struct {
		name          string
		codigo        string
		errorEsperado error
	}{
		{
			name:          "Todo Ok",
			codigo:        codigoTest,
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

			codigoSAPManager.Cargar(codigoTest, nil)
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
	descipcionTest := "DescripcionFiltro1"
	if err := codigoSAPManager.Cargar("codigo1", &descipcionTest); err != nil {
		t.Errorf("Error al crear colaborador. Detalle: %v", err)
	}
	descipcionTest = "DescripcionFiltro2"
	if err := codigoSAPManager.Cargar("codigo2", &descipcionTest); err != nil {
		t.Errorf("Error al crear colaborador. Detalle: %v", err)
	}
	descipcionTest = "DescripcionFiltro3"
	if err := codigoSAPManager.Cargar("codigo3", &descipcionTest); err != nil {
		t.Errorf("Error al crear colaborador. Detalle: %v", err)
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

	tests := []struct {
		name                  string
		codigo                string
		descripcionModificada string
		funcionCarga          func(*codigoSAPRepository, string)
		errorEsperado         error
	}{
		{
			name:                  "Todo Ok",
			codigo:                randomString(5),
			descripcionModificada: "description modificación de prueba",
			funcionCarga: func(cs *codigoSAPRepository, codigo string) {
				descripcionInicial := "description inicial de prueba"
				cs.Cargar(codigo, &descripcionInicial)
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
			funcionCarga: func(cs *codigoSAPRepository, codigo string) {
				cs.Cargar(codigo, nil)
			},
			descripcionModificada: "description modificación de prueba-nil",
			errorEsperado:         nil,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			db := NewGestotorDb("../../../database/app_test.db")
			defer db.Close()
			codigoSAP := NewCodigoSAPRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(codigoSAP, test.codigo)
			}

			err := codigoSAP.ModificarDescripcion(test.codigo, test.descripcionModificada)

			if !errors.Is(err, test.errorEsperado) {

				t.Errorf("Error no esperado.\nSe esperaba: \n --- %v \nse obtuvo: \n --- %v",
					test.errorEsperado, err)
			}
		})
	}
}
