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
	estado             string     `db:"estado"`
	fechaPedido        time.Time  `db:"fecha_pedido"`
	fechaActualizacion *time.Time `db:"fecha_actualizacion"`
}
