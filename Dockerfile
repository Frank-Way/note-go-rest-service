FROM golang:alpine AS builder

WORKDIR /build
ADD go.mod .

COPY . .

RUN go build -o main cmd/main/main.go
FROM alpine
WORKDIR /build

COPY --from=builder /build/main /build/main
COPY config.yaml .

EXPOSE 5000

CMD ["./main"]
