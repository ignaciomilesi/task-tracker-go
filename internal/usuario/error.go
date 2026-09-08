package usuario

import (
	"errors"
	"net/http"
)

var (
	errModelDeCargaNil     = errors.New("modelo de carga es nil") //ok
	errParametrosNoValidos = errors.New("parámetros no válidos")  //ok

	errUsuarioDuplicado    = errors.New("usuario ya existe")     //ok
	errUsuarioNoEncontrado = errors.New("usuario no encontrado") //ok
)

func StatusCode(err error) int {
	switch {
	case errors.Is(err, errModelDeCargaNil):
		return http.StatusInternalServerError

	case errors.Is(err, errParametrosNoValidos):
		return http.StatusBadRequest

	case errors.Is(err, errUsuarioNoEncontrado):
		return http.StatusNotFound

	default:
		return http.StatusInternalServerError
	}
}
