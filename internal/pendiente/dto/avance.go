package dto

import "time"

type AvanceRequest struct {
	PendienteID int
	Descripcion string
	Fecha       time.Time
	MailPath    *string
}

type AvanceDetalleRespuesta struct {
	Descripcion string
	Fecha       time.Time
	MailPath    *string
}
