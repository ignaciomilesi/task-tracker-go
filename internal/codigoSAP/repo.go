package codigoSAP

import (
	"context"
	"strings"

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
type CodigoSAPDB struct {
	Codigo      string  `db:"codigo"`
	Descripcion *string `db:"descripcion"`
}

// repo y funciones
type repo struct {
	db storageInterface
}

func newRepo(db storageInterface) *repo {
	return &repo{db: db}
}

func (r *repo) Cargar(ctx context.Context, NuevoCodigoSap *CodigoSAP) error {

	if NuevoCodigoSap == nil {
		return errModelDeCargaNil
	}

	var nuevoCodigo = CodigoSAPDB{
		Codigo:      NuevoCodigoSap.codigo,
		Descripcion: NuevoCodigoSap.descripcion,
	}

	query := `INSERT INTO codigo_SAP (codigo, descripcion) 
		VALUES (:codigo, :descripcion)`

	if err := r.db.Cargar(ctx, &nuevoCodigo, query); err != nil {
		return err
	}

	return nil
}

func (r *repo) ObtenerDetalle(ctx context.Context, codigo string) (*CodigoSAP, error) {

	if strings.TrimSpace(codigo) == "" {
		return nil, errParametrosNoValidos
	}

	var registro CodigoSAPDB

	query := `SELECT codigo, descripcion 
		FROM codigo_SAP 
		WHERE codigo = ?`

	err := r.db.ObtenerDetalle(ctx, &registro, query, codigo)
	if err != nil {
		return nil, err
	}

	return &CodigoSAP{
		codigo:      registro.Codigo,
		descripcion: registro.Descripcion,
	}, nil
}

func (r *repo) BuscarPorDescripcion(ctx context.Context, parametro string) ([]CodigoSAP, error) {

	if strings.TrimSpace(parametro) == "" {
		return nil, errParametrosNoValidos
	}

	var lista []CodigoSAPDB

	query := `SELECT codigo, descripcion 
		FROM codigo_SAP 
		WHERE descripcion LIKE ?`

	err := r.db.ObtenerLista(ctx, &lista, query, "%"+parametro+"%")
	if err != nil {
		return nil, err
	}

	var codigos []CodigoSAP

	for _, registro := range lista {
		codigos = append(codigos, CodigoSAP{
			codigo:      registro.Codigo,
			descripcion: registro.Descripcion,
		})
	}

	return codigos, nil
}

func (r *repo) ModificarDescripcion(ctx context.Context, codigo string, nuevaDescripcion string) error {

	if err := validarCodigo(codigo); err != nil {
		return err
	}

	query := `UPDATE codigo_SAP
	 SET descripcion = ? 
	 WHERE codigo = ?`

	err := r.db.ModificarRegistro(ctx, query, nuevaDescripcion, codigo)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) Listar(ctx context.Context, limit, offset int) ([]CodigoSAP, error) {

	var lista []CodigoSAPDB

	query := `SELECT codigo, descripcion 
		FROM codigo_SAP ORDER BY codigo
		LIMIT ? OFFSET ?`

	err := r.db.ObtenerLista(ctx, &lista, query, limit, offset)
	if err != nil {
		return nil, err
	}

	var codigos []CodigoSAP

	for _, registro := range lista {
		codigos = append(codigos, CodigoSAP{
			codigo:      registro.Codigo,
			descripcion: registro.Descripcion,
		})
	}

	return codigos, nil
}

type pendienteDB struct {
	ID          int    `db:"id"`
	Descripcion string `db:"descripcion"`
}

func (r *repo) PendientesRelacionados(ctx context.Context, codigoSAP string) ([]pendiente, error) {

	if strings.TrimSpace(codigoSAP) == "" {
		return nil, errParametrosNoValidos
	}

	query := `SELECT
				p.id,
				p.descripcion
			FROM pendientes p
			INNER JOIN ti_pendientes_codigo_sap t
				ON t.pendiente_id = p.id
			WHERE t.codigo_sap_codigo = ?;`

	var lista []pendienteDB

	err := r.db.ObtenerLista(ctx, &lista, query, codigoSAP)
	if err != nil {
		return nil, err
	}

	var pendientes []pendiente

	for _, registro := range lista {
		pendientes = append(pendientes, pendiente{
			ID:          registro.ID,
			Descripcion: registro.Descripcion,
		})
	}

	return pendientes, nil
}
