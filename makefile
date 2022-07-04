
build: frontend broker auth

frontend:
	cd frontend && env CGO_ENABLED=0 go build ./cmd/web

broker:
	cd broker-service && env CGO_ENABLED=0 go build ./cmd/api

auth:
	cd auth-service && env CGO_ENABLED=0 go build ./cmd/api

start: build
	docker-compose up --build

clean:
	rm frontend/web broker-service/api

.PHONY: frontend
