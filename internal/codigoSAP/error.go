package codigoSAP

import (
	"errors"
	"net/http"
)

var (
	errModelDeCargaNil     = errors.New("modelo de carga es nil") //ok
	errParametrosNoValidos = errors.New("parámetros no válidos")  //ok

	errCodigoSAPDuplicado    = errors.New("código SAP duplicado")     //ok
	errCodigoSAPNoEncontrado = errors.New("código SAP no encontrado") //ok
)

func StatusCode(err error) int {
	switch {
	case errors.Is(err, errModelDeCargaNil):
		return http.StatusInternalServerError

	case errors.Is(err, errParametrosNoValidos):
		return http.StatusBadRequest

	case errors.Is(err, errCodigoSAPDuplicado),
		errors.Is(err, errCodigoSAPNoEncontrado):
		return http.StatusNotFound

	default:
		return http.StatusInternalServerError
	}
}
