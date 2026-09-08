package usuario

import (
	"context"
	"fmt"
	"strings"
)

type repoInterface interface {
	Crear(context.Context, *usuario) error
	ObtenerIDPorNombre(context.Context, string) (int, error)
	Buscar(context.Context, string) ([]usuario, error)
	Listar(context.Context, int, int) ([]usuario, error)
}

type service struct {
	repo repoInterface
}

func newService(repo repoInterface) *service {
	return &service{repo: repo}
}

// funciona como DTO
type nuevoUsuarioRequest struct {
	Nombre      string
	Colaborador bool
}

func (s *service) Crear(ctx context.Context, nuevoSolicitudUsuario *nuevoUsuarioRequest) error {
	if nuevoSolicitudUsuario == nil {
		return fmt.Errorf("service: %w", errModelDeCargaNil)
	}

	if strings.TrimSpace(nuevoSolicitudUsuario.Nombre) == "" {
		return fmt.Errorf("service: %w, parametro: nombre", errParametrosNoValidos)
	}

	usuario := NewUsuario(nuevoSolicitudUsuario.Nombre, nuevoSolicitudUsuario.Colaborador)

	return s.repo.Crear(ctx, usuario)
}

func (s *service) ObtenerIDPorNombre(ctx context.Context, nombre string) (int, error) {
	if strings.TrimSpace(nombre) == "" {
		return 0, fmt.Errorf("service: %w, nombre vacío", errParametrosNoValidos)
	}
	return s.repo.ObtenerIDPorNombre(ctx, nombre)
}

func (s *service) Buscar(ctx context.Context, parametro string) ([]usuario, error) {
	if strings.TrimSpace(parametro) == "" {
		return nil, fmt.Errorf("service: %w, parámetro vacío", errParametrosNoValidos)
	}
	return s.repo.Buscar(ctx, parametro)
}

func (s *service) Listar(ctx context.Context, limit, offset int) ([]usuario, error) {
	if limit <= 0 || offset < 0 {
		return nil, fmt.Errorf("service: %w, parámetros de búsqueda inválidos", errParametrosNoValidos)
	}
	return s.repo.Listar(ctx, limit, offset)
}
