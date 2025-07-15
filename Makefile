# Variables for the application container
APP_CONTAINER_IMAGE=my-go-app
APP_CONTAINER_NAME=go-app
APP_DOCKER_CONTEXT=.
APP_DOCKERFILE=./docker/app/Dockerfile
APP_ENV_FILE=.env
APP_PORT=1000

# Variables for the Redis container
REDIS_CONTAINER_IMAGE=my-redis-server
REDIS_CONTAINER_NAME=redis-server
REDIS_DOCKERFILE=./docker/redis/Dockerfile
REDIS_DOCKER_CONTEXT=.
REDIS_PORT=6379

# Variables for the PostgreSQL container
POSTGRES_CONTAINER_IMAGE=my-postgres-server
POSTGRES_CONTAINER_NAME=postgres-server
POSTGRES_DOCKER_CONTEXT=./docker/postgres
POSTGRES_DOCKERFILE=./docker/postgres/Dockerfile
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=P@ssw0rd
POSTGRES_DB=payment_service

# Network for the application and RabbitMQ containers
NETWORK=app-network

## ---- Development Commands ----
# Install dependencies
tidy:
	@echo -e "Running go mod tidy..."
	@go mod tidy

# Run the application in development mode
run:
	@echo -e "Running the application..."
	@dotenv -e .env -- go run ./cmd/main.go

# Test the application
test:
	@echo -e "Running tests..."
	@dotenv -e .env -- go test -v ./tests/...




## ---- Docker related targets ----
# Create a Docker network if it doesn't exist
docker-create-network:
	docker network inspect $(NETWORK) >NUL 2>&1 || docker network create $(NETWORK)

docker-remove-network:
	docker network rm $(NETWORK)




## --- PostgreSQL related targets ---
# Build PostgreSQL Docker image
docker-build-postgres:
	docker build -f $(POSTGRES_DOCKERFILE) -t $(POSTGRES_CONTAINER_IMAGE) $(POSTGRES_DOCKER_CONTEXT)

# Run PostgreSQL container
docker-run-postgres:
	docker run --name $(POSTGRES_CONTAINER_NAME) --network $(NETWORK) -p $(POSTGRES_PORT):$(POSTGRES_PORT) \
	-e POSTGRES_DB=$(POSTGRES_DB) \
	-e POSTGRES_USER=$(POSTGRES_USER) \
	-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	-d $(POSTGRES_CONTAINER_IMAGE)

# Build and run PostgreSQL container
docker-build-run-postgres: docker-build-postgres docker-run-postgres

# Remove PostgreSQL container
docker-remove-postgres:
	docker stop $(POSTGRES_CONTAINER_NAME)
	docker rm $(POSTGRES_CONTAINER_NAME)




## --- Redis related targets ---
# Build Redis Docker image
docker-build-redis:
	docker build -f $(REDIS_DOCKERFILE) -t $(REDIS_CONTAINER_IMAGE) $(REDIS_DOCKER_CONTEXT)

# Run Redis container
docker-run-redis:
	docker run --name $(REDIS_CONTAINER_NAME) --network $(NETWORK) -p $(REDIS_PORT):$(REDIS_PORT) \
	-d $(REDIS_CONTAINER_IMAGE)

# Build and run Redis container
docker-build-run-redis: docker-build-redis docker-run-redis

# Remove Redis container
docker-remove-redis:
	docker stop $(REDIS_CONTAINER_NAME)
	docker rm $(REDIS_CONTAINER_NAME)




## --- Application related targets ---
docker-build-app:
	docker build -f $(APP_DOCKERFILE) -t $(APP_CONTAINER_IMAGE) $(APP_DOCKER_CONTEXT)

docker-run-app:
	docker run --name $(APP_CONTAINER_NAME) --network $(NETWORK) -p $(APP_PORT):$(APP_PORT) \
	--env-file $(APP_ENV_FILE) \
	--link $(POSTGRES_CONTAINER_NAME):$(POSTGRES_CONTAINER_NAME) \
	--link $(REDIS_CONTAINER_NAME):$(REDIS_CONTAINER_NAME) \
	-v cert:/app/cert \
	-v keys:/app/keys \
	-v logs:/app/logs \
	-d $(APP_CONTAINER_IMAGE)

# Build and run the application container
docker-build-run-app: docker-build-app docker-run-app

docker-remove-app:
	docker stop $(APP_CONTAINER_NAME)
	docker rm $(APP_CONTAINER_NAME)

docker-up: docker-create-network \
	docker-build-run-postgres \
	docker-build-run-redis \
	docker-build-run-app

docker-down: docker-remove-app \
	docker-remove-redis \
	docker-remove-postgres \
	docker-remove-network

.PHONY: tidy run test \
	docker-create-network docker-remove-network \
	docker-build-postgres docker-run-postgres docker-build-run-postgres docker-remove-postgres \
	docker-build-redis docker-run-redis docker-build-run-redis docker-remove-redis \
	docker-build-app docker-run-app docker-build-run-app docker-remove-app \
	docker-up docker-down