# Build Step
FROM golang:1.18-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o broker-service ./cmd/api

RUN chmod +x /app/broker-service



# Run Step
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/broker-service /app

CMD [ "/app/broker-service" ]
