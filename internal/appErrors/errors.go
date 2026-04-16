package appErrors

import "errors"

var (
	SolicitanteNoEncontrado = errors.New("Solicitante no fue encontrado")
	SolicitanteDuplicado    = errors.New("El Solicitante ya existe en la base de datos")

	ColaboradorNoEncontrado = errors.New("El Colaborador no fue encontrado")
	ColaboradorDuplicado    = errors.New("El Colaborador ya existe en la base de datos")

	CodigoSAPNoEncontrado = errors.New("El Código SAP no fue encontrado")
	CodigoSAPDuplicado    = errors.New("El Código SAP ya existe en la base de datos")
	CodigoSAPVacio        = errors.New("El código SAP que se desea dar de alta esta en blanco")
)
