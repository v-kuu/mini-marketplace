FROM golang:1.25.3-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN apk add --no-cache build-base
COPY internal ./internal
COPY cmd ./cmd
RUN CGO_ENABLED=1 GOOS=linux go build -o api ./cmd/server

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/api .
COPY products.db ./products.db
RUN apk --no-cache add wget
EXPOSE 8080
CMD ["./api"]
