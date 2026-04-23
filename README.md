# task-tracker-go
Sistema simple de seguimiento de pendientes en Go con base de datos SQLite


### Estructura del proyecto
```text
api/
  └── main.go              # Entry point

database/
  ├── app.db               # Base datos
  ├── esquema.sql          # Esquema de la base de datos
  └── init.go              # Genera la base de datos

internal/
  ├── handlers/            # HTTP (Gin)
  ├── services/            # Lógica de negocio
  ├── repositories/
  │     └── db_manager/    # Acceso a datos
  ├── models/              # Entidades 
  └── appErrors/           # Errores customs

config/                    # Configuración
```