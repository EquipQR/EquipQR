# Stage 1: Build
FROM golang:1.23-alpine AS BUILDER

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd

WORKDIR /app/cmd/equipqr
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

FROM alpine:latest

WORKDIR /root/

COPY --from=BUILDER /main .

EXPOSE 3000

CMD ["./main"]
