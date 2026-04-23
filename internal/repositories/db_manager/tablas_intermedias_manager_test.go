package dbmanager

import (
	"database/sql"
	"errors"
	"task-tracker-go/internal/appErrors"
	"testing"
)

func TestVincularPendienteDocumento(t *testing.T) {

	tests := []struct {
		name          string
		pendienteID   int
		documentoID   int
		setup         func(*sql.DB)
		errorEsperado error
	}{
		{
			name:        "Todo OK",
			pendienteID: 1,
			documentoID: 1,
			setup: func(db *sql.DB) {
				_, _ = db.Exec("PRAGMA foreign_keys = OFF")
			},
			errorEsperado: nil,
		},
		{
			name:          "FK inexistente",
			pendienteID:   999,
			documentoID:   999,
			errorEsperado: appErrors.FkNoEncontrado,
		},
		{
			name:        "Relacion duplicada",
			pendienteID: 1,
			documentoID: 1,
			setup: func(db *sql.DB) {
				_, _ = db.Exec("PRAGMA foreign_keys = OFF")
				_, _ = db.Exec(`INSERT INTO ti_pendientes_documento VALUES (1,1)`)
			},
			errorEsperado: appErrors.RelacionDuplicada,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("ti_pendientes_documento")
			defer db.Close()

			if test.setup != nil {
				test.setup(db)
			}

			repo := NewTablaIntermediaRepository(db)

			err := repo.VincularPendienteDocumento(test.pendienteID, test.documentoID)

			if !errors.Is(err, test.errorEsperado) {
				t.Errorf("Esperado: %v\nObtenido: %v", test.errorEsperado, err)
			}
		})
	}
}

func TestListarDocumentosPorPendiente(t *testing.T) {

	tests := []struct {
		name          string
		setup         func(*sql.DB)
		largoEsperado int
	}{
		{
			name:          "Sin resultados",
			largoEsperado: 0,
		},
		{
			name: "Con resultados",
			setup: func(db *sql.DB) {
				_, _ = db.Exec("PRAGMA foreign_keys = OFF")
				_, _ = db.Exec(`INSERT INTO ti_pendientes_documento VALUES (1,1)`)
				_, _ = db.Exec(`INSERT INTO ti_pendientes_documento VALUES (1,2)`)
			},
			largoEsperado: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("ti_pendientes_documento")
			defer db.Close()

			if test.setup != nil {
				test.setup(db)
			}

			repo := NewTablaIntermediaRepository(db)

			lista, err := repo.ListarDocumentosPorPendiente(1)

			if err != nil {
				t.Errorf("Error inesperado: %v", err)
			}

			if len(lista) != test.largoEsperado {
				t.Errorf("Esperado: %d\nObtenido: %d", test.largoEsperado, len(lista))
			}
		})
	}
}

func TestVincularPendienteCodigoSAP(t *testing.T) {

	tests := []struct {
		name          string
		pendienteID   int
		codigoSAP     string
		setup         func(*sql.DB)
		errorEsperado error
	}{
		{
			name:        "Todo OK",
			pendienteID: 1,
			codigoSAP:   "A",
			setup: func(db *sql.DB) {
				_, _ = db.Exec("PRAGMA foreign_keys = OFF")
			},
			errorEsperado: nil,
		},
		{
			name:          "FK inexistente",
			pendienteID:   999,
			codigoSAP:     "BBB",
			errorEsperado: appErrors.FkNoEncontrado,
		},
		{
			name:        "Duplicado",
			pendienteID: 1,
			codigoSAP:   "A",
			setup: func(db *sql.DB) {
				_, _ = db.Exec("PRAGMA foreign_keys = OFF")
				_, _ = db.Exec(`INSERT INTO ti_pendientes_codigo_sap VALUES (1,'A')`)
			},
			errorEsperado: appErrors.RelacionDuplicada,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("ti_pendientes_codigo_sap")
			defer db.Close()

			if test.setup != nil {
				test.setup(db)
			}

			repo := NewTablaIntermediaRepository(db)

			err := repo.VincularPendienteCodigoSAP(test.pendienteID, test.codigoSAP)

			if !errors.Is(err, test.errorEsperado) {
				t.Errorf("Esperado: %v\nObtenido: %v", test.errorEsperado, err)
			}
		})
	}
}

