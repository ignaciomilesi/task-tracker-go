package usuario

type usuario struct {
	id          int
	nombre      string
	colaborador bool
}

func NewUsuario(nombre string, colaborador bool) *usuario {
	return &usuario{
		nombre:      nombre,
		colaborador: colaborador,
	}
}
