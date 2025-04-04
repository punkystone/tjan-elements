FROM golang:alpine AS builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/elements -ldflags="-s -w" cmd/elements/main.go

FROM scratch
COPY --from=builder /build/bin/elements /elements
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/elements"]