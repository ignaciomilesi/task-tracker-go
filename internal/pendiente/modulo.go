package pendiente

import (
	"task-tracker-go/internal/pendiente/handler"
	"task-tracker-go/internal/pendiente/repo"
	"task-tracker-go/internal/pendiente/service"
)

type modulo struct {
	Repo             *repo.Repo
	ServiceAdjunto   *service.ServiceAdjunto
	ServiceAvance    *service.ServiceAvance
	ServicePendiente *service.ServicePendiente
	HandlerAdjunto   *handler.HandlerAdjunto
	HandlerAvance    *handler.HandlerAvance
	HandlerPendiente *handler.HandlerPendiente
}

func NuevoModulo(db repo.StorageInterface) *modulo {

	repo := repo.NewRepo(db)
	ServiceAdjunto := service.NewServiceAdjunto(repo)
	ServiceAvance := service.NewServiceAvance(repo)
	ServicePendiente := service.NewServicePendiente(repo, repo, repo)
	HandlerAdjunto := handler.NewHandlerAdjunto(ServiceAdjunto)
	HandlerAvance := handler.NewHandlerAvance(ServiceAvance)
	HandlerPendiente := handler.NewHandlerPendiente(ServicePendiente)

	return &modulo{
		Repo:             repo,
		ServiceAdjunto:   ServiceAdjunto,
		ServiceAvance:    ServiceAvance,
		ServicePendiente: ServicePendiente,
		HandlerAdjunto:   HandlerAdjunto,
		HandlerAvance:    HandlerAvance,
		HandlerPendiente: HandlerPendiente,
	}
}
