package appErrors

import "errors"

var (
	SolicitanteNoEncontrado = errors.New("Solicitante no encontrado")
	ColaboradorNoEncontrado = errors.New("Colaborador no encontrado")
)
