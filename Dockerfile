FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

RUN apk add --no-cache \
    python3 \
    py3-pip \
    ffmpeg

RUN pip3 install --upgrade yt-dlp

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]