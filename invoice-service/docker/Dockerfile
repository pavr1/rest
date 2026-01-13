# Invoice Service Dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /build
RUN apk add --no-cache git

# Copy shared module
COPY shared/ /build/shared/

# Copy invoice-service
WORKDIR /build/invoice-service
COPY invoice-service/go.mod invoice-service/go.sum ./
RUN go mod download

COPY invoice-service/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
WORKDIR /root/
COPY --from=builder /build/invoice-service/main .
EXPOSE 8092
CMD ["./main"]