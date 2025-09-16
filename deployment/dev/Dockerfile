FROM golang:alpine AS builder

WORKDIR /build

COPY . .

RUN go build -mod=vendor -o kado-server ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/kado-server .

EXPOSE 8000

ENTRYPOINT ["./kado-server"]
