package usuario

import (
	"context"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type serviceInterface interface {
	Crear(context.Context, *nuevoUsuarioRequest) error
	ObtenerIDPorNombre(context.Context, string) (int, error)
	Buscar(context.Context, string) ([]usuario, error)
	Listar(context.Context, int, int) ([]usuario, error)
}

type handler struct {
	service serviceInterface
}

func newHandler(service serviceInterface) *handler {
	return &handler{service: service}
}

func (h *handler) Registrar(c *gin.Context) {

	solicitudRegistro := struct {
		Nombre      string `json:"nombre" binding:"required"`
		Colaborador bool   `json:"colaborador" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&solicitudRegistro); err != nil {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": err,
		})
		return
	}

	nuevoUsuario := &nuevoUsuarioRequest{
		Nombre:      solicitudRegistro.Nombre,
		Colaborador: solicitudRegistro.Colaborador,
	}

	err := h.service.Crear(c.Request.Context(), nuevoUsuario)
	if err != nil {
		c.JSON(StatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(204)
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

func (h *handler) ConsultarIDporNombre(c *gin.Context) {

	var nombreParaConsultar string

	if err := c.ShouldBindJSON(&nombreParaConsultar); err != nil {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": err,
		})
		return
	}

	id, err := h.service.ObtenerIDPorNombre(c.Request.Context(), nombreParaConsultar)
	if err != nil {
		c.JSON(StatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"id": id,
	})
}

func (h *handler) BuscarPorNombre(c *gin.Context) {
	// obtengo el campo completada de la query
	query := c.Query("nombre")
	if strings.TrimSpace(query) == "" {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": "El campo 'nombre' es requerido",
		})
		return
	}

	lista, err := h.service.Buscar(c.Request.Context(), query)
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
