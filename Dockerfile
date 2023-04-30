FROM golang:alpine AS builder
RUN apk update
RUN apk add musl-dev gcc
WORKDIR /

ENV CGO_ENABLED=1
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build

FROM alpine as runner
RUN apk update
RUN apk add sqlite ffmpeg --no-cache
WORKDIR /

COPY --from=builder /prince .

CMD ["./prince"]
