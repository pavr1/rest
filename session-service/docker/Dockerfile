FROM golang:1.24-alpine AS builder
WORKDIR /build
RUN apk add --no-cache git

# Copy shared module (go.mod has: replace shared => ../shared)
COPY shared/ /build/shared/

# Copy data-service module (go.mod has: replace data-service => ../data-service)
COPY data-service/ /build/data-service/

# Copy session-service into a subdirectory
WORKDIR /build/session-service
COPY session-service/go.mod session-service/go.sum ./
RUN go mod download

COPY session-service/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
WORKDIR /root/
COPY --from=builder /build/session-service/main .
EXPOSE 8087
CMD ["./main"]
