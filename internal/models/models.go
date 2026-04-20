package models

import "time"

type Solicitante struct {
	ID     int    `db:"id"`
	Nombre string `db:"nombre"`
}

type Colaborador struct {
	ID     int    `db:"id"`
	Nombre string `db:"nombre"`
}

type CodigoSAP struct {
	Codigo      string  `db:"codigo"`
	Descripcion *string `db:"descripcion"`
}

type CodigoID struct {
	Codigo             string     `db:"codigo"`
	Descripcion        *string    `db:"descripcion"`
	Estado             string     `db:"estado"`
	FechaPedido        time.Time  `db:"fecha_pedido"`
	FechaActualizacion *time.Time `db:"fecha_actualizacion"`
}

type Documento struct {
	ID            int     `db:"id"`
	Codigo        string  `db:"codigo"`
	Titulo        string  `db:"titulo"`
	Tipo          string  `db:"tipo"`
	UbicacionPath *string `db:"ubicacion_path"`
}
