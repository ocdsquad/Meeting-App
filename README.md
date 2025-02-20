# README

## Structure
```
project/
├── cmd/                   # Entry points (main files)
│   └── main.go            # Main application entry
├── configs/               # Configuration files
│   └── config.yaml
├── internal/              # Application-specific code
│   ├── domain/            # Core business logic (Entities and Interfaces)
│   │   ├── entity/        # Entity definitions
│   │   │   └── student.go
│   │   └── repository/    # Repository interfaces
│   │       └── student_repo.go
│   └── usecase/           # Use case implementations
│       ├── student_usecase.go
│       └── usecase.go     # Use case interface
├── presenter/             # Input and Output layers
│   ├── handler/           # API handlers
│   │   ├── handler_http.go
│   │   ├── handler.go     # Handler interface
│   │   └── docs/          # API documentation (Swagger, etc.)
│   └── swagger.yaml
│   ├── middleware/        # HTTP middleware (auth, logging, etc.)
│   └── presenter/         # Presenter (formatting the output)
│       ├── student_presenter.go
│       └── presenter.go   # Presenter interface
├── pkg/                   # Shared utility packages (e.g., database, logger, errors)
│   └── database/
│       └── postgres.go
└── go.mod                 # Dependency management

```

### Upload File
```
create folder upload in root folder.
```