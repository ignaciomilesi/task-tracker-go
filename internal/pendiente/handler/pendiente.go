package handler

import (
	"context"
	"strconv"
	"strings"
	"time"

	"task-tracker-go/internal/pendiente/dto"
	"task-tracker-go/internal/pendiente/model"

	"github.com/gin-gonic/gin"
)

type servicePendienteInterface interface {
	Crear(ctx context.Context, nuevaSolicitudPendiente *dto.PendienteRequest) error
	ObtenerDetalle(ctx context.Context, id int) (*dto.PendienteDetalleCompleto, error)
	Actualizar(ctx context.Context, solicitud *dto.ActualizacionPendienteRequest) error

	BuscarPorTituloDescripcion(ctx context.Context, texto string, finalizado bool) ([]dto.PendienteDetalleNombre, error)
	Listar(ctx context.Context, asignadoID *int, finalizado bool, limit, offset int) ([]dto.PendienteDetalleNombre, error)

	VincularCodigoID(ctx context.Context, pendienteID int, codigoID string) error
	VincularCodigoSAP(ctx context.Context, pendienteID int, codigoSAP string) error
	VincularDocumento(ctx context.Context, pendienteID int, documentoID string) error
}

type HandlerPendiente struct {
	service servicePendienteInterface
}

func NewHandlerPendiente(service servicePendienteInterface) *HandlerPendiente {
	return &HandlerPendiente{service: service}
}

func (h *HandlerPendiente) Registrar(c *gin.Context) {

	var req struct {
		Titulo        string  `json:"titulo" binding:"required"`
		Descripcion   string  `json:"descripcion" binding:"required"`
		SolicitanteID int     `json:"solicitante_id" binding:"required"`
		FechaPedido   string  `json:"fecha_pedido" binding:"required"`
		MailPath      *string `json:"mail_path"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": err.Error(),
		})
		return
	}

	fechaParseada, err := time.Parse("02/01/2006", req.FechaPedido)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   "fecha inválida, formato esperado DD/MM/YYYY",
			"detalle": err.Error(),
		})
		return
	}

	nuevoPendiente := dto.PendienteRequest{
		Titulo:        req.Titulo,
		Descripcion:   req.Descripcion,
		SolicitanteID: req.SolicitanteID,
		FechaPedido:   fechaParseada,
		MailPath:      req.MailPath,
	}

	if err := h.service.Crear(c.Request.Context(), &nuevoPendiente); err != nil {
		c.JSON(model.StatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(204)

}

func (h *HandlerPendiente) Actualizar(c *gin.Context) {

	var req struct {
		ID            int    `json:"id" binding:"required"`
		Titulo        string `json:"titulo" binding:"required"`
		Descripcion   string `json:"descripcion" binding:"required"`
		SolicitanteID int    `json:"solicitante_id" binding:"required"`
		FechaPedido   string `json:"fecha_pedido" binding:"required"`

		AsignadoID    *int       `json:"asignado_id"`
		FechaAsignado *time.Time `json:"fecha_asignado"`

		Finalizado  bool       `json:"finalizado"`
		Cierre      *string    `json:"cierre"`
		FechaCierre *time.Time `json:"fecha_cierre"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": err.Error(),
		})
		return
	}

	fechaPedidoParseada, err := time.Parse("02/01/2006", req.FechaPedido)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   "fecha inválida, formato esperado DD/MM/YYYY",
			"detalle": err.Error(),
		})
		return
	}
	var fechaAsignadoParseada *time.Time
	if req.FechaAsignado != nil {
		fecha, err := time.Parse("02/01/2006", req.FechaAsignado.Format("02/01/2006"))
		if err != nil {
			c.JSON(400, gin.H{
				"error":   "fecha de asignado inválida, formato esperado DD/MM/YYYY",
				"detalle": err.Error(),
			})
			return
		}
		fechaAsignadoParseada = &fecha
	}

	var fechaCierreParseada *time.Time
	if req.FechaCierre != nil {
		fecha, err := time.Parse("02/01/2006", req.FechaCierre.Format("02/01/2006"))
		if err != nil {
			c.JSON(400, gin.H{
				"error":   "fecha de cierre inválida, formato esperado DD/MM/YYYY",
				"detalle": err.Error(),
			})
			return
		}
		fechaCierreParseada = &fecha
	}

	actualizacion := dto.ActualizacionPendienteRequest{
		ID:            req.ID,
		Titulo:        req.Titulo,
		Descripcion:   req.Descripcion,
		SolicitanteID: req.SolicitanteID,
		FechaPedido:   fechaPedidoParseada,
		AsignadoID:    req.AsignadoID,
		FechaAsignado: fechaAsignadoParseada,
		Finalizado:    req.Finalizado,
		Cierre:        req.Cierre,
		FechaCierre:   fechaCierreParseada,
	}

	if err := h.service.Actualizar(c.Request.Context(), &actualizacion); err != nil {
		c.JSON(model.StatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(204)
}

func (h *HandlerPendiente) ObtenerDetalle(c *gin.Context) {

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(400, gin.H{
			"error":   "ID de pendiente inválido",
			"detalle": err,
		})
		return
	}

	detalle, err := h.service.ObtenerDetalle(c.Request.Context(), id)
	if err != nil {
		c.JSON(model.StatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"detalle": detalle,
	})
}

