package handler

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"

	"task-tracker-go/internal/pendiente/dto"
	"task-tracker-go/internal/pendiente/model"
)

type serviceAdjuntoInterface interface {
	CargarAdjunto(ctx context.Context, nuevoAdjunto *dto.AdjuntoRequest) error
	EliminarAdjunto(ctx context.Context, id int) error
	ObtenerListaAdjuntos(ctx context.Context, pendienteID int) ([]dto.AdjuntoDetalleRespuesta, error)
}

type HandlerAdjunto struct {
	service serviceAdjuntoInterface
}

func NewHandlerAdjunto(service serviceAdjuntoInterface) *HandlerAdjunto {
	return &HandlerAdjunto{service: service}
}

func (h *HandlerAdjunto) RegistrarAdjunto(c *gin.Context) {

	solicitudRegistro := struct {
		PendienteID int    `json:"pendiente_id" binding:"required"`
		Descripcion string `json:"descripcion" binding:"required"`
		ArchivoPath string `json:"archivo_path" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&solicitudRegistro); err != nil {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": err,
		})
		return
	}

	nuevoAdjunto := dto.AdjuntoRequest{
		PendienteID: solicitudRegistro.PendienteID,
		Descripcion: solicitudRegistro.Descripcion,
		ArchivoPath: solicitudRegistro.ArchivoPath,
	}

	err := h.service.CargarAdjunto(c.Request.Context(), &nuevoAdjunto)
	if err != nil {
		c.JSON(model.StatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(200)
}

func (h *HandlerAdjunto) ListarAdjunto(c *gin.Context) {

	idParam := c.Param("id")
	idPendiente, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   "ID de pendiente inválido",
			"detalle": err.Error(),
		})
		return
	}

	lista, err := h.service.ObtenerListaAdjuntos(c.Request.Context(), idPendiente)
	if err != nil {
		c.JSON(model.StatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"lista": lista,
	})
}
