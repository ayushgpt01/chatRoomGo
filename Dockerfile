# ---------- build stage ----------
FROM golang:1.26.1-alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server


# ---------- runtime stage ----------
FROM alpine:3.23

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]