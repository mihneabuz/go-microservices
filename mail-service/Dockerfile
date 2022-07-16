FROM alpine:latest

RUN mkdir /app

COPY ./api /app
COPY templates /app/templates

WORKDIR /app

CMD [ "/app/api" ]
