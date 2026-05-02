FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o order-service ./cmd

CMD ["./order-service"]