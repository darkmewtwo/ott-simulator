FROM golang:1.24

RUN apt-get update && \
    apt-get install -y curl xz-utils && \
    rm -rf /var/lib/apt/lists/*

RUN curl -L https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz \
    -o ffmpeg.tar.xz && \
    tar -xf ffmpeg.tar.xz && \
    cp ffmpeg-*-static/ffmpeg /usr/local/bin/ && \
    cp ffmpeg-*-static/ffprobe /usr/local/bin/ && \
    rm -rf ffmpeg*tar.xz ffmpeg-*-static

WORKDIR /app

COPY . .

RUN go install github.com/air-verse/air@v1.62

RUN go mod download || true