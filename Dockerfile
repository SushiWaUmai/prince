FROM golang:alpine AS builder
RUN apk update
RUN apk add musl-dev gcc alsa-lib-dev
WORKDIR /

ENV CGO_ENABLED=1
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build

FROM alpine as runner
RUN apk update
RUN apk add sqlite ffmpeg py3-pip --no-cache

RUN pip install yt-dlp spotdl

WORKDIR /

COPY --from=builder /prince .

CMD ["./prince"]
