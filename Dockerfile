FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o snippetbox ./cmd/web

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/snippetbox .
COPY --from=builder /app/ui ./ui
COPY --from=builder /app/security ./security
COPY --from=builder /app/internal ./internal

EXPOSE 4000
CMD ["./snippetbox"]