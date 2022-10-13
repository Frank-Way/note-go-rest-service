FROM golang:alpine AS builder

WORKDIR /build
ADD go.mod .

COPY . .

RUN go build -o main cmd/main/main.go
RUN go build -o health cmd/health/health.go
FROM alpine
WORKDIR /build

COPY --from=builder /build/main /build/main
COPY --from=builder /build/health /build/health
COPY config.yaml .

HEALTHCHECK --interval=1s --timeout=1s --start-period=2s --retries=3 CMD [ "/healthcheck" ]

CMD ["./main"]
