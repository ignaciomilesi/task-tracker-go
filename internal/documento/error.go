package documento

import (
	"errors"
	"net/http"
)

var (
	errModelDeCargaNil     = errors.New("modelo de carga es nil") //ok
	errParametrosNoValidos = errors.New("parámetros no válidos")  //ok

	errDocumentoDuplicado    = errors.New("documento ya existe")     //ok
	errDocumentoNoEncontrado = errors.New("documento no encontrado") //ok
)

func StatusCode(err error) int {
	switch {
	case errors.Is(err, errModelDeCargaNil):
		return http.StatusInternalServerError

	case errors.Is(err, errParametrosNoValidos):
		return http.StatusBadRequest

	case errors.Is(err, errDocumentoDuplicado),
		errors.Is(err, errDocumentoNoEncontrado):
		return http.StatusNotFound

	default:
		return http.StatusInternalServerError
	}
}
