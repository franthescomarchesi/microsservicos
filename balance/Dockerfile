FROM golang:1.20

WORKDIR /app/

COPY . .

RUN apt-get update && apt-get install -y librdkafka-dev

RUN go build -o ./cmd/balance/main ./cmd/balance/

CMD ["./cmd/balance/main"]