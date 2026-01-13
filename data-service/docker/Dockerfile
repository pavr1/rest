FROM golang:1.25-alpine AS builder
WORKDIR /build
RUN apk add --no-cache git

# Copy shared module (go.mod has: replace shared => ../shared)
COPY shared/ /build/shared/

# Copy data-service into a subdirectory
WORKDIR /build/data-service
COPY data-service/go.mod data-service/go.sum ./
RUN go mod download

COPY data-service/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
WORKDIR /root/
COPY --from=builder /build/data-service/main .
EXPOSE 8086
CMD ["./main"]
