package model

import (
	"time"
)

type Pendiente struct {
	ID            int
	Titulo        string
	Descripcion   string
	SolicitanteID int
	FechaPedido   time.Time

	AsignadoID    *int //colaborador
	FechaAsignado *time.Time

	Finalizado  bool
	Cierre      *string
	FechaCierre *time.Time
}

type Avance struct {
	ID          int
	PendienteID int
	Descripcion string
	Fecha       time.Time
	MailPath    *string
}

type Adjunto struct {
	ID          int
	PendienteID int
	Descripcion string
	ArchivoPath string
}
