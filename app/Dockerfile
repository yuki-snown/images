ARG GO_VERSION=1.13
FROM golang:${GO_VERSION}-alpine AS builder
RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*
WORKDIR /app
ADD . .
RUN go mod download
RUN go build -o ./app ./app.go

FROM alpine:latest
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /image
ADD . /image
COPY --from=builder /app/app /image
EXPOSE 3000
ENTRYPOINT /image/app --port 3000