# TimeWise DMS

## Prerequisites

- Go v1.22.5

### Installing

1. Install dependencies
```bash
go mod download
```

2. Copy the `.env.example` file to `.env` and fill in the necessary information.
```bash
cp .env.example .env
```

3. Run the application (server will be running on port `8080`)
```bash
go run main.go
```