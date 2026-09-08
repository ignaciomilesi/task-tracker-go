package codigoSAP

import (
	"fmt"
	"strconv"
)

type CodigoSAP struct {
	codigo      string
	descripcion *string
}

func NewCodigoID(
	codigo string,
	descripcion *string,
) (*CodigoSAP, error) {

	err := validarCodigo(codigo)
	if err != nil {
		return nil, err
	}

	if descripcion != nil && *descripcion == "" {
		descripcion = nil
	}

	nuevoCodigo := CodigoSAP{
		codigo:      codigo,
		descripcion: descripcion,
	}

	return &nuevoCodigo, nil
}

func validarCodigo(codigo string) error {
	if len(codigo) != 10 {
		return fmt.Errorf("%w: el códigoID no posee el largo esperado", errParametrosNoValidos)
	}

	if _, err := strconv.Atoi(codigo); err != nil {
		return fmt.Errorf("%w: no es un código válido", errParametrosNoValidos)
	}

	return nil
}
