FROM golang:1.24

WORKDIR /app

COPY . .

RUN go install github.com/air-verse/air@v1.62

RUN go mod download || true

# RUN go build -o streamer .

# CMD ["./streamer"]