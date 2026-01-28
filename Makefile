# Variables
BINARY_NAME = payly
BUILD_DIR = ./bin
REMOTE_SERVER = payly
REMOTE_PATH = /var/www/payly/api/
ENV_FILE = .env
SERVICE_NAME = payly.service

# Default target
all: build deploy

# Run the app in development mode
dev:
	go run main.go

# Build the application binary
build:
	@echo "Building the application..."
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)

# Deploy the binary and .env file to the remote server
deploy: build
	@echo "Stopping the service on remote server..."
	ssh -t $(REMOTE_SERVER) "sudo systemctl stop $(SERVICE_NAME)"
	@echo "Deploying binary and .env file to $(REMOTE_SERVER)..."
	scp $(BUILD_DIR)/$(BINARY_NAME) $(ENV_FILE) $(REMOTE_SERVER):$(REMOTE_PATH)
	@echo "Starting the service on remote server..."
	ssh -t $(REMOTE_SERVER) "sudo systemctl start $(SERVICE_NAME)"
	@echo "Deployment complete."

# Deploy all (API + client) using root Makefile
deploy-all:
	$(MAKE) -C .. deploy

# Clean the build artifacts
clean:
	rm -rf $(BUILD_DIR)

# Help target
help:
	@echo "Usage:"
	@echo "  make dev        - Run the app in development mode"
	@echo "  make build      - Build the application binary"
	@echo "  make deploy     - Deploy binary and .env file to the remote server, and restart service"
	@echo "  make deploy-all - Deploy both API and client (uses root Makefile)"
	@echo "  make clean      - Remove build artifacts"
