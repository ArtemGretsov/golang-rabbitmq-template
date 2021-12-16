FROM golang:alpine AS builder
WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/server/main.go

FROM scratch
WORKDIR /root/
COPY certs/* /etc/ssl/certs/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/config/ ./config
COPY --from=builder /app/app ./
CMD ["./app"]