FROM golang:alpine

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o ./cmd/tender-service-api/main ./cmd/tender-service-api

CMD ["./cmd/tender-service-api/main"]
