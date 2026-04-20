package dbmanager

import (
	"errors"
	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"

	"testing"
)

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

			db := generarGestorDbLimpio("solicitante")
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

			db := generarGestorDbLimpio("solicitante")
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

// probar la búsqueda con filtro
func TestBuscar(t *testing.T) {

	tests := []struct {
		name          string
		filtro        string
		funcionCarga  func(*solicitanteRepository)
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
			funcionCarga: func(r *solicitanteRepository) {
				if _, err := r.Crear(
					&models.Solicitante{
						Nombre: "SolicitanteFiltro1",
					}); err != nil {
					t.Errorf("Error al crear solicitante. Detalle: %v", err)
				}
				if _, err := r.Crear(
					&models.Solicitante{
						Nombre: "SolicitanteFiltro2",
					}); err != nil {
					t.Errorf("Error al crear solicitante. Detalle: %v", err)
				}
			},
			largoEsperado: 2,
			errorEsperado: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("solicitante")
			defer db.Close()

			repo := NewsolicitanteRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(repo)
			}

			lista, err := repo.Buscar(test.filtro)

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

func TestListar(t *testing.T) {

	tests := []struct {
		name          string
		limit         int
		offset        int
		funcionCarga  func(*solicitanteRepository)
		largoEsperado int
		errorEsperado error
	}{
		{
			name:   "Lista con datos",
			limit:  10,
			offset: 0,
			funcionCarga: func(sr *solicitanteRepository) {
				sr.Crear(&models.Solicitante{Nombre: randomString(6)})
				sr.Crear(&models.Solicitante{Nombre: randomString(6)})
			},
			errorEsperado: nil,
		},
		{
			name:   "Respeta limit",
			limit:  1,
			offset: 0,
			funcionCarga: func(sr *solicitanteRepository) {
				sr.Crear(&models.Solicitante{Nombre: randomString(6)})
				sr.Crear(&models.Solicitante{Nombre: randomString(6)})
			},
			largoEsperado: 1,
			errorEsperado: nil,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("solicitante")
			defer db.Close()

			solicitanteManager := NewsolicitanteRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(solicitanteManager)
			}

			lista, err := solicitanteManager.Listar(test.limit, test.offset)

			if !errors.Is(err, test.errorEsperado) {
				t.Errorf("Error no esperado.\nSe esperaba:\n --- %v\nSe obtuvo:\n --- %v",
					test.errorEsperado, err)
			}

			if test.largoEsperado != 0 && len(lista) != test.largoEsperado {
				t.Errorf("Cantidad incorrecta.\nSe esperaba:\n --- %v\nSe obtuvo:\n --- %v",
					test.largoEsperado, len(lista))
			}
		})
	}
}
