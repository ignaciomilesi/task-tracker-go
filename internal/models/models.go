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

type Pendientes struct {
	ID          int    `db:"id"`
	Titulo      string `db:"titulo"`
	Descripcion string `db:"descripcion"`

	SolicitanteID int       `db:"solicitante_id"`
	FechaPedido   time.Time `db:"fecha_pedido"`

	AsignadoID    *int       `db:"asignado_id"` //colaborador
	FechaAsignado *time.Time `db:"fecha_asignado"`
	Finalizado    bool       `db:"finalizado"`

	Cierre      *string    `db:"cierre"`
	FechaCierre *time.Time `db:"fecha_cierre"`

	IdentificacionTablaPendiente *string `db:"identificacion_tabla_pendiente"`
}

type Avance struct {
	ID          int       `db:"id"`
	PendienteID int       `db:"pendiente_id"`
	Descripcion string    `db:"descripcion"`
	Fecha       time.Time `db:"fecha"`
	MailPath    *string   `db:"mail_path"`
}

type Adjunto struct {
	ID          int    `db:"id"`
	PendienteID int    `db:"pendiente_id"`
	Descripcion string `db:"descripcion"`
	ArchivoPath string `db:"archivo_path"`
}
