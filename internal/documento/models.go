package documento

import (
	"fmt"
	"slices"
)

type Documento struct {
	codigo        string
	emision       string
	titulo        string
	tipo          string
	ubicacionPath *string
	backupPath    *string
}

var tipos = []string{"plano", "et", "dtc", "informe"}

func NewDocumento(
	codigo string,
	emision string,
	titulo string,
	tipo string,
	ubicacionPath *string,
	backupPath *string,
) (*Documento, error) {

	if err := validarTipo(tipo); err != nil {
		return nil, err
	}

	if ubicacionPath != nil && *ubicacionPath == "" {
		ubicacionPath = nil
	}

	if backupPath != nil && *backupPath == "" {
		backupPath = nil
	}

	nuevoDoc := Documento{
		codigo:        codigo,
		emision:       emision,
		titulo:        titulo,
		tipo:          tipo,
		ubicacionPath: ubicacionPath,
		backupPath:    backupPath,
	}

	return &nuevoDoc, nil
}

func validarTipo(tipo string) error {
	if !slices.Contains(tipos[:], tipo) {
		return fmt.Errorf("%w: el tipo no es valido", errParametrosNoValidos)
	}
	return nil
}
