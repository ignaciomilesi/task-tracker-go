package repo

import (
	"context"
	"task-tracker-go/internal/pendiente/dto"
	"task-tracker-go/internal/pendiente/model"
)

type adjuntoDB struct {
	ID          int    `db:"id"`
	PendienteID int    `db:"pendiente_id"`
	Descripcion string `db:"descripcion"`
	ArchivoPath string `db:"archivo_path"`
}

func (c adjuntoDB) generarDesdeDomain(av model.Adjunto) adjuntoDB {
	c.ID = av.ID
	c.PendienteID = av.PendienteID
	c.Descripcion = av.Descripcion
	c.ArchivoPath = av.ArchivoPath
	return c
}

func (d adjuntoDB) devolverDomain() model.Adjunto {
	return model.Adjunto{
		ID:          d.ID,
		PendienteID: d.PendienteID,
		Descripcion: d.Descripcion,
		ArchivoPath: d.ArchivoPath,
	}
}

func (r *Repo) CargarAdjunto(ctx context.Context, nuevoAdjunto *model.Adjunto) (int, error) {

	if nuevoAdjunto == nil {
		return -1, model.ErrModelDeCargaNil
	}

	var nuevaSolicitudAdjunto adjuntoDB
	nuevaSolicitudAdjunto.generarDesdeDomain(*nuevoAdjunto)

	query := `INSERT INTO adjunto (pendiente_id, descripcion, archivo_path)
		 VALUES (:pendiente_id, :descripcion, :archivo_path)`

	return r.db.Cargar(ctx, &nuevaSolicitudAdjunto, query)
}

func (r *Repo) EliminarAdjunto(ctx context.Context, id int) error {

	if id <= 0 {
		return model.ErrParametrosNoValidos
	}

	query := `DELETE FROM adjunto WHERE id = ?`

	err := r.db.ModificarRegistro(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) ObtenerListaAdjuntos(ctx context.Context, pendienteID int) ([]dto.AdjuntoDetalleRespuesta, error) {

	if pendienteID <= 0 {
		return nil, model.ErrParametrosNoValidos
	}
	var lista []adjuntoDB

	query := `SELECT descripcion, archivo_path
              FROM adjunto
              WHERE pendiente_id = ?`

	err := r.db.ObtenerLista(ctx, &lista, query, pendienteID)
	if err != nil {
		return nil, err
	}

	var adjuntos []dto.AdjuntoDetalleRespuesta

	for _, registro := range lista {
		adjuntos = append(adjuntos, dto.AdjuntoDetalleRespuesta{
			Descripcion: registro.Descripcion,
			ArchivoPath: registro.ArchivoPath,
		})
	}

	return adjuntos, nil
}
