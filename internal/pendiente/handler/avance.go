package handler

import (
	"context"
	"strconv"
	"time"

	"task-tracker-go/internal/pendiente/dto"
	"task-tracker-go/internal/pendiente/model"

	"github.com/gin-gonic/gin"
)

type serviceAvanceInterface interface {
	CargarAvance(ctx context.Context, nuevoAvance *dto.AvanceRequest) error
	EliminarAvance(ctx context.Context, id int) error
	ObtenerListaAvances(ctx context.Context, pendienteID int) ([]dto.AvanceDetalleRespuesta, error)
}

type HandlerAvance struct {
	service serviceAvanceInterface
}

func NewHandlerAvance(service serviceAvanceInterface) *HandlerAvance {
	return &HandlerAvance{service: service}
}

func (h *HandlerAvance) RegistrarAvance(c *gin.Context) {

	solicitudRegistro := struct {
		PendienteID int     `json:"pendiente_id" binding:"required"`
		Descripcion string  `json:"descripcion" binding:"required"`
		Fecha       string  `json:"fecha" binding:"required"`
		MailPath    *string `json:"mail_path"`
	}{}

	if err := c.ShouldBindJSON(&solicitudRegistro); err != nil {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": err.Error(),
		})
		return
	}

	fechaParseada, err := time.Parse("02/01/2006", solicitudRegistro.Fecha)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   "fecha inválida, formato esperado DD/MM/YYYY",
			"detalle": err.Error(),
		})
		return
	}

	nuevoAvance := dto.AvanceRequest{
		PendienteID: solicitudRegistro.PendienteID,
		Descripcion: solicitudRegistro.Descripcion,
		Fecha:       fechaParseada,
		MailPath:    solicitudRegistro.MailPath,
	}

	if err := h.service.CargarAvance(c.Request.Context(), &nuevoAvance); err != nil {
		c.JSON(model.StatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(200)
}

func (h *HandlerAvance) ListarAvance(c *gin.Context) {
	idParam := c.Param("id")
	idPendiente, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   "ID de pendiente inválido",
			"detalle": err.Error(),
		})
		return
	}

	lista, err := h.service.ObtenerListaAvances(c.Request.Context(), idPendiente)
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
