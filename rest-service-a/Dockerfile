# builder image
FROM golang:1.21.5-alpine as builder

ARG COMMIT_HASH

ENV COMMIT_HASH=${COMMIT_HASH}
ENV GOCACHE=/root/.cache/go-build

WORKDIR /app
COPY . /app

ENV CGO_ENABLED=0

RUN apk add --no-cache git make build-base
RUN go build \
    -o /app/app \
    -v \
    -a \
    -ldflags '-extldflags "-static"' \
    -ldflags "-X 'github.com/eldius/rest-api/rest-service-a/internal/config.Version=${COMMIT_HASH}'" \
      ./cmd/

RUN chmod +x /app/app

FROM gcr.io/distroless/base

EXPOSE 8080

WORKDIR /app

COPY --chown=0:0 --from=builder /app/app /app
COPY rest-api-config.yaml /app/rest-api-config.yaml
COPY db/migrations /app/db/migrations

ENTRYPOINT [ "./app", "--config", "/app/rest-api-config.yaml" ]
