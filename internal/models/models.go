package models

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
