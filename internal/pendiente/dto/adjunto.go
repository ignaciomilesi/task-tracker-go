package dto

type AdjuntoRequest struct {
	PendienteID int
	Descripcion string
	ArchivoPath string
}

type AdjuntoDetalleRespuesta struct {
	Descripcion string
	ArchivoPath string
}
