package repo

import (
	"context"
	"task-tracker-go/internal/pendiente/dto"
	"task-tracker-go/internal/pendiente/model"
	"time"
)

type avanceDB struct {
	ID          int       `db:"id"`
	PendienteID int       `db:"pendiente_id"`
	Descripcion string    `db:"descripcion"`
	Fecha       time.Time `db:"fecha"`
	MailPath    *string   `db:"mail_path"`
}

func (c avanceDB) generarDesdeDomain(av model.Avance) avanceDB {
	c.ID = av.ID
	c.PendienteID = av.PendienteID
	c.Descripcion = av.Descripcion
	c.Fecha = av.Fecha
	c.MailPath = av.MailPath
	return c
}

func (d avanceDB) devolverDomain() model.Avance {
	return model.Avance{
		ID:          d.ID,
		PendienteID: d.PendienteID,
		Descripcion: d.Descripcion,
		Fecha:       d.Fecha,
		MailPath:    d.MailPath,
	}
}

func (r *Repo) CargarAvance(ctx context.Context, nuevoavance *model.Avance) (int, error) {

	if nuevoavance == nil {
		return -1, model.ErrModelDeCargaNil
	}

	var nuevaSolicitudAvance avanceDB
	nuevaSolicitudAvance.generarDesdeDomain(*nuevoavance)

	query := `INSERT INTO avance (pendiente_id, descripcion, fecha, mail_path)
		 VALUES (:pendiente_id, :descripcion, :fecha, :mail_path)`

	return r.db.Cargar(ctx, &nuevaSolicitudAvance, query)
}

func (r *Repo) EliminarAvance(ctx context.Context, id int) error {

	if id <= 0 {
		return model.ErrParametrosNoValidos
	}

	query := `DELETE FROM avance WHERE id = ?`

	err := r.db.ModificarRegistro(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) ObtenerListaAvances(ctx context.Context, pendienteID int) ([]dto.AvanceDetalleRespuesta, error) {

	if pendienteID <= 0 {
		return nil, model.ErrParametrosNoValidos
	}

	var lista []avanceDB

	query := `SELECT descripcion, fecha, mail_path
		 FROM avance
		 WHERE pendiente_id = ?
		 ORDER BY fecha DESC`

	err := r.db.ObtenerLista(ctx, &lista, query, pendienteID)
	if err != nil {
		return nil, err
	}

	var avances []dto.AvanceDetalleRespuesta

	for _, registro := range lista {
		avances = append(avances, dto.AvanceDetalleRespuesta{
			Descripcion: registro.Descripcion,
			Fecha:       registro.Fecha,
			MailPath:    registro.MailPath,
		})
	}

	return avances, nil
}
