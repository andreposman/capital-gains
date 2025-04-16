# Capital Gains Tax Calculator (Go)

This project implements a command-line interface (CLI) application in Go to calculate capital gains tax for stock market operations based on a series of transactions provided via standard input.


## Prerequisites

*   **If you want to run locally**
    * **Go:** Version 1.24.3 or higher.
    * **Make:** (Optional, for convenience using Makefile commands)
    

* Using containers
  *   **Docker:** 
  *   **Docker Compose:** 
  *   **Make:** (Optional, for convenience using Makefile commands)

## Setup

1.  **Clone the repository:**
    ```bash
    git clone <your-repository-url>
    cd <repository-directory>
    ```


## Makefile Commands

A Makefile is provided for convenience:

```bash
make run-local     # Run the application locally with input.txt
make build         # Build the application binary
make test          # Run all tests
make run           # Build and run the application via Docker, piping input.txt (std input from challenge)
make run2          # Build and run the application via Docker, piping input2.txt (which has more data)
make down          # Stop and clean up Docker Compose containers
make help          # Show available commands
```


### Running with Docker (recommended)

Uses Docker Compose to build the image and run the application in a container, piping input from input.txt. 
```bash
make run
```

To run with input2.txt which has more data:

```bash
make run2
```

To run manually with Docker Compose and a different input file:


#### Ensure the image is built
```bash
docker-compose build
```

#### Pipe input to the container

```bash
cat input.txt | docker-compose run --rm -T capital-gains
```



To stop and remove the container managed by Docker Compose (if it was run without --rm or detached):

```bash
make down
```

To run the unit tests:

```bash
make test
```
This command executes `go test -v ./....`



## Running Locally

Use Make command:

```bash
  make run-local
  #or
  make run-local2 # if you want an input with more data
```

To run manually you can use: 

```bash
go run cmd/main.go < input.txt
```

Or build the binary using: 
```bash
make build
```

and using the built binary run :
```bash
./bin/capital-gains < input.txt
```

## Project Structure

```bash
.
├── Dockerfile             # For building the Docker image
├── Makefile               # Convenience commands for build, run, test
├── README.md              # This file
├── bin/                   # Compiled binaries (created by 'make build')
├── cmd/                   # Main application entrypoint
│   └── main.go
├── docker-compose.yml     # Docker Compose configuration
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
├── input.txt              # Sample input file
├── internal/              # Internal application code (not reusable)
│   ├── application/       # Application logic/use cases (OperationProcessor)
│   ├── domain/            # Core business logic (Portfolio, Tax rules)
│   └── infra/             # Infrastructure concerns (CLI handler, JSON parsing)
└── pkg/                   # Shared library code (reusable, e.g., helpers)
└── helpers/
```

## Architecture

The project attempts to follow principles inspired by Clean Architecture and Domain-Driven Design (DDD):

* Domain (internal/domain): Contains the core business logic and state (Portfolio, Tax calculation rules), independent of other layers.


* Application (internal/application): Orchestrates the use cases (processing operations). Depends on Domain.


* Infrastructure (internal/infra): Handles external concerns like CLI interaction (stdin/stdout) and JSON parsing. Depends on Application.


* Cmd (cmd): The entry point that wires everything together.


* Pkg (pkg): Utility functions potentially reusable across projects.







    