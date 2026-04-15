package dbmanager

import (
	"errors"
	"task-tracker-go/internal/appErrors"

	"testing"
)

func TestCrearColaborador(t *testing.T) {

	nombreTest := randomString(6)

	tests := []struct {
		name          string
		Usuario       string
		errorEsperado error
	}{
		{
			name:          "Todo Ok",
			Usuario:       nombreTest,
			errorEsperado: nil,
		},
		{
			name:          "Colaborador duplicado",
			Usuario:       nombreTest,
			errorEsperado: appErrors.ColaboradorDuplicado,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			db := NewGestotorDb("../../../database/app_test.db")
			defer db.Close()
			colaboradorManager := NewColaboradorRepository(db)

			_, err := colaboradorManager.Crear(test.Usuario)

			if !errors.Is(err, test.errorEsperado) {

				t.Errorf("Error no esperado.\nSe esperaba: \n --- %v \nse obtuvo: \n --- %v",
					test.errorEsperado, err)
			}
		})
	}
}

func TestObtenerIDPorNombreColaborador(t *testing.T) {
	tests := []struct {
		name          string
		Usuario       string
		errorEsperado error
	}{
		{
			name:          "Todo Ok",
			Usuario:       "ColaboradorTestExistente",
			errorEsperado: nil,
		},
		{
			name:          "Colaborador no existe",
			Usuario:       randomString(8),
			errorEsperado: appErrors.ColaboradorNoEncontrado,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			db := NewGestotorDb("../../../database/app_test.db")
			defer db.Close()
			colaboradorManager := NewColaboradorRepository(db)

			_, err := colaboradorManager.ObtenerIDPorNombre(test.Usuario)

			if !errors.Is(err, test.errorEsperado) {

				t.Errorf("Error no esperado.\nSe esperaba: \n --- %v \nse obtuvo: \n --- %v",
					test.errorEsperado, err)
			}
		})
	}
}

// probar la búsqueda con filtro y sin filtro (parámetro en blanco, trae lista completa)
func TestBuscarColaborador(t *testing.T) {

	db := NewGestotorDb("../../../database/app_test.db")
	defer db.Close()
	colaboradorManager := NewColaboradorRepository(db)

	// creamos 3 colaboradores para realizar la prueba
	if _, err := colaboradorManager.Crear("ColaboradorFiltro1"); err != nil {
		t.Errorf("Error al crear colaborador. Detalle: %v", err)
	}
	if _, err := colaboradorManager.Crear("ColaboradorFiltro2"); err != nil {
		t.Errorf("Error al crear colaborador. Detalle: %v", err)
	}
	if _, err := colaboradorManager.Crear("ColaboradorFiltro3"); err != nil {
		t.Errorf("Error al crear colaborador. Detalle: %v", err)
	}

	// Busco filtrando
	listaFiltrada, err := colaboradorManager.Buscar("Filtro")
	if err != nil {
		t.Errorf("Error no esperado.\nSe esperaba: \n --- nil \nse obtuvo: \n --- %v", err)
	}

	if len(listaFiltrada) != 3 {
		t.Errorf("Error en la busqueda, se esperaban 3 resultados pero se obtuvo %d. Detalle: \n %v", len(listaFiltrada), listaFiltrada)
	}

	// Busco toda la lista
	listaSinFiltrar, err := colaboradorManager.Buscar("")
	if err != nil {
		t.Errorf("Error no esperado.\nSe esperaba: \n --- nil \nse obtuvo: \n --- %v", err)
	}

	if len(listaFiltrada) >= len(listaSinFiltrar) {
		t.Errorf("Error en la busqueda, se esperaban más resultados sin filtro pero se obtuvo %d", len(listaSinFiltrar))
	}
}
