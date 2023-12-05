FROM golang:1.18-alpine as builder

WORKDIR /app
COPY ./ .

RUN go mod download
RUN go mod verify
RUN go build -o matic-jagar

FROM golang:1.18-alpine

# Create a group and user for the Go application
RUN groupadd -g 10000 maticjager && \
    useradd -u 10000 -g maticjager -s /bin/false maticjager

RUN mkdir -p /app
WORKDIR /app

COPY --from=builder /app/matic-jagar /app

# Set the user and group ownership for the built binary
RUN chown -R maticjager:maticjager /app

# Use an unprivileged user.
USER maticjager

CMD ["./matic-jagar"]