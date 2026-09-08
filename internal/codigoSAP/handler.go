package codigoSAP

import (
	"context"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type serviceInterface interface {
	Cargar(context.Context, *nuevoCodigoSAPRequest) error
	ObtenerDetalle(context.Context, string) (*CodigoSAP, error)
	BuscarPorDescripcion(context.Context, string) ([]CodigoSAP, error)
	ModificarDescripcion(context.Context, string, string) error
	Listar(context.Context, int, int) ([]CodigoSAP, error)
	PendientesRelacionados(context.Context, string) ([]pendiente, error)
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
	}{}

	if err := c.ShouldBindJSON(&solicitudRegistro); err != nil {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": err,
		})
		return
	}

	nuevoCodigoSAP := nuevoCodigoSAPRequest{
		Codigo:      solicitudRegistro.Codigo,
		Descripcion: solicitudRegistro.Descripcion,
	}

	err := h.service.Cargar(c.Request.Context(), &nuevoCodigoSAP)
	if err != nil {
		c.JSON(StatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"mensaje": "Código SAP registrado exitosamente",
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
	query := c.Query("codigoSAP")
	if strings.TrimSpace(query) == "" {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": "El campo 'codigoSAP' es requerido",
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

func (h *handler) BuscarPorDescripcion(c *gin.Context) {

	// obtengo el campo de la query
	query := c.Query("parametroBusqueda")
	if strings.TrimSpace(query) == "" {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": "El campo 'parametroBusqueda' es requerido",
		})
		return
	}

	lista, err := h.service.BuscarPorDescripcion(c.Request.Context(), query)
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
		Codigo           string `json:"codigo" binding:"required"`
		NuevaDescripcion string `json:"nueva_descripcion" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&actualizacion); err != nil {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": err,
		})
		return
	}

	err := h.service.ModificarDescripcion(c.Request.Context(), actualizacion.Codigo, actualizacion.NuevaDescripcion)
	if err != nil {
		c.JSON(StatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"mensaje": "Descripcion actualizado exitosamente",
	})
}

func (h *handler) PendientesRelacionados(c *gin.Context) {

	documentoID := c.Query("documentoID")
	if strings.TrimSpace(documentoID) == "" {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": "El campo 'documentoID' es requerido",
		})
		return
	}

	lista, err := h.service.PendientesRelacionados(c.Request.Context(), documentoID)
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
