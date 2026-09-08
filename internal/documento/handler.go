package documento

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
)

type serviceInterface interface {
	Cargar(context.Context, *nuevoDocumentoRequest) error
	ObtenerDetalle(context.Context, string) (*Documento, error)
	FiltrarPorTipo(context.Context, string) ([]Documento, error)
	FiltrarPorTitulo(context.Context, string) ([]Documento, error)
	ActualizarPath(context.Context, string, string, string) error
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
		Codigo        string  `json:"codigo" binding:"required"`
		Emision       string  `json:"emision" binding:"required"`
		Titulo        string  `json:"titulo" binding:"required"`
		Tipo          string  `json:"tipo" binding:"required"`
		UbicacionPath *string `json:"ubicacion_path"`
		BackupPath    *string `json:"backup_path"`
	}{}

	if err := c.ShouldBindJSON(&solicitudRegistro); err != nil {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": err,
		})
		return
	}

	nuevoDocumento := nuevoDocumentoRequest{
		Codigo:        solicitudRegistro.Codigo,
		Emision:       solicitudRegistro.Emision,
		Titulo:        solicitudRegistro.Titulo,
		Tipo:          solicitudRegistro.Tipo,
		UbicacionPath: solicitudRegistro.UbicacionPath,
		BackupPath:    solicitudRegistro.BackupPath,
	}

	err := h.service.Cargar(c.Request.Context(), &nuevoDocumento)
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

func (h *handler) ListarPorTipo(c *gin.Context) {
	// obtengo el campo completada de la query
	tipo := c.Query("tipo")
	if strings.TrimSpace(tipo) == "" {
		c.JSON(422, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": "parametro 'tipo' no puede ser vacío",
		})
		return
	}

	lista, err := h.service.FiltrarPorTipo(c.Request.Context(), tipo)
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
	query := c.Query("codigo")
	if strings.TrimSpace(query) == "" {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": "El campo 'codigo' es requerido",
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

func (h *handler) FiltrarPorTitulo(c *gin.Context) {

	// obtengo el campo de la query
	query := c.Query("titulo")
	if strings.TrimSpace(query) == "" {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": "El campo 'titulo' es requerido",
		})
		return
	}

	lista, err := h.service.FiltrarPorTitulo(c.Request.Context(), query)
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

func (h *handler) ActualizarPath(c *gin.Context) {

	// obtengo el campo de la query
	path := c.Query("path")
	if strings.TrimSpace(path) == "" {
		c.JSON(422, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": "parametro 'path' no puede ser vacío",
		})
		return
	}

	actualizacion := struct {
		Codigo    string `json:"codigo" binding:"required"`
		NuevoPath string `json:"nuevo_path" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&actualizacion); err != nil {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": err,
		})
		return
	}

	err := h.service.ActualizarPath(c.Request.Context(), actualizacion.Codigo, actualizacion.NuevoPath, path)
	if err != nil {
		c.JSON(StatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"mensaje": "Path actualizado exitosamente",
	})
}

func (h *handler) PendientesRelacionados(c *gin.Context) {

	codigoDocumento := c.Query("codigo_documento")
	if strings.TrimSpace(codigoDocumento) == "" {
		c.JSON(400, gin.H{
			"error":   "Campos requeridos no válidos",
			"detalle": "El campo 'codigo_documento' es requerido",
		})
		return
	}

	lista, err := h.service.PendientesRelacionados(c.Request.Context(), codigoDocumento)
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
