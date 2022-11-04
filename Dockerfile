FROM golang:1.18-alpine as builder

WORKDIR /app
COPY ./ .

RUN go mod download
RUN go mod verify
RUN go build -o matic-jagar

FROM golang:1.18-alpine

RUN mkdir -p /app
WORKDIR /app

COPY --from=builder /app/matic-jagar /app

CMD ["./matic-jagar"]