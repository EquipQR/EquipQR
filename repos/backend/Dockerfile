# Stage 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

# Copy everything (not just ./cmd)
COPY . .

WORKDIR /app/cmd/equipqr

RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /main .

EXPOSE 3000

CMD ["./main"]
