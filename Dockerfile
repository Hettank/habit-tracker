# ---------- Builder Stage ----------
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o habit-tracker ./cmd/api

# ---------- Runtime Stage ----------
FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/habit-tracker .

EXPOSE 8080

CMD ["./habit-tracker"]