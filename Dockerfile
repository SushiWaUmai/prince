FROM golang:alpine AS builder
RUN apk update
WORKDIR /

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build

FROM alpine as runner
RUN apk update
WORKDIR /

COPY --from=builder /prince .

CMD ["./prince"]
