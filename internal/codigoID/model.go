package codigoID

import (
	"fmt"
	"strconv"
	"time"
)

type CodigoID struct {
	codigo             string
	descripcion        *string
	estado             string
	fechaPedido        time.Time
	fechaActualizacion *time.Time
}

func NewCodigoID(
	codigo string,
	descripcion *string,
	estado string,
	fechaPedido time.Time,
	fechaActualizacion *time.Time,
) (*CodigoID, error) {

	err := validarCodigo(codigo)
	if err != nil {
		return nil, err
	}

	if descripcion != nil && *descripcion == "" {
		descripcion = nil
	}

	if fechaPedido.After(time.Now()) {
		return nil, fmt.Errorf("%w: la fecha de pedido no puede ser futura", errParametrosNoValidos)
	}

	if fechaActualizacion != nil && fechaActualizacion.After(time.Now()) {
		fechaActualizacion = nil
	}

	nuevoCodigo := CodigoID{
		codigo:             codigo,
		descripcion:        descripcion,
		estado:             estado,
		fechaPedido:        fechaPedido,
		fechaActualizacion: fechaActualizacion,
	}

	return &nuevoCodigo, nil
}

func validarCodigo(codigo string) error {
	if len(codigo) > 8 {
		return fmt.Errorf("%w: el códigoID no posee el largo esperado", errParametrosNoValidos)
	}

	if _, err := strconv.Atoi(codigo); err != nil {
		return fmt.Errorf("%w: no es un código válido", errParametrosNoValidos)
	}

	return nil
}
