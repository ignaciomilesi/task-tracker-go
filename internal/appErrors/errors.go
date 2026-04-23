package appErrors

import "errors"

var (
	ParametroDeCargaVacio    = errors.New("El parámetro de carga no puede ser vacío")
	ParametroDeBusquedaVacio = errors.New("El parámetro de búsqueda no puede ser vacío")

	SolicitanteNoEncontrado = errors.New("El Solicitante no fue encontrado")
	SolicitanteDuplicado    = errors.New("El Solicitante ya existe en la base de datos")

	ColaboradorNoEncontrado = errors.New("El Colaborador no fue encontrado")
	ColaboradorDuplicado    = errors.New("El Colaborador ya existe en la base de datos")

	CodigoSAPNoEncontrado = errors.New("El Código SAP no fue encontrado")
	CodigoSAPDuplicado    = errors.New("El Código SAP ya existe en la base de datos")
	CodigoSAPVacio        = errors.New("El código SAP que se desea dar de alta esta en blanco")

	CodigoIDNoEncontrado = errors.New("El Código ID no fue encontrado")
	CodigoIDDuplicado    = errors.New("El Código ID ya existe en la base de datos")
	CodigoIDVacio        = errors.New("El código ID que se desea dar de alta esta en blanco")

	DocumentoDuplicado    = errors.New("El documento ya existe en la base de datos")
	DocumentoNoEncontrado = errors.New("El documento no fue encontrado")

	PendienteNoEncontrado = errors.New("El Pendiente no fue encontrado")

	AvanceNoEncontrado = errors.New("El avance no fue encontrado")

	AdjuntoNoEncontrado = errors.New("El adjunto no fue encontrado")

	FkNoEncontrado    = errors.New("Uno de los IDs no existe (foreign key)")
	RelacionDuplicada = errors.New("La relación entre ID ya existe en la base de datos")
)
