FROM python:3.13-slim

# RUN apt-get update && \
#     apt-get install -y curl xz-utils && \
#     rm -rf /var/lib/apt/lists/*

# RUN curl -L https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz \
#     -o ffmpeg.tar.xz && \
#     tar -xf ffmpeg.tar.xz && \
#     cp ffmpeg-*-static/ffmpeg /usr/local/bin/ && \
#     cp ffmpeg-*-static/ffprobe /usr/local/bin/ && \
#     rm -rf ffmpeg*tar.xz ffmpeg-*-static

WORKDIR /app

COPY requirements.txt .
# COPY . .

RUN pip install --no-cache-dir -r requirements.txt

# CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8100"]