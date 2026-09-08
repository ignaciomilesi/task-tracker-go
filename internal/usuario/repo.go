package usuario

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Interfaz para la capa de almacenamiento
type storageInterface interface {
	Cargar(ctx context.Context, strct any, query string) error
	ModificarRegistro(ctx context.Context, query string, arg ...interface{}) error
	ObtenerDetalle(ctx context.Context, strct any, query string, arg ...interface{}) error
	ObtenerLista(ctx context.Context, strct any, query string, arg ...interface{}) error
}

// Estructura que representa la tabla de usuario en la base de datos
type usuarioDB struct {
	ID          int    `db:"id"`
	Nombre      string `db:"nombre"`
	Colaborador bool   `db:"colaborador"`
}

// repo y funciones
type repo struct {
	db storageInterface
}

func newRepo(db storageInterface) *repo {
	return &repo{db: db}
}

func (r *repo) Crear(ctx context.Context, nuevoUsuario *usuario) error {

	if nuevoUsuario == nil {
		return errModelDeCargaNil
	}

	var nuevoUsuarioDB = usuarioDB{
		Nombre:      nuevoUsuario.nombre,
		Colaborador: nuevoUsuario.colaborador,
	}

	query := `INSERT INTO usuario (nombre, colaborador) 
		VALUES (:nombre, :colaborador);`

	if err := r.db.Cargar(ctx, &nuevoUsuarioDB, query); err != nil {
		return err
	}
	return nil
}

func (r *repo) ObtenerIDPorNombre(ctx context.Context, nombre string) (int, error) {

	if nombre == "" {
		return 0, errParametrosNoValidos
	}

	var id int

	query := "SELECT id FROM usuario WHERE nombre = ?"

	err := r.db.ObtenerDetalle(ctx, &id, query, nombre)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errUsuarioNoEncontrado
		}
		return 0, fmt.Errorf("Error inesperado, detalle: %w", err)
	}

	return id, nil
}

func (r *repo) Buscar(ctx context.Context, parametro string) ([]usuario, error) {
	if parametro == "" {
		return nil, errParametrosNoValidos
	}

	var lista []usuarioDB

	query := `SELECT id, nombre FROM usuario WHERE nombre LIKE ?`

	err := r.db.ObtenerLista(ctx, &lista, query, "%"+parametro+"%")
	if err != nil {
		return nil, err
	}

	var usuarios []usuario

	for _, registro := range lista {
		usuarios = append(usuarios, usuario{
			id:          registro.ID,
			nombre:      registro.Nombre,
			colaborador: registro.Colaborador,
		})
	}

	return usuarios, nil
}

func (r *repo) Listar(ctx context.Context, limit, offset int) ([]usuario, error) {

	var lista []usuarioDB

	query := `SELECT id, nombre FROM usuario ORDER BY id
		LIMIT ? OFFSET ?`

	err := r.db.ObtenerLista(ctx, &lista, query, limit, offset)
	if err != nil {
		return nil, err
	}

	var usuarios []usuario

	for _, registro := range lista {
		usuarios = append(usuarios, usuario{
			id:          registro.ID,
			nombre:      registro.Nombre,
			colaborador: registro.Colaborador,
		})
	}

	return usuarios, nil
}
