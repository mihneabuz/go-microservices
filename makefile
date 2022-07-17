
build: frontend broker auth logger mail listener

frontend:
	cd frontend && env CGO_ENABLED=0 go build ./cmd/web

broker:
	cd broker-service && env CGO_ENABLED=0 go build ./cmd/api

auth:
	cd auth-service && env CGO_ENABLED=0 go build ./cmd/api

logger:
	cd logger-service && env CGO_ENABLED=0 go build ./cmd/api

mail:
	cd mail-service && env CGO_ENABLED=0 go build ./cmd/api

listener:
	cd listener-service && env CGO_ENABLED=0 go build ./cmd/main.go

start: build
	docker-compose up --build

clean:
	rm frontend/web broker-service/api auth-service/api mail-service/api logger-service/api listener-service/main

.PHONY: frontend
