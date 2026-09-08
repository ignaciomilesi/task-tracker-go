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
### Estructura de la base de datos

[link para visualizarla](https://www.drawdb.app/editor/diagrams/8a5e1b2c-69f8-4985-8ee2-a2ad752b8bb2)

### Objetivo del API

- Pendiente:
  - [x] [x] Crear [UseCase]
  - [x] [x] Asignar [UseCase]
  - [x] [x] Finalizar [UseCase]

  - [x] [x] Registrar código SAP con vinculación [UseCase]
  - [x] [x] Registrar código ID con vinculación [UseCase]
  - [x] [x] Registrar Documento con vinculación [UseCase]

  - [x] [x] Ver detalle completo de un pendiente [UseCase]

  - [x] [x] Registrar avance [Service]
  - [x] [x] Cargar Adjunto [Service]

- Por fuera de pendiente
  - [x] [x] Cargar nuevo colaborador [Service]
  - [x] [x] Cargar nuevo solicitante [Service]
  - [x] [x] Cargar código SAP sin vinculación [Service]
  - [x] [x] Cargar código ID sin vinculación [Service]
  - [x] [x] Cargar Documento sin vinculación [Service]


- [x] [x] Ver lista de pendientes no finalizados [Service]
- [x] Listar colaboradores [x], solicitante [x], SAP [x], ID [x], Documentos [x][Service]