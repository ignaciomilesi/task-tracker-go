package dbmanager

import (
	"errors"
	"math/rand"
	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"

	"testing"
)

func randomString(n int) string {

	const letras = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, n)
	for i := range b {
		b[i] = letras[rand.Intn(len(letras))]
	}
	return string(b)
}

func TestCrear(t *testing.T) {

	tests := []struct {
		name          string
		solicitante   *models.Solicitante
		funcionCarga  func(*solicitanteRepository, *models.Solicitante)
		errorEsperado error
	}{
		{
			name: "Todo Ok",
			solicitante: &models.Solicitante{
				Nombre: randomString(6),
			},
			errorEsperado: nil,
		},
		{
			name: "Solicitante duplicado",
			solicitante: &models.Solicitante{
				Nombre: randomString(6),
			},
			funcionCarga: func(sr *solicitanteRepository, solicitante *models.Solicitante) {
				sr.Crear(solicitante)
			},
			errorEsperado: appErrors.SolicitanteDuplicado,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			db := NewGestotorDb("../../../database/app_test.db")
			defer db.Close()
			solicitanteManager := NewsolicitanteRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(solicitanteManager, test.solicitante)
			}
			_, err := solicitanteManager.Crear(test.solicitante)

			if !errors.Is(err, test.errorEsperado) {

				t.Errorf("Error no esperado.\nSe esperaba: \n --- %v \nse obtuvo: \n --- %v",
					test.errorEsperado, err)
			}
		})
	}

}

func TestObtenerIDPorNombre(t *testing.T) {
	tests := []struct {
		name          string
		Usuario       string
		funcionCarga  func(*solicitanteRepository, *models.Solicitante)
		errorEsperado error
	}{
		{
			name:    "Todo Ok",
			Usuario: randomString(8),
			funcionCarga: func(sr *solicitanteRepository, s *models.Solicitante) {
				sr.Crear(s)
			},
			errorEsperado: nil,
		},
		{
			name:          "Solicitante no existe",
			Usuario:       randomString(8),
			errorEsperado: appErrors.SolicitanteNoEncontrado,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			db := NewGestotorDb("../../../database/app_test.db")
			defer db.Close()
			solicitanteManager := NewsolicitanteRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(solicitanteManager, &models.Solicitante{
					Nombre: test.Usuario,
				})
			}

			_, err := solicitanteManager.ObtenerIDPorNombre(test.Usuario)

			if !errors.Is(err, test.errorEsperado) {

				t.Errorf("Error no esperado.\nSe esperaba: \n --- %v \nse obtuvo: \n --- %v",
					test.errorEsperado, err)
			}
		})
	}
}

// probar la búsqueda con filtro y sin filtro (parámetro en blanco, trae lista completa)
func TestBuscar(t *testing.T) {

	db := NewGestotorDb("../../../database/app_test.db")
	defer db.Close()
	solicitanteManager := NewsolicitanteRepository(db)

	// creamos 3 solicitantes para realizar la prueba
	if _, err := solicitanteManager.Crear(
		&models.Solicitante{
			Nombre: "SolicitanteFiltro1",
		}); err != nil {
		t.Errorf("Error al crear solicitante. Detalle: %v", err)
	}
	if _, err := solicitanteManager.Crear(
		&models.Solicitante{
			Nombre: "SolicitanteFiltro2",
		}); err != nil {
		t.Errorf("Error al crear solicitante. Detalle: %v", err)
	}
	if _, err := solicitanteManager.Crear(
		&models.Solicitante{
			Nombre: "SolicitanteFiltro3",
		}); err != nil {
		t.Errorf("Error al crear solicitante. Detalle: %v", err)
	}

	// Busco filtrando
	listaFiltrada, err := solicitanteManager.Buscar("Filtro")
	if err != nil {
		t.Errorf("Error no esperado.\nSe esperaba: \n --- nil \nse obtuvo: \n --- %v", err)
	}

	if len(listaFiltrada) != 3 {
		t.Errorf("Error en la busqueda, se esperaban 3 resultados pero se obtuvo %d. Detalle: \n %v", len(listaFiltrada), listaFiltrada)
	}

	// Busco toda la lista
	listaSinFiltrar, err := solicitanteManager.Buscar("")
	if err != nil {
		t.Errorf("Error no esperado.\nSe esperaba: \n --- nil \nse obtuvo: \n --- %v", err)
	}

	if len(listaFiltrada) >= len(listaSinFiltrar) {
		t.Errorf("Error en la busqueda, se esperaban 3 resultados pero se obtuvo %d", len(listaSinFiltrar))
	}

}
