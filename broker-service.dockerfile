FROM golang:1.18-alpine AS builder

RUN mkdir -p /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

RUN chmod +x /app/brokerApp

#build a tiny docker image
FROM alpine:latest

RUN mkdir /app 

COPY --from=builder /app/brokerApp /app

CMD ["/app/brokerApp"]