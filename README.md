# Reborn

Currently look into [this discussion](https://github.com/oj-lab/reborn/discussions/1) for knowledge about the project.

## Docker Usage

### Building the Docker Image

```bash
# Build the Docker image
make docker-build

# Or build with docker-compose
make docker-compose-build
```

### Running with Docker

```bash
# Run the container directly
make docker-run

# Or run with docker-compose (recommended)
make docker-compose-up
```

The application will be available at `http://localhost:8080`.

### Docker Commands

- `make docker-build` - Build the Docker image
- `make docker-run` - Run the container directly
- `make docker-compose-up` - Start services with docker-compose
- `make docker-compose-down` - Stop services with docker-compose
- `make docker-clean` - Clean up Docker containers and images

### Configuration

The application uses the configuration file at `configs/default.toml`. You can modify this file or mount your own configuration when running the container.

## Contributing

See more in [CONTRIBUTING.md](CONTRIBUTING.md).
