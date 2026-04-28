FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src/ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o server ./cmd/api

FROM alpine:3.20

RUN apk --no-cache add ca-certificates wget

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
