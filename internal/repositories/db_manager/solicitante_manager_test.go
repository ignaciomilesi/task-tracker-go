package dbmanager

import (
	"context"
	"errors"
	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"

	"testing"
)

func TestCrear(t *testing.T) {

	tests := []struct {
		name          string
		solicitante   *models.Solicitante
		funcionCarga  func(context.Context, *solicitanteRepository, *models.Solicitante)
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
			funcionCarga: func(ctx context.Context, sr *solicitanteRepository, solicitante *models.Solicitante) {
				sr.Crear(ctx, solicitante)
			},
			errorEsperado: appErrors.SolicitanteDuplicado,
		},
	}

	ctx := t.Context()
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("solicitante")
			defer db.Close()

			solicitanteManager := NewsolicitanteRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(ctx, solicitanteManager, test.solicitante)
			}
			_, err := solicitanteManager.Crear(ctx, test.solicitante)

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
		funcionCarga  func(context.Context, *solicitanteRepository, *models.Solicitante)
		errorEsperado error
	}{
		{
			name:    "Todo Ok",
			Usuario: randomString(8),
			funcionCarga: func(ctx context.Context, sr *solicitanteRepository, s *models.Solicitante) {
				sr.Crear(ctx, s)
			},
			errorEsperado: nil,
		},
		{
			name:          "Solicitante no existe",
			Usuario:       randomString(8),
			errorEsperado: appErrors.SolicitanteNoEncontrado,
		},
	}

	ctx := t.Context()
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("solicitante")
			defer db.Close()

			solicitanteManager := NewsolicitanteRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(ctx, solicitanteManager, &models.Solicitante{
					Nombre: test.Usuario,
				})
			}

			_, err := solicitanteManager.ObtenerIDPorNombre(ctx, test.Usuario)

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
		funcionCarga  func(context.Context, *solicitanteRepository)
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
			funcionCarga: func(ctx context.Context, r *solicitanteRepository) {
				if _, err := r.Crear(ctx,
					&models.Solicitante{
						Nombre: "SolicitanteFiltro1",
					}); err != nil {
					t.Errorf("Error al crear solicitante. Detalle: %v", err)
				}
				if _, err := r.Crear(ctx,
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

	ctx := t.Context()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("solicitante")
			defer db.Close()

			repo := NewsolicitanteRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(ctx, repo)
			}

			lista, err := repo.Buscar(ctx, test.filtro)

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
		funcionCarga  func(context.Context, *solicitanteRepository)
		largoEsperado int
		errorEsperado error
	}{
		{
			name:   "Lista con datos",
			limit:  10,
			offset: 0,
			funcionCarga: func(ctx context.Context, sr *solicitanteRepository) {
				sr.Crear(ctx, &models.Solicitante{Nombre: randomString(6)})
				sr.Crear(ctx, &models.Solicitante{Nombre: randomString(6)})
			},
			errorEsperado: nil,
		},
		{
			name:   "Respeta limit",
			limit:  1,
			offset: 0,
			funcionCarga: func(ctx context.Context, sr *solicitanteRepository) {
				sr.Crear(ctx, &models.Solicitante{Nombre: randomString(6)})
				sr.Crear(ctx, &models.Solicitante{Nombre: randomString(6)})
			},
			largoEsperado: 1,
			errorEsperado: nil,
		},
	}

	ctx := t.Context()
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("solicitante")
			defer db.Close()

			solicitanteManager := NewsolicitanteRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(ctx, solicitanteManager)
			}

			lista, err := solicitanteManager.Listar(ctx, test.limit, test.offset)

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
