package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

type gestorDbSQLite struct {
	db *sqlx.DB
}

func NewGestorDb(ctx context.Context, dataBasePath string) (*gestorDbSQLite, error) {
	db, err := sqlx.Open("sqlite3", dataBasePath+"?_foreign_keys=on") // "./database/app.db"
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return &gestorDbSQLite{db: db}, nil
}

func (g *gestorDbSQLite) Close() error {
	return g.db.Close()
}

// Carga en la db. El campo "strct" debe ser un puntero a un struct con las etiquetas db
func (g *gestorDbSQLite) Cargar(ctx context.Context, strct any, query string) (int, error) {

	result, err := g.db.NamedExecContext(ctx, query, strct)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return -1, ErrRegistroYaExistente
			}
		}
		return -1, fmt.Errorf("error inesperado: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error obteniendo ID del registro creado: %v", err)
	}

	return int(id), nil
}

// Obtiene un elemento en la db. El campo "strct" debe ser un puntero a un struct con las etiquetas db
func (g *gestorDbSQLite) ObtenerDetalle(ctx context.Context, strct any, query string, arg ...interface{}) error {

	err := g.db.GetContext(ctx, strct, query, arg...)

	if err != nil {
		if err == sql.ErrNoRows {
			return ErrRegistroNoEncontrado
		}
		return fmt.Errorf("error inesperado: %v", err)
	}

	return nil
}

// Obtiene una lista de elemento en la db. El campo "strct" debe ser un puntero a un struct con las etiquetas db
func (g *gestorDbSQLite) ObtenerLista(ctx context.Context, strct any, query string, arg ...interface{}) error {

	err := g.db.SelectContext(ctx, strct, query, arg...)
	if err != nil {
		return fmt.Errorf("error inesperado: %v", err)
	}

	return nil
}

// Modifica un registro en la db.
func (g *gestorDbSQLite) ModificarRegistro(ctx context.Context, query string, arg ...interface{}) error {

	result, err := g.db.ExecContext(ctx, query, arg...)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {

			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintForeignKey {
				return ErrRegistroNoEncontrado
			}

		}
		return fmt.Errorf("Error inesperado, detalle: %w", err)

	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error inesperado, detalle: %w", err)
	}
	if rows == 0 {
		return ErrNingunRegistroEncontrado
	}

	return nil
}
