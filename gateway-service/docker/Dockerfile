FROM golang:1.25-alpine AS builder
WORKDIR /build
RUN apk add --no-cache git

# Copy shared module (go.mod has: replace shared => ../shared)
COPY shared/ /build/shared/

# Copy gateway-service into a subdirectory
WORKDIR /build/gateway-service
COPY gateway-service/go.mod gateway-service/go.sum ./
RUN go mod download

COPY gateway-service/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
WORKDIR /root/
COPY --from=builder /build/gateway-service/main .
EXPOSE 8082
CMD ["./main"]

