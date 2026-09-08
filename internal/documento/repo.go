package documento

import (
	"context"
	"strings"
)

type storageInterface interface {
	Cargar(ctx context.Context, strct any, query string) error
	ModificarRegistro(ctx context.Context, query string, arg ...interface{}) error
	ObtenerDetalle(ctx context.Context, strct any, query string, arg ...interface{}) error
	ObtenerLista(ctx context.Context, strct any, query string, arg ...interface{}) error
}

type documentoDB struct {
	Codigo        string  `db:"codigo"`
	Emision       string  `db:"emision"`
	Titulo        string  `db:"titulo"`
	Tipo          string  `db:"tipo"`
	UbicacionPath *string `db:"ubicacion_path"`
	BackupPath    *string `db:"backup_path"`
}

type repo struct {
	db storageInterface
}

func newRepo(db storageInterface) *repo {
	return &repo{db: db}
}

func (r *repo) Cargar(ctx context.Context, doc *Documento) error {

	if doc == nil {
		return errModelDeCargaNil
	}

	nuevoDocumento := documentoDB{
		Codigo:        doc.codigo,
		Emision:       doc.emision,
		Titulo:        doc.titulo,
		Tipo:          doc.tipo,
		UbicacionPath: doc.ubicacionPath,
		BackupPath:    doc.backupPath,
	}

	query := `INSERT INTO documento (codigo, emision, titulo, tipo, ubicacion_path, backup_path)
		 VALUES (:codigo, :emision, :titulo, :tipo, :ubicacion_path, :backup_path)`

	if err := r.db.Cargar(ctx, &nuevoDocumento, query); err != nil {
		return err
	}
	return nil
}

func (r *repo) ObtenerDetalle(ctx context.Context, codigo string) (*Documento, error) {

	if strings.TrimSpace(codigo) == "" {
		return nil, errParametrosNoValidos
	}

	registro := documentoDB{}

	query := `SELECT codigo, emision, titulo, tipo, ubicacion_path, backup_path
			  FROM documento
			  WHERE codigo = ?`

	err := r.db.ObtenerDetalle(ctx, &registro, query, codigo)
	if err != nil {
		return nil, err
	}

	doc := Documento{
		codigo:        registro.Codigo,
		emision:       registro.Emision,
		titulo:        registro.Titulo,
		tipo:          registro.Tipo,
		ubicacionPath: registro.UbicacionPath,
		backupPath:    registro.BackupPath,
	}
	return &doc, nil
}

func (r *repo) FiltrarPorTipo(ctx context.Context, tipo string) ([]Documento, error) {

	if strings.TrimSpace(tipo) == "" {
		return nil, errParametrosNoValidos
	}

	var lista []documentoDB

	query := `SELECT codigo, emision, titulo, tipo, ubicacion_path, backup_path
			  FROM documento
			  WHERE tipo = ?`

	err := r.db.ObtenerLista(ctx, &lista, query, tipo)
	if err != nil {
		return nil, err
	}

	var documentos []Documento

	for _, registro := range lista {
		documentos = append(documentos, Documento{
			codigo:        registro.Codigo,
			emision:       registro.Emision,
			titulo:        registro.Titulo,
			tipo:          registro.Tipo,
			ubicacionPath: registro.UbicacionPath,
			backupPath:    registro.BackupPath,
		})
	}

	return documentos, nil
}

func (r *repo) FiltrarPorTitulo(ctx context.Context, titulo string) ([]Documento, error) {

	if strings.TrimSpace(titulo) == "" {
		return nil, errParametrosNoValidos
	}

	var lista []documentoDB

	query := `SELECT codigo, emision, titulo, tipo, ubicacion_path, backup_path
			  FROM documento
			  WHERE titulo LIKE ?`

	err := r.db.ObtenerLista(ctx, &lista, query, titulo)
	if err != nil {
		return nil, err
	}

	var documentos []Documento

	for _, registro := range lista {
		documentos = append(documentos, Documento{
			codigo:        registro.Codigo,
			emision:       registro.Emision,
			titulo:        registro.Titulo,
			tipo:          registro.Tipo,
			ubicacionPath: registro.UbicacionPath,
			backupPath:    registro.BackupPath,
		})
	}

	return documentos, nil
}

func (r *repo) ActualizarPath(ctx context.Context, codigo string, nuevoPath string) error {

	if strings.TrimSpace(codigo) == "" {
		return errParametrosNoValidos
	}

	query := `UPDATE documento
		 SET ubicacion_path = ?
		 WHERE codigo = ?`

	err := r.db.ModificarRegistro(ctx, query, nuevoPath, codigo)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) ActualizarBackupPath(ctx context.Context, codigo string, nuevoPath string) error {

	if strings.TrimSpace(codigo) == "" {
		return errParametrosNoValidos
	}

	query := `UPDATE documento
		 SET backup_path = ?
		 WHERE codigo = ?`

	err := r.db.ModificarRegistro(ctx, query, nuevoPath, codigo)
	if err != nil {
		return err
	}

	return nil
}

type pendienteDB struct {
	ID          int    `db:"id"`
	Descripcion string `db:"descripcion"`
}

func (r *repo) PendientesRelacionados(ctx context.Context, codigoDocumento string) ([]pendiente, error) {

	if strings.TrimSpace(codigoDocumento) == "" {
		return nil, errParametrosNoValidos
	}

	query := `SELECT
				p.id,
				p.descripcion
			FROM pendientes p
			INNER JOIN ti_pendientes_documento t
				ON t.pendiente_id = p.id
			WHERE t.documento_id = ?;`

	var lista []pendienteDB

	err := r.db.ObtenerLista(ctx, &lista, query, codigoDocumento)
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
