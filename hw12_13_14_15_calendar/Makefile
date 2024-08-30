BIN := "./bin/calendar"
SCHEDULE_BIN := "./bin/calendar_scheduler"
SENDER_BIN := "./bin/calendar_sender"
DOCKER_CALENDAR_IMG="calendar:develop"
DOCKER_SCHEDULER_IMG="scheduler:develop"
DOCKER_SENDER_IMG="sender:develop"
DOCKER_INT_IMG="integration:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar && \
	go build -v -o $(SCHEDULE_BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler && \
	go build -v -o $(SENDER_BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

run-calendar: build
	$(BIN) -config ./configs/calendar_config.yml

run-scheduler: build
	$(SCHEDULE_BIN) -config ./configs/scheduler_config.yml

run-sender: build
	$(SENDER_BIN) -config ./configs/sender_config.yml

run: build
	$(BIN) -config ./configs/calendar_config.yml && \
	$(SCHEDULE_BIN) -config ./configs/scheduler_config.yml && \
	$(SENDER_BIN) -config ./configs/sender_config.yml

run-all: build
	# docker run --name postgres -e POSTGRES_PASSWORD=1234 -p 5432:5432 postgres & \
	# PGPASSWORD=1234 psql -h localhost -p 5432 -U postgres -a -f ./migrations/init.sql postgres & \
	# docker run -d --hostname my-rabbit --name some-rabbit -p 8080:15672 -p 5672:5672 rabbitmq:3-management & \
	$(BIN) -config ./configs/calendar_config.yml & \
	$(SCHEDULE_BIN) -config ./configs/scheduler_config.yml

build-calendar-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_CALENDAR_IMG) \
		-f build/CalendarDockerfile .

build-scheduler-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_SCHEDULER_IMG) \
		-f build/SchedulerDockerfile .

build-sender-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_SENDER_IMG) \
		-f build/SenderDockerfile .

build-int-img:
	docker build -t $(DOCKER_INT_IMG) -f build/IntegrationDockerfile .

build-img: build-calendar-img build-scheduler-img build-sender-img

up: build-img
	docker-compose -f docker-compose.yml up -d

down: up
	docker-compose -f docker-compose.yml down

integration-tests: up build-int-img
	docker-compose -f docker-compose-integrations.yml up && \
	docker-compose -f docker-compose.yml down --remove-orphans

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.60.1

lint: install-lint-deps
	golangci-lint run ./...

lint-fix: install-lint-deps
	golangci-lint run --fix ./...

generate:
	rm -rf ./internal/server/grpc/pb
	mkdir -p ./internal/server/grpc/pb

	protoc \
		--proto_path=api/grpc/ \
		--go_out=internal/server/grpc/pb \
		--go-grpc_out=internal/server/grpc/pb \
		api/grpc/*.proto

.PHONY: build run run-all build-img run-img version test lint generate
