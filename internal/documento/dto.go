package documento

type nuevoDocumentoRequest struct {
	Codigo        string
	Emision       string
	Titulo        string
	Tipo          string
	UbicacionPath *string
	BackupPath    *string
}

type pendiente struct {
	ID          int
	Descripcion string
}