func TestListarCodigosSAPPorPendiente(t *testing.T) {

	tests := []struct {
		name          string
		setup         func(*sql.DB)
		largoEsperado int
	}{
		{
			name:          "Sin resultados",
			largoEsperado: 0,
		},
		{
			name: "Con resultados",
			setup: func(db *sql.DB) {
				_, _ = db.Exec("PRAGMA foreign_keys = OFF")
				_, _ = db.Exec(`INSERT INTO ti_pendientes_codigo_sap VALUES (1,'A')`)
				_, _ = db.Exec(`INSERT INTO ti_pendientes_codigo_sap VALUES (1,'B')`)
			},
			largoEsperado: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("ti_pendientes_codigo_sap")
			defer db.Close()

			if test.setup != nil {
				test.setup(db)
			}

			repo := NewTablaIntermediaRepository(db)

			lista, err := repo.ListarCodigosSAPPorPendiente(1)

			if err != nil {
				t.Errorf("Error inesperado: %v", err)
			}

			if len(lista) != test.largoEsperado {
				t.Errorf("Esperado: %d\nObtenido: %d", test.largoEsperado, len(lista))
			}
		})
	}
}

func TestVincularPendienteCodigoID(t *testing.T) {

	tests := []struct {
		name          string
		pendienteID   int
		codigoID      string
		setup         func(*sql.DB)
		errorEsperado error
	}{
		{
			name:        "Todo OK",
			pendienteID: 1,
			codigoID:    "A",
			setup: func(db *sql.DB) {
				_, _ = db.Exec("PRAGMA foreign_keys = OFF")
			},
			errorEsperado: nil,
		},
		{
			name:          "FK inexistente",
			pendienteID:   999,
			codigoID:      "BBB",
			errorEsperado: appErrors.FkNoEncontrado,
		},
		{
			name:        "Duplicado",
			pendienteID: 1,
			codigoID:    "A",
			setup: func(db *sql.DB) {
				_, _ = db.Exec("PRAGMA foreign_keys = OFF")
				_, _ = db.Exec(`INSERT INTO ti_pendientes_codigo_id VALUES (1,'A')`)
			},
			errorEsperado: appErrors.RelacionDuplicada,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("ti_pendientes_codigo_id")
			defer db.Close()

			if test.setup != nil {
				test.setup(db)
			}

			repo := NewTablaIntermediaRepository(db)

			err := repo.VincularPendienteCodigoID(test.pendienteID, test.codigoID)

			if !errors.Is(err, test.errorEsperado) {
				t.Errorf("Esperado: %v\nObtenido: %v", test.errorEsperado, err)
			}
		})
	}
}

func TestListarCodigosIDPorPendiente(t *testing.T) {

	tests := []struct {
		name          string
		setup         func(*sql.DB)
		largoEsperado int
	}{
		{
			name:          "Sin resultados",
			largoEsperado: 0,
		},
		{
			name: "Con resultados",
			setup: func(db *sql.DB) {
				_, _ = db.Exec("PRAGMA foreign_keys = OFF")
				_, _ = db.Exec(`INSERT INTO ti_pendientes_codigo_id VALUES (1,'A')`)
				_, _ = db.Exec(`INSERT INTO ti_pendientes_codigo_id VALUES (1,'B')`)
			},
			largoEsperado: 2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			db := generarGestorDbLimpio("ti_pendientes_codigo_id")
			defer db.Close()

			if test.setup != nil {
				test.setup(db)
			}

			repo := NewTablaIntermediaRepository(db)

			lista, err := repo.ListarCodigosIDPorPendiente(1)

			if err != nil {
				t.Errorf("Error inesperado: %v", err)
			}

			if len(lista) != test.largoEsperado {
				t.Errorf("Esperado: %d\nObtenido: %d", test.largoEsperado, len(lista))
			}
		})
	}
}
