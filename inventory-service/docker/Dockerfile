FROM golang:1.24-alpine AS builder
WORKDIR /build
RUN apk add --no-cache git

# Copy shared module (go.mod has: replace shared => ../shared)
COPY shared/ /build/shared/

# Copy inventory-service into a subdirectory
WORKDIR /build/inventory-service
COPY inventory-service/go.mod inventory-service/go.sum ./
RUN go mod download

COPY inventory-service/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
WORKDIR /root/
COPY --from=builder /build/inventory-service/main .
RUN chmod +x ./main
EXPOSE 8090
CMD ["./main"]
