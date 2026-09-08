package codigoID

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type serviceInterface interface {
	Cargar(context.Context, *nuevoCodigoIDRequest) error
	ObtenerDetalle(context.Context, string) (*CodigoID, error)
	FiltrarPorEstado(context.Context, string) ([]CodigoID, error)
	ActualizarEstado(ctx context.Context, codigo string, nuevoEstado string, fecha time.Time) error
	Listar(context.Context, int, int) ([]CodigoID, error)
}

type handler struct {
	service serviceInterface
}

func newHandler(service serviceInterface) *handler {
	return &handler{service: service}
}

func (h *handler) Registrar(c *gin.Context) {
	solicitudRegistro := struct {
		Codigo      string  `json:"codigo" binding:"required"`
		Descripcion *string `json:"descripcion"`
		FechaPedido string  `json:"fecha_pedido" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&solicitudRegistro); err != nil {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": err,
		})
		return
	}

	fechaPedidoParseada, err := time.Parse("02/01/2006", solicitudRegistro.FechaPedido)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   "fecha pedido inválida, formato esperado DD/MM/YYYY",
			"detalle": err,
		})
		return
	}

	nuevaSolicitud := nuevoCodigoIDRequest{
		Codigo:      solicitudRegistro.Codigo,
		Descripcion: solicitudRegistro.Descripcion,
		FechaPedido: fechaPedidoParseada,
	}

	err = h.service.Cargar(c.Request.Context(), &nuevaSolicitud)
	if err != nil {
		c.JSON(StatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"mensaje": "Código ID registrado exitosamente",
	})
}

func (h *handler) Listar(c *gin.Context) {
	// obtengo el campo completada de la query
	query := c.Query("limit")
	limit, err := strconv.Atoi(query)
	if err != nil {
		limit = 100 // valor por defecto si viene mal
	}
	query = c.Query("offset")
	offset, err := strconv.Atoi(query)
	if err != nil {
		offset = 0 // valor por defecto si viene mal
	}

	lista, err := h.service.Listar(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(StatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"lista": lista,
	})
}

func (h *handler) ObtenerDetalle(c *gin.Context) {

	// obtengo el campo de la query
	query := c.Query("codigoID")
	if strings.TrimSpace(query) == "" {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": "El campo 'codigoID' es requerido",
		})
		return
	}

	detalle, err := h.service.ObtenerDetalle(c.Request.Context(), query)
	if err != nil {
		c.JSON(StatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"detalle": detalle,
	})
}

func (h *handler) FiltrarPorEstado(c *gin.Context) {

	// obtengo el campo de la query
	query := c.Query("estado")
	if strings.TrimSpace(query) == "" {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": "El campo 'estado' es requerido",
		})
		return
	}

	lista, err := h.service.FiltrarPorEstado(c.Request.Context(), query)
	if err != nil {
		c.JSON(StatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"lista": lista,
	})
}

func (h *handler) Actualizar(c *gin.Context) {

	actualizacion := struct {
		Codigo      string `json:"codigo" binding:"required"`
		NuevoEstado string `json:"nuevo_estado" binding:"required"`
		Fecha       string `json:"fecha" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&actualizacion); err != nil {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": err,
		})
		return
	}

	fechaParseada, err := time.Parse("02/01/2006", actualizacion.Fecha)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   "fecha inválida, formato esperado DD/MM/YYYY",
			"detalle": err,
		})
		return
	}

	err = h.service.ActualizarEstado(c.Request.Context(), actualizacion.Codigo, actualizacion.NuevoEstado, fechaParseada)
	if err != nil {
		c.JSON(StatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"mensaje": "Estado actualizado exitosamente",
	})
}
