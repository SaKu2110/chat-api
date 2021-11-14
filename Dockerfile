FROM golang:1.13-alpine AS build-env
WORKDIR /go/src/chat/v1/
COPY ./ ./
RUN go build -o server cmd/main.go

FROM alpine:latest
RUN apk add --no-cache --update ca-certificates
COPY --from=build-env /go/src/chat/v1/server /usr/local/bin/server
ENV API_PORT 8080

EXPOSE 8080
CMD ["/usr/local/bin/server"]