func (h *HandlerPendiente) Buscar(c *gin.Context) {

	// obtengo el campo de la query
	query := c.Query("parametroBusqueda")
	if strings.TrimSpace(query) == "" {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": "El campo 'parametroBusqueda' es requerido",
		})
		return
	}

	// obtengo el campo completada de la query
	query = c.Query("completadas")
	completadas, err := strconv.ParseBool(query)
	if err != nil {
		completadas = false // valor por defecto si viene mal
	}

	lista, err := h.service.BuscarPorTituloDescripcion(c.Request.Context(), query, completadas)
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

func (h *HandlerPendiente) Listar(c *gin.Context) {

	var asignadoID *int
	// obtengo el campo asignado de la query
	query := c.Query("asignadoID")
	valor, err := strconv.Atoi(query)
	if err == nil {
		asignadoID = &valor
	}
	// obtengo el campo completada de la query
	query = c.Query("completadas")
	completadas, err := strconv.ParseBool(query)
	if err != nil {
		completadas = false // valor por defecto si viene mal
	}
	query = c.Query("limit")
	limit, err := strconv.Atoi(query)
	if err != nil {
		limit = 100 // valor por defecto si viene mal
	}
	query = c.Query("offset")
	offset, err := strconv.Atoi(query)
	if err != nil {
		offset = 0 // valor por defecto si viene mal
	}

	lista, err := h.service.Listar(c.Request.Context(), asignadoID, completadas, limit, offset)
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

/*----------------------------*/

func (h *HandlerPendiente) VincularCodigoID(c *gin.Context) {

	query := c.Query("pendiente_id")
	id, err := strconv.Atoi(query)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   "ID de pendiente inválido",
			"detalle": err,
		})
		return
	}

	codigoID := c.Query("codigo_id")
	if codigoID == "" {
		c.JSON(400, gin.H{
			"error": "codigo_id es requerido",
		})
		return
	}

	if err := h.service.VincularCodigoID(c.Request.Context(), id, codigoID); err != nil {
		c.JSON(model.StatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(204)
}

func (h *HandlerPendiente) VincularCodigoSAP(c *gin.Context) {

	query := c.Query("pendiente_id")
	id, err := strconv.Atoi(query)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   "ID de pendiente inválido",
			"detalle": err,
		})
		return
	}

	codigoSAP := c.Query("codigo_sap")
	if codigoSAP == "" {
		c.JSON(400, gin.H{
			"error": "codigo_sap es requerido",
		})
		return
	}

	if err := h.service.VincularCodigoSAP(c.Request.Context(), id, codigoSAP); err != nil {
		c.JSON(model.StatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(204)
}

func (h *HandlerPendiente) VincularDocumento(c *gin.Context) {

	query := c.Query("pendiente_id")
	id, err := strconv.Atoi(query)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   "ID de pendiente inválido",
			"detalle": err,
		})
		return
	}

	documentoID := c.Query("documento_id")
	if documentoID == "" {
		c.JSON(400, gin.H{
			"error": "documento_id es requerido",
		})
		return
	}

	if err := h.service.VincularDocumento(c.Request.Context(), id, documentoID); err != nil {
		c.JSON(model.StatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(204)
}
