FROM golang:1.22.1

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN go build -o ./dynamic_qr_server ./cmd/dynamic_qr/main.go
EXPOSE 8080
CMD ["./dynamic_qr_server"]
