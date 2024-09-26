FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o receipt-uploader

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/receipt-uploader .
EXPOSE 8080
CMD ["./receipt-uploader"]