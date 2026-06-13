FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /streamer .

FROM scratch

COPY --from=builder /streamer /streamer

EXPOSE 8080
ENTRYPOINT ["/streamer"]
