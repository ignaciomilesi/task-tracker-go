package codigoID

import (
	"errors"
	"net/http"
)

var (
	errModelDeCargaNil     = errors.New("modelo de carga es nil")
	errParametrosNoValidos = errors.New("parámetros no válidos")

	errCodigoIDDuplicado    = errors.New("código ID duplicado")
	errCodigoIDNoEncontrado = errors.New("código ID no encontrado")
)

func StatusCode(err error) int {
	switch {
	case errors.Is(err, errModelDeCargaNil):
		return http.StatusInternalServerError

	case errors.Is(err, errParametrosNoValidos):
		return http.StatusBadRequest

	case errors.Is(err, errCodigoIDDuplicado),
		errors.Is(err, errCodigoIDNoEncontrado):
		return http.StatusNotFound

	default:
		return http.StatusInternalServerError
	}
}
