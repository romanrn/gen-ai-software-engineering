FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY main.go .
COPY main_test.go .
RUN go build -o server .
# Allow docker compose run --rm app go test ./... to work against this stage
CMD ["./server"]

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]