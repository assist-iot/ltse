ARG API_VERSION=0.2.0
ARG ALPINE_VERSION

#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY ./src .
RUN go get -d -v ./...
RUN go build -o /go/bin/app

#final stage
FROM alpine:${ALPINE_VERSION:-latest}
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
ENTRYPOINT /app
LABEL Name=ltseapi Version=$API_VERSION
EXPOSE 8080
