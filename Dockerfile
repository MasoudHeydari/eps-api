# Build
FROM golang:alpine as builder

RUN apk update --no-cache && apk add --no-cache tzdata
WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .
RUN go build -o /app/eps .

FROM alpine:latest
USER root

COPY --from=builder /app/eps /usr/local/bin/eps

ENTRYPOINT ["eps"]

