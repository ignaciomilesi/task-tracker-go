package codigoID

import (
	"context"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Interfaz para la capa de almacenamiento
type storageInterface interface {
	Cargar(ctx context.Context, strct any, query string) error
	ModificarRegistro(ctx context.Context, query string, arg ...interface{}) error
	ObtenerDetalle(ctx context.Context, strct any, query string, arg ...interface{}) error
	ObtenerLista(ctx context.Context, strct any, query string, arg ...interface{}) error
}

// Estructura que representa la tabla de código_ID en la base de datos
type codigoIDDB struct {
	Codigo             string     `db:"codigo"`
	Descripcion        *string    `db:"descripcion"`
	Estado             string     `db:"estado"`
	FechaPedido        time.Time  `db:"fecha_pedido"`
	FechaActualizacion *time.Time `db:"fecha_actualizacion"`
}

// repo y funciones
type repo struct {
	db storageInterface
}

func newRepo(db storageInterface) *repo {
	return &repo{db: db}
}

func (r *repo) Cargar(ctx context.Context, nuevoCodigoID *CodigoID) error {

	if nuevoCodigoID == nil {
		return errModelDeCargaNil
	}

	var nuevoCodigo = codigoIDDB{
		Codigo:             nuevoCodigoID.codigo,
		Descripcion:        nuevoCodigoID.descripcion,
		Estado:             nuevoCodigoID.estado,
		FechaPedido:        nuevoCodigoID.fechaPedido,
		FechaActualizacion: nuevoCodigoID.fechaActualizacion,
	}

	query := `INSERT INTO codigo_ID (codigo, descripcion, estado, fecha_pedido, fecha_actualizacion)
		 VALUES (:codigo, :descripcion, :estado, :fecha_pedido, :fecha_actualizacion)`

	if err := r.db.Cargar(ctx, &nuevoCodigo, query); err != nil {
		return err
	}
	return nil

}

func (r *repo) ObtenerDetalle(ctx context.Context, codigo string) (*CodigoID, error) {

	if strings.TrimSpace(codigo) == "" {
		return nil, errParametrosNoValidos
	}

	var registro codigoIDDB

	query := `SELECT codigo, descripcion, estado, fecha_pedido, fecha_actualizacion
		 FROM codigo_ID
		 WHERE codigo = ?`

	err := r.db.ObtenerDetalle(ctx, &registro, query, codigo)
	if err != nil {
		return nil, err
	}

	codigoID := CodigoID{
		codigo:             registro.Codigo,
		descripcion:        registro.Descripcion,
		estado:             registro.Estado,
		fechaPedido:        registro.FechaPedido,
		fechaActualizacion: registro.FechaActualizacion,
	}
	return &codigoID, nil
}

func (r *repo) ActualizarEstado(ctx context.Context, codigo string, nuevoEstado string, fecha time.Time) error {

	if err := validarCodigo(codigo); err != nil {
		return err
	}

	query := `UPDATE codigo_ID
		 SET estado = ?, fecha_actualizacion = ?
		 WHERE codigo = ?`

	err := r.db.ModificarRegistro(ctx, query, nuevoEstado, fecha, codigo)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) FiltrarPorEstado(ctx context.Context, estado string) ([]CodigoID, error) {

	if strings.TrimSpace(estado) == "" {
		return nil, errParametrosNoValidos
	}

	var lista []codigoIDDB

	query := `SELECT codigo, descripcion, estado, fecha_pedido, fecha_actualizacion
		 FROM codigo_ID
		 WHERE estado = ?`

	err := r.db.ObtenerLista(ctx, &lista, query, estado)
	if err != nil {
		return nil, err
	}

	var codigos []CodigoID

	for _, registro := range lista {
		codigos = append(codigos, CodigoID{
			codigo:             registro.Codigo,
			descripcion:        registro.Descripcion,
			estado:             registro.Estado,
			fechaPedido:        registro.FechaPedido,
			fechaActualizacion: registro.FechaActualizacion,
		})
	}

	return codigos, nil
}

func (r *repo) Listar(ctx context.Context, limit, offset int) ([]CodigoID, error) {

	var lista []codigoIDDB

	query := `SELECT codigo, descripcion, estado, fecha_pedido, fecha_actualizacion
		 FROM codigo_ID
		LIMIT ? OFFSET ?`

	err := r.db.ObtenerLista(ctx, &lista, query, limit, offset)
	if err != nil {
		return nil, err
	}

	var codigos []CodigoID

	for _, registro := range lista {
		codigos = append(codigos, CodigoID{
			codigo:             registro.Codigo,
			descripcion:        registro.Descripcion,
			estado:             registro.Estado,
			fechaPedido:        registro.FechaPedido,
			fechaActualizacion: registro.FechaActualizacion,
		})
	}

	return codigos, nil
}
