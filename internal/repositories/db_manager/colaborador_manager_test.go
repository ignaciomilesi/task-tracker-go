package dbmanager

import (
	"errors"
	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"

	"testing"
)

func TestCrearColaborador(t *testing.T) {

	tests := []struct {
		name          string
		colaborador   *models.Colaborador
		funcionCarga  func(*colaboradorRepository, *models.Colaborador)
		errorEsperado error
	}{
		{
			name: "Todo Ok",
			colaborador: &models.Colaborador{
				Nombre: randomString(6),
			},
			errorEsperado: nil,
		},
		{
			name: "Colaborador duplicado",
			colaborador: &models.Colaborador{
				Nombre: randomString(6),
			},
			funcionCarga: func(cr *colaboradorRepository, c *models.Colaborador) {
				cr.Crear(c)
			},
			errorEsperado: appErrors.ColaboradorDuplicado,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			db := NewGestotorDb("../../../database/app_test.db")
			defer db.Close()
			colaboradorManager := NewColaboradorRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(colaboradorManager, test.colaborador)
			}

			_, err := colaboradorManager.Crear(test.colaborador)

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
		funcionCarga  func(*colaboradorRepository, *models.Colaborador)
		errorEsperado error
	}{
		{
			name:    "Todo Ok",
			Usuario: randomString(8),
			funcionCarga: func(sr *colaboradorRepository, s *models.Colaborador) {
				sr.Crear(s)
			},
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

			if test.funcionCarga != nil {
				test.funcionCarga(colaboradorManager, &models.Colaborador{
					Nombre: test.Usuario,
				})
			}

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
	if _, err := colaboradorManager.Crear(
		&models.Colaborador{
			Nombre: "SolicitanteFiltro1",
		}); err != nil {
		t.Errorf("Error al crear colaborador. Detalle: %v", err)
	}
	if _, err := colaboradorManager.Crear(
		&models.Colaborador{
			Nombre: "ColaboradorFiltro2",
		}); err != nil {
		t.Errorf("Error al crear colaborador. Detalle: %v", err)
	}
	if _, err := colaboradorManager.Crear(
		&models.Colaborador{
			Nombre: "ColaboradorFiltro3",
		}); err != nil {
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
