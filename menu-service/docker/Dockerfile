FROM golang:1.25-alpine AS builder
WORKDIR /build
RUN apk add --no-cache git

# Copy shared module (go.mod has: replace shared => ../shared)
COPY shared/ /build/shared/

# Copy menu-service into a subdirectory
WORKDIR /build/menu-service
COPY menu-service/go.mod menu-service/go.sum ./
RUN go mod download

COPY menu-service/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
WORKDIR /root/
COPY --from=builder /build/menu-service/main .
EXPOSE 8088
CMD ["./main"]
