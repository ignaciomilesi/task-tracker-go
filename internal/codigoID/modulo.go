package codigoID

type modulo struct {
	Repo     *repo
	Service  *service
	Handlers *handler
}

func NuevoModulo(db storageInterface) *modulo {

	repo := newRepo(db)
	service := newService(repo)
	handlers := newHandler(service)

	return &modulo{
		Repo:     repo,
		Service:  service,
		Handlers: handlers,
	}
}
