FROM golang:1.24-alpine AS builder

WORKDIR /src

COPY backend/go.mod ./
RUN go mod download

COPY backend/ ./
RUN go build -o /bin/api ./cmd/api

FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache wget

COPY --from=builder /bin/api /app/api
COPY backend/migrations /app/migrations

EXPOSE 8080

CMD ["/app/api"]
