FROM golang:1.19-alpine AS base
ARG MODULE
ENV GOPRIVATE=gitlab.com
RUN apk add --no-cache git

FROM base AS builder
COPY ./src /go/src
WORKDIR /go/src/$MODULE
RUN env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./application

FROM alpine:3.15
ARG MODULE
COPY --from=builder /go/src/$MODULE/config /app/config
COPY --from=builder /go/src/$MODULE/application /app
EXPOSE 8080
WORKDIR /app
ENTRYPOINT ["./application"]