package dto

import (
	"time"
)

type PendienteDetalleNombre struct {
	ID                int
	Titulo            string
	Descripcion       string
	SolicitanteNombre string
	FechaPedido       time.Time

	AsignadoNombre *string
	FechaAsignado  *time.Time

	Finalizado  bool
	Cierre      *string
	FechaCierre *time.Time
}

type PendienteDetalleCompleto struct {
	PendienteDetalleNombre
	Avances  []AvanceDetalleRespuesta
	Adjuntos []AdjuntoDetalleRespuesta
}

type PendienteRequest struct {
	Titulo        string
	Descripcion   string
	SolicitanteID int
	FechaPedido   time.Time
	MailPath      *string
}

type ActualizacionPendienteRequest struct {
	ID            int
	Titulo        string
	Descripcion   string
	SolicitanteID int
	FechaPedido   time.Time

	AsignadoID    *int
	FechaAsignado *time.Time

	Finalizado  bool
	Cierre      *string
	FechaCierre *time.Time
}
