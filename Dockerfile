FROM golang:1.22-alpine3.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o app main.go

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/app .
CMD ["./app"]