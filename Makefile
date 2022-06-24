FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LOGGER_BINARY=loggerApp
MAILER_BINARY=mailerApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	sudo docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_auth build_logger build_mailer
	@echo "Stopping docker images (if running...)"
	sudo docker-compose down
	@echo "Building (when required) and starting docker images..."
	sudo docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	sudo docker-compose down
	@echo "Done!"

## build_mailer: builds the mailer binary as a linux executable
build_mailer:
	@echo "Building mailer binary..."
	cd ./mail-svc/bin && env GOOS=linux CGO_ENABLED=0 go build -o ${MAILER_BINARY} ../cmd
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ./broker-svc/bin && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ../cmd
	@echo "Done!"

## build_auth: builds the auth binary as a linux executable
build_auth:
	@echo "Building authentication binary..."
	cd ./authentication-svc/bin && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ../cmd
	@echo "Done!"

build_logger:
	@echo "Building logger binary..."
	cd ./logger-svc/bin && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ../cmd
	@echo "Done!"

## build_front: builds the frone end binary
build_front:
	@echo "Building front end binary..."
	cd ./front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

## start: starts the front end
start: build_front
	@echo "Starting front end"
	cd ./front-end && ./${FRONT_END_BINARY} &

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"[jerome@thearch micro]$ 