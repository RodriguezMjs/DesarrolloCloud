# DesarrolloCloud
Backend repo for cloud development project

1. To build the app and create the containers run 

```yaml
docker compose up -d --build
```

This command build the app and build the container with the PostgreSql Database for the project application

2. Estructura deseada del proyecto:

```text
project-root/
│
├── cmd/                        # Punto de entrada de la app
│   └── api/
│       └── main.go
│
├── internal/                   # Código privado (core del sistema)
│
│   ├── domain/                # Dominio 
│   │   ├── entities/
│   │   │   ├── user.go
│   │   │   ├── role.go
│   │   │   ├── course.go
│   │   │   ├── project.go
│   │   │   ├── assignment.go        # vinculación
│   │   │   ├── task.go
│   │   │   └── week.go
│   │   │
│   │   ├── valueobjects/
│   │   │   ├── hours.go
│   │   │   ├── task_status.go
│   │   │   └── week_range.go
│   │   │
│   │   ├── services/          # Domain Services (reglas complejas)
│   │   │   ├── hours_validator.go
│   │   │   ├── report_policy.go
│   │   │   └── week_service.go
│   │   │
│   │   ├── repositories/      # Interfaces (NO implementación)
│   │   │   ├── user_repository.go
│   │   │   ├── task_repository.go
│   │   │   ├── assignment_repository.go
│   │   │   └── report_repository.go
│   │   │
│   │   └── events/            # Eventos de dominio
│   │       ├── task_created.go
│   │       ├── report_generated.go
│   │       └── user_not_reported.go
│   │
│   ├── application/           # Casos de uso (use cases)
│   │   ├── usecases/
│   │   │   ├── create_task.go
│   │   │   ├── update_task.go
│   │   │   ├── generate_reports.go
│   │   │   ├── assign_user.go
│   │   │   └── login.go
│   │   │
│   │   ├── dto/               # Request / Response models
│   │   │   ├── task_dto.go
│   │   │   ├── user_dto.go
│   │   │   └── report_dto.go
│   │   │
│   │   └── services/          # Orquestación (Application Services)
│   │       ├── task_app_service.go
│   │       ├── report_app_service.go
│   │       └── notification_service.go
│   │
│   ├── infrastructure/        # 🔌 Implementaciones externas
│   │   ├── persistence/
│   │   │   ├── postgres/
│   │   │   │   ├── user_repository.go
│   │   │   │   ├── task_repository.go
│   │   │   │   └── assignment_repository.go
│   │   │
│   │   ├── http/              # Framework (Gin)
│   │   │   ├── handlers/
│   │   │   │   ├── task_handler.go
│   │   │   │   ├── user_handler.go
│   │   │   │   └── report_handler.go
│   │   │   │
│   │   │   ├── middleware/
│   │   │   │   ├── auth_middleware.go
│   │   │   │   └── role_middleware.go
│   │   │   │
│   │   │   └── router.go
│   │   │
│   │   ├── ai/                # Integración IA (reportes)
│   │   │   └── report_generator.go
│   │   │
│   │   ├── pdf/
│   │   │   └── pdf_generator.go
│   │   │
│   │   ├── queue/             # Asynq (opcional)
│   │   │   ├── client.go
│   │   │   └── worker.go
│   │   │
│   │   └── notifications/
│   │       └── notification_service.go
│   │
│   └── interfaces/            # Adaptadores (entrada)
│       └── api/
│           ├── controllers/
│           │   ├── task_controller.go
│           │   ├── user_controller.go
│           │   └── report_controller.go
│           │
│           └── routes.go
│
├── pkg/                       # Librerías reutilizables (opcional)
│   ├── logger/
│   └── utils/
│
├── configs/                   # Configuración
│   ├── config.yaml
│   └── database.go
│
├── migrations/                # SQL (PostgreSQL)
│
├── deployments/               # Docker / Nginx
│   ├── Dockerfile
│   └── docker-compose.yml
│
├── scripts/                   # Scripts auxiliares
│
├── go.mod
└── README.md
```
