package repo

import "context"

// Interfaz para la capa de almacenamiento
type StorageInterface interface {
	Cargar(ctx context.Context, strct any, query string) (int, error)
	ModificarRegistro(ctx context.Context, query string, arg ...interface{}) error
	ObtenerDetalle(ctx context.Context, strct any, query string, arg ...interface{}) error
	ObtenerLista(ctx context.Context, strct any, query string, arg ...interface{}) error
}

// repo y funciones
type Repo struct {
	db StorageInterface
}

func NewRepo(db StorageInterface) *Repo {
	return &Repo{db: db}
}
