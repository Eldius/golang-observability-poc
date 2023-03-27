# builder image
FROM golang:1.20-alpine as builder

WORKDIR /app
COPY . /app

ENV CGO_ENABLED=0

RUN apk add --no-cache git make build-base
RUN go build -o /app/app -v -a -ldflags '-extldflags "-static"' ./cmd/
RUN ls -lha
RUN chmod +x /app/app

#FROM alpine
FROM gcr.io/distroless/base

EXPOSE 8080

WORKDIR /app

COPY --chown=0:0 --from=builder /app/app /app
COPY rest-api-config.yaml /app/rest-api-config.yaml
COPY db/migrations /app/db/migrations

ENTRYPOINT [ "./app", "--config", "/app/rest-api-config.yaml" ]
