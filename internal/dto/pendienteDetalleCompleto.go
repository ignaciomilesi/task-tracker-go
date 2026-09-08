package dto

import (
	"task-tracker-go/internal/codigoID"
	"task-tracker-go/internal/codigoSAP"
	"task-tracker-go/internal/documento"
	pendiente "task-tracker-go/internal/pendiente/model"
)

type PendientesDetalleCompleto struct {
	pendiente.Pendiente

	Avances    []pendiente.Avance
	Adjuntos   []pendiente.Adjunto
	CodigosID  []codigoID.CodigoID
	CodigosSAP []codigoSAP.CodigoSAP
	Documentos []documento.Documento
}
