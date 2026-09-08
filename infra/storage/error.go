package storage

import "errors"

var (
	ErrRegistroNoEncontrado     = errors.New("El registro no existe en la base de datos")
	ErrRegistroYaExistente      = errors.New("El registro ya existe en la base de datos")
	ErrNingunRegistroEncontrado = errors.New("No se encontró algún registro para modificar")
)
