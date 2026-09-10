FROM golang:1.27.1-alpine3.24 AS base
WORKDIR /app
RUN apk add --no-cache ca-certificates git

FROM base AS dev
RUN go install github.com/air-verse/air@latest
COPY . .
RUN go mod tidy
EXPOSE 3000
CMD ["air", "-c", ".air.toml"]

FROM base AS builder
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/go-api ./cmd/api

FROM alpine:3.24 AS prod
WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin/go-api ./go-api

USER appuser

EXPOSE 3000

CMD ["./go-api"]