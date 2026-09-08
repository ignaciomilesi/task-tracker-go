package model

import (
	"errors"
	"net/http"
)

var (
	ErrModelDeCargaNil     = errors.New("modelo de carga es nil")
	ErrParametrosNoValidos = errors.New("parámetros no válidos")

	ErrPendienteNoEncontrado = errors.New("pendiente no encontrado")
	ErrAvanceNoEncontrado    = errors.New("avance no encontrado")
	ErrAdjuntoNoEncontrado   = errors.New("adjunto no encontrado")
	ErrUsuarioNoEncontrado   = errors.New("usuario no encontrado")

	ErrFkNoEncontrado    = errors.New("Uno de los IDs no existe (foreign key)")
	ErrRelacionDuplicada = errors.New("La relación entre ID ya existe en la base de datos")
)

func StatusCode(err error) int {
	switch {
	case errors.Is(err, ErrModelDeCargaNil):
		return http.StatusInternalServerError

	case errors.Is(err, ErrParametrosNoValidos):
		return http.StatusBadRequest

	case errors.Is(err, ErrPendienteNoEncontrado),
		errors.Is(err, ErrAvanceNoEncontrado),
		errors.Is(err, ErrAdjuntoNoEncontrado),
		errors.Is(err, ErrUsuarioNoEncontrado),
		errors.Is(err, ErrFkNoEncontrado):
		return http.StatusNotFound

	case errors.Is(err, ErrRelacionDuplicada):
		return http.StatusConflict

	default:
		return http.StatusInternalServerError
	}
}
