FROM golang:1.24.5

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o user-service ./cmd/user-service
RUN go build -o product-service ./cmd/product-service
RUN go build -o order-service ./cmd/order-service
RUN go build -o inventory-service ./cmd/inventory-service
RUN go build -o payment-service ./cmd/payment-service

EXPOSE 8080 8081 8082 8083 8084

CMD ["/bin/sh", "-c", "if [ \"$DOCKER_ENV\" = \"user\" ]; then ./user-service; elif [ \"$DOCKER_ENV\" = \"product\" ]; then ./product-service; elif [ \"$DOCKER_ENV\" = \"order\" ]; then ./order-service; elif [ \"$DOCKER_ENV\" = \"payment\" ]; then ./payment-service; elif [ \"$DOCKER_ENV\" = \"inventory\" ]; then ./inventory-service; fi"]