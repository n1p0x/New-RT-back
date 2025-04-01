FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o dep_floor ./cmd/dep_floor/main.go

FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache dcron
COPY --from=builder /app/dep_floor .

RUN echo "0 * * * * /app/dep_floor" > /etc/crontabs/root

CMD ["crond", "-f"]