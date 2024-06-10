FROM golang:1.20

WORKDIR /app/

COPY . .

RUN apt-get update && apt-get install -y librdkafka-dev

RUN go build -o ./cmd/walletcore/main ./cmd/walletcore/

CMD ["./cmd/walletcore/main"]