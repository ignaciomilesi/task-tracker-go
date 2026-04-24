package dbmanager

import (
	"context"
	"errors"
	"task-tracker-go/internal/appErrors"
	"task-tracker-go/internal/models"

	"testing"
)

func TestCrearColaborador(t *testing.T) {

	tests := []struct {
		name          string
		colaborador   *models.Colaborador
		funcionCarga  func(context.Context, *colaboradorRepository, *models.Colaborador)
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
			funcionCarga: func(ctx context.Context, cr *colaboradorRepository, c *models.Colaborador) {
				cr.Crear(ctx, c)
			},
			errorEsperado: appErrors.ColaboradorDuplicado,
		},
	}

	ctx := t.Context()
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("colaborador")
			defer db.Close()

			colaboradorManager := NewColaboradorRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(ctx, colaboradorManager, test.colaborador)
			}

			_, err := colaboradorManager.Crear(ctx, test.colaborador)

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
		funcionCarga  func(context.Context, *colaboradorRepository, *models.Colaborador)
		errorEsperado error
	}{
		{
			name:    "Todo Ok",
			Usuario: randomString(8),
			funcionCarga: func(ctx context.Context, sr *colaboradorRepository, s *models.Colaborador) {
				sr.Crear(ctx, s)
			},
			errorEsperado: nil,
		},
		{
			name:          "Colaborador no existe",
			Usuario:       randomString(8),
			errorEsperado: appErrors.ColaboradorNoEncontrado,
		},
	}

	ctx := t.Context()
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("colaborador")
			defer db.Close()

			colaboradorManager := NewColaboradorRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(ctx, colaboradorManager, &models.Colaborador{
					Nombre: test.Usuario,
				})
			}

			_, err := colaboradorManager.ObtenerIDPorNombre(ctx, test.Usuario)

			if !errors.Is(err, test.errorEsperado) {

				t.Errorf("Error no esperado.\nSe esperaba: \n --- %v \nse obtuvo: \n --- %v",
					test.errorEsperado, err)
			}
		})
	}
}

func TestBuscarColaborador(t *testing.T) {

	tests := []struct {
		name          string
		filtro        string
		funcionCarga  func(context.Context, *colaboradorRepository)
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
			funcionCarga: func(ctx context.Context, r *colaboradorRepository) {
				if _, err := r.Crear(ctx,
					&models.Colaborador{
						Nombre: "ColaboradorFiltro1",
					}); err != nil {
					t.Errorf("Error al crear colaborador. Detalle: %v", err)
				}
				if _, err := r.Crear(ctx,
					&models.Colaborador{
						Nombre: "ColaboradorFiltro2",
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

			db := generarGestorDbLimpio("colaborador")
			defer db.Close()

			repo := NewColaboradorRepository(db)

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

func TestListarColaborador(t *testing.T) {

	tests := []struct {
		name          string
		limit         int
		offset        int
		funcionCarga  func(context.Context, *colaboradorRepository)
		largoEsperado int
		errorEsperado error
	}{
		{
			name:   "Lista con datos",
			limit:  10,
			offset: 0,
			funcionCarga: func(ctx context.Context, sr *colaboradorRepository) {
				sr.Crear(ctx, &models.Colaborador{Nombre: randomString(6)})
				sr.Crear(ctx, &models.Colaborador{Nombre: randomString(6)})
			},
			largoEsperado: 2,
			errorEsperado: nil,
		},
		{
			name:   "Respeta limit",
			limit:  1,
			offset: 0,
			funcionCarga: func(ctx context.Context, sr *colaboradorRepository) {
				sr.Crear(ctx, &models.Colaborador{Nombre: randomString(6)})
				sr.Crear(ctx, &models.Colaborador{Nombre: randomString(6)})
			},
			largoEsperado: 1,
			errorEsperado: nil,
		},
	}

	ctx := t.Context()
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("colaborador")
			defer db.Close()

			colaboradorManager := NewColaboradorRepository(db)

			if test.funcionCarga != nil {
				test.funcionCarga(ctx, colaboradorManager)
			}

			lista, err := colaboradorManager.Listar(ctx, test.limit, test.offset)

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
