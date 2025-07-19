# Axe Server - Microservices Docker Setup

This repository contains a complete microservices architecture for the Axe transport management system with a **database-per-service** pattern for better isolation and scalability.

## Architecture

### Services

- **Routing Service** (Port 3001) - Journey planning and route optimization
- **VTC Service** (Port 8080) - Vehicle for hire service
- **Data API** (Port 8000) - FastAPI service for data operations

### Databases (One per Service)

- **Routing DB** (Port 5432) - PostgreSQL with PostGIS for routing service
- **VTC DB** (Port 5433) - PostgreSQL for VTC service data
- **Data DB** (Port 5434) - PostgreSQL for data service operations

### Development Tools

- **PGAdmin** (Port 5050) - Database administration (development only)
- **Nginx** (Port 80) - API Gateway (development proxy)

## Quick Start

### Development Environment

```bash
# Clone and setup
git clone <repository-url>
cd axe-server

# Setup environment file
cp .env.example .env

# Start development services with hot reloading
make dev

# Or start in background
make dev-detached

# Initialize databases with schema and sample data
make init-dbs
```

### Production Environment

```bash
# Start production services
make prod

# Start with monitoring
make prod-with-monitoring
```

## Database Management

### Per-Service Database Commands

```bash
# Initialize all databases
make init-dbs

# Reset all databases (destructive)
make reset-dbs

# Backup all databases
make backup-dbs

# Connect to specific databases
make connect-routing-db
make connect-vtc-db
make connect-data-db

# View database logs
make logs-routing-db
make logs-vtc-db
make logs-data-db
```

### PGAdmin Access

In development mode, PGAdmin is available at:
- **URL**: http://localhost:5050
- **Email**: admin@axe.com
- **Password**: admin

All three databases are pre-configured in PGAdmin for easy access.

## Available Commands

Run `make help` to see all available commands:

```bash
make help
```

Key commands:
- `make dev` - Start development environment
- `make prod` - Start production environment
- `make logs` - View all service logs
- `make init-dbs` - Initialize all databases
- `make reset-dbs` - Reset all databases
- `make clean` - Clean up Docker resources

## Service URLs

After starting the services:

### APIs
- **Routing API**: http://localhost:3001
- **VTC API**: http://localhost:8080
- **Data API**: http://localhost:8000

### Development Gateway (Nginx)
- **API Gateway**: http://localhost (routes to all services)
  - `/api/routing/` -> Routing Service
  - `/api/vtc/` -> VTC Service
  - `/api/data/` -> Data API

### Databases (Development)
- **Routing DB**: postgres://routing_user:routing_pass@localhost:5432/routing_db
- **VTC DB**: postgres://vtc_user:vtc_pass@localhost:5433/vtc_db
- **Data DB**: postgres://data_user:data_pass@localhost:5434/data_db

### Management Tools
- **PgAdmin**: http://localhost:5050 (development only)

## Environment Configuration

Copy `.env.example` to `.env` and modify as needed:

```bash
cp .env.example .env
```

The setup supports per-service database configuration:

```bash
# Routing Service Database
ROUTING_DB_HOST=routing-db
ROUTING_DB_USER=routing_user
ROUTING_DB_PASSWORD=routing_pass
ROUTING_DB_NAME=routing_db

# VTC Service Database
VTC_DB_HOST=vtc-db
VTC_DB_USER=vtc_user
VTC_DB_PASSWORD=vtc_pass
VTC_DB_NAME=vtc_db

# Data Service Database
DATA_DB_HOST=data-db
DATA_DB_USER=data_user
DATA_DB_PASSWORD=data_pass
DATA_DB_NAME=data_db
```

## Docker Profiles

The setup supports different profiles for different use cases:

### Core Services (default)
```bash
docker-compose up -d
```

### Development
```bash
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d
```

### Production
```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

### With Monitoring
```bash
docker-compose --profile monitoring up -d
```

### With Reverse Proxy
```bash
docker-compose --profile proxy up -d
```

### With Caching
```bash
docker-compose --profile cache up -d
```

## Health Checks

All services include health checks. Check status with:

```bash
docker-compose ps
# or
make health
```

## Database Management

### Connect to database
```bash
make db-shell
```

### Backup database
```bash
make backup-db
```

### Restore database
```bash
make restore-db FILE=backup_20240101_120000.sql
```

## Development Features

### Hot Reloading

Development environment includes hot reloading using [Air](https://github.com/cosmtrek/air):
- Code changes automatically trigger rebuilds
- No need to restart containers during development

### Development Tools

When running in development mode, additional tools are available:
- PgAdmin at http://localhost:5050
- Enhanced logging and debugging
- Volume mounts for live code editing

## Monitoring

Enable monitoring with:

```bash
make prod-with-monitoring
```

This provides:
- **Prometheus** - Metrics collection
- **Grafana** - Metrics visualization
- **Node Exporter** - System metrics
- **cAdvisor** - Container metrics

## Scaling

Production compose includes scaling configuration:

```bash
# Scale routing service to 3 replicas
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d --scale routing-service=3
```

## Troubleshooting

### Check service logs
```bash
make logs
make logs-routing
make logs-vtc
make logs-db
```

### Restart specific service
```bash
docker-compose restart routing-service
```

### Rebuild specific service
```bash
docker-compose up -d --build routing-service
```

### Clean everything and start fresh
```bash
make clean
make dev
```

## API Testing

Test the APIs with:

```bash
# Test routing service
curl "http://localhost:3001/routing/best?fromLat=36.75&fromLng=3.05&toLat=36.76&toLng=3.06"

# Test VTC service
curl "http://localhost:8080/api/estimate"

# Or use the make command
make test-api
```

## Security Notes

### For Production:

1. **Change default passwords** in `.env` file
2. **Use secrets management** for sensitive data
3. **Enable TLS/SSL** for external access
4. **Configure firewall rules** appropriately
5. **Regular security updates** for base images

### Network Security:

Services communicate through isolated Docker networks. Only necessary ports are exposed to the host.

## Contributing

1. Create feature branch
2. Make changes
3. Test with `make test`
4. Submit pull request

## License

[Your License Here]